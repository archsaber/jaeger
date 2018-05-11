package ddtrace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/tinylib/msgp/msgp"

	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	ddmodel "github.com/jaegertracing/jaeger/proto-gen/ddtrace"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

const (
	maxRequestBodyLength = 10 * 1024 * 1024
	tagTraceHandler      = "handler:traces"
	tagServiceHandler    = "handler:services"
)

// APIVersion is a dumb way to version our collector handlers
type APIVersion string

const (
	// v01 DEPRECATED, FIXME[1.x]
	// Traces: JSON, slice of spans
	// Services: JSON, map[string]map[string][string]
	v01 APIVersion = "v0.1"
	// v02 DEPRECATED, FIXME[1.x]
	// Traces: JSON, slice of traces
	// Services: JSON, map[string]map[string][string]
	v02 APIVersion = "v0.2"
	// v03
	// Traces: msgpack/JSON (Content-Type) slice of traces
	// Services: msgpack/JSON, map[string]map[string][string]
	v03 APIVersion = "v0.3"
	// v04
	// Traces: msgpack/JSON (Content-Type) slice of traces + returns service sampling ratios
	// Services: msgpack/JSON, map[string]map[string][string]
	v04 APIVersion = "v0.4"
)

// ProcessorConfig is the config for DDTraceProcessor
type ProcessorConfig struct {
	Enabled         bool   `yaml:"enabled"`
	HostPort        string `yaml:"hostPort" validate:"nonzero"`
	ConnectionLimit int    `yaml:"connectionLimit"`
	ReceiverTimeout int    `yaml:"receiverTimeout"`
	NumProcessors   int    `yaml:"workers"`
}

// ShouldStartDDTraceProcessor checks if the processor is enabled and
// has non zero workers in the config
func (c ProcessorConfig) ShouldStartDDTraceProcessor() bool {
	return c.Enabled && c.NumProcessors > 0
}

// Processor is a collector that uses HTTP protocol and just holds
// a chan where the spans received are sent one by one
type Processor struct {
	Traces   chan ddmodel.Trace
	Services chan ddmodel.ServicesMetadata

	conf     ProcessorConfig
	server   *http.Server
	reporter reporter.Reporter

	// Underlying transport for this processor to get spans from. If nil, this processor is a no-op.
	ln *RateLimitedListener

	processing           sync.WaitGroup
	maxRequestBodyLength int64
	debug                bool
}

// NewDDTraceProcessor returns a pointer to a new DDServer
func (c ProcessorConfig) NewDDTraceProcessor(
	reporter reporter.Reporter) (*Processor, error) {
	// use buffered channels so that handlers are not waiting on downstream processing
	ddtraceProcessor := Processor{
		Traces:   make(chan ddmodel.Trace, 5000),
		Services: make(chan ddmodel.ServicesMetadata, 50),

		conf:                 c,
		reporter:             reporter,
		maxRequestBodyLength: maxRequestBodyLength,
		debug:                false,
	}

	// initialize the underlying transport
	err := ddtraceProcessor.initTransport()
	if err != nil {
		return nil, err
	}
	return &ddtraceProcessor, nil
}

// Serve starts doing the HTTP server and is ready to receive traces
func (r *Processor) Serve() {
	r.processing.Add(r.conf.NumProcessors)
	for i := 0; i < r.conf.NumProcessors; i++ {
		go r.processDDTraces()
	}

	// FIXME[1.x]: remove all those legacy endpoints + code that goes with it
	http.HandleFunc("/spans", r.httpHandleWithVersion(v01, r.handleTraces))
	http.HandleFunc("/services", r.httpHandleWithVersion(v01, r.handleServices))
	http.HandleFunc("/v0.1/spans", r.httpHandleWithVersion(v01, r.handleTraces))
	http.HandleFunc("/v0.1/services", r.httpHandleWithVersion(v01, r.handleServices))
	http.HandleFunc("/v0.2/traces", r.httpHandleWithVersion(v02, r.handleTraces))
	http.HandleFunc("/v0.2/services", r.httpHandleWithVersion(v02, r.handleServices))
	http.HandleFunc("/v0.3/traces", r.httpHandleWithVersion(v03, r.handleTraces))
	http.HandleFunc("/v0.3/services", r.httpHandleWithVersion(v03, r.handleServices))

	// current collector API
	http.HandleFunc("/v0.4/traces", r.httpHandleWithVersion(v04, r.handleTraces))
	http.HandleFunc("/v0.4/services", r.httpHandleWithVersion(v04, r.handleServices))

	r.server.Serve(r.ln)
}

func (r *Processor) processDDTraces() {
	for t := range r.Traces {
		jaegerSpans := make([]*jaeger.Span, 0, len(t))
		for _, s := range t {
			jaegerSpans = append(jaegerSpans, convertDDSpanToJaeger(s))
		}
		env := t.GetEnv()
		r.reporter.EmitBatch(
			&jaeger.Batch{
				Process: &jaeger.Process{
					ServiceName: t.GetRoot().GetService(),
					Tags: []*jaeger.Tag{
						&jaeger.Tag{
							Key:   "env",
							VType: jaeger.TagType_STRING,
							VStr:  &env,
						},
					},
				},
				Spans: jaegerSpans,
			},
		)
	}
}

func convertDDSpanToJaeger(s *ddmodel.Span) *jaeger.Span {
	jaegerTags := ddMetaInfoToTags(s.GetMeta())
	ddName := s.GetName()
	jaegerTags = append(jaegerTags, &jaeger.Tag{
		Key:   "operation",
		VType: jaeger.TagType_STRING,
		VStr:  &ddName,
	})
	jaegerSpan := &jaeger.Span{
		TraceIdHigh:   0,
		TraceIdLow:    int64(s.GetTraceID()),
		SpanId:        int64(s.GetSpanID()),
		ParentSpanId:  int64(s.GetParentID()),
		OperationName: s.GetResource(),
		Flags:         1,
		StartTime:     s.GetStart() / 1e3,
		Duration:      s.GetDuration() / 1e3,
		Tags:          jaegerTags,
	}
	return jaegerSpan
}

func ddMetaInfoToTags(meta map[string]string) []*jaeger.Tag {
	jaegerTags := make([]*jaeger.Tag, 0, len(meta)+1)
	for k := range meta {
		val := meta[k]
		jaegerTags = append(jaegerTags, &jaeger.Tag{
			Key:   k,
			VType: jaeger.TagType_STRING,
			VStr:  &val,
		})
	}
	return jaegerTags
}

func (r *Processor) initTransport() error {
	listener, err := net.Listen("tcp", r.conf.HostPort)
	if err != nil {
		return fmt.Errorf("cannot listen on %s: %v", r.conf.HostPort, err)
	}

	ln, err := NewRateLimitedListener(listener, r.conf.ConnectionLimit)
	if err != nil {
		return fmt.Errorf("cannot create listener: %v", err)
	}
	r.ln = ln

	timeout := 5 * time.Second
	if r.conf.ReceiverTimeout > 0 {
		timeout = time.Duration(r.conf.ReceiverTimeout) * time.Second
	}
	r.server = &http.Server{
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	go func() {
		ln.Refresh(r.conf.ConnectionLimit)
	}()

	return nil
}

// Stop stops the processor
func (r *Processor) Stop() {
	// close the server
	expiry := time.Now().Add(20 * time.Second) // give it 20 seconds
	ctx, cancel := context.WithDeadline(context.Background(), expiry)
	defer cancel()
	go r.server.Shutdown(ctx)

	// close the workers by closing the channels from which they read
	close(r.Traces)
	close(r.Services)
	r.processing.Wait()
}

func (r *Processor) httpHandle(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.Body = ddmodel.NewLimitedReader(req.Body, r.maxRequestBodyLength)
		defer req.Body.Close()

		fn(w, req)
	}
}

func (r *Processor) httpHandleWithVersion(v APIVersion, f func(APIVersion, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return r.httpHandle(func(w http.ResponseWriter, req *http.Request) {
		contentType := req.Header.Get("Content-Type")
		if contentType == "application/msgpack" && (v == v01 || v == v02) {
			// msgpack is only supported for versions 0.3
			// log.Errorf("rejecting client request, unsupported media type %q", contentType)
			HTTPFormatError([]string{tagTraceHandler, fmt.Sprintf("v:%s", v)}, w)
			return
		}

		f(v, w, req)
	})
}

func (r *Processor) replyTraces(v APIVersion, w http.ResponseWriter) {
	switch v {
	case v01:
		fallthrough
	case v02:
		fallthrough
	case v03:
		// Simple response, simply acknowledge with "OK"
		HTTPOK(w)
	case v04:
		// Return the recommended sampling rate for each service as a JSON.
		// TODO(apoorv): HTTPRateByService(w, r.dynConf)
	}
}

// handleTraces knows how to handle a bunch of traces
func (r *Processor) handleTraces(v APIVersion, w http.ResponseWriter, req *http.Request) {

	traces, ok := getTraces(v, w, req)
	if !ok {
		return
	}

	// We successfuly decoded the payload
	r.replyTraces(v, w)

	// normalize data
	for i := range traces {
		normTrace, err := ddmodel.NormalizeTrace(traces[i])
		if err != nil {

			errorMsg := fmt.Sprintf("dropping trace reason: %s (debug for more info), %v", err, normTrace)

			// avoid truncation in DEBUG mode
			if len(errorMsg) > 150 && !r.debug {
				errorMsg = errorMsg[:150] + "..."
			}
			// log.Errorf(errorMsg)
		} else {

			select {
			case r.Traces <- normTrace:
				// if our downstream consumer is slow, we drop the trace on the floor
				// this is a safety net against us using too much memory
				// when clients flood us
			default:

				// log.Errorf("dropping trace reason: rate-limited")
			}
		}
	}
}

// handleServices handle a request with a list of several services
func (r *Processor) handleServices(v APIVersion, w http.ResponseWriter, req *http.Request) {
	var servicesMeta ddmodel.ServicesMetadata

	contentType := req.Header.Get("Content-Type")
	if err := decodeReceiverPayload(req.Body, &servicesMeta, v, contentType); err != nil {
		// log.Errorf("cannot decode %s services payload: %v", v, err)
		HTTPDecodingError(err, []string{tagServiceHandler, fmt.Sprintf("v:%s", v)}, w)
		return
	}

	HTTPOK(w)

	r.Services <- servicesMeta
}

func getTraces(v APIVersion, w http.ResponseWriter, req *http.Request) (ddmodel.Traces, bool) {
	var traces ddmodel.Traces
	contentType := req.Header.Get("Content-Type")

	switch v {
	case v01:
		// We cannot use decodeReceiverPayload because []model.Span does not
		// implement msgp.Decodable. This hack can be removed once we
		// drop v01 support.
		if contentType != "application/json" && contentType != "text/json" && contentType != "" {
			// log.Errorf("rejecting client request, unsupported media type %q", contentType)
			HTTPFormatError([]string{tagTraceHandler, fmt.Sprintf("v:%s", v)}, w)
			return nil, false
		}

		// in v01 we actually get spans that we have to transform in traces
		var spans []ddmodel.Span
		if err := json.NewDecoder(req.Body).Decode(&spans); err != nil {
			// log.Errorf("cannot decode %s traces payload: %v", v, err)
			HTTPDecodingError(err, []string{tagTraceHandler, fmt.Sprintf("v:%s", v)}, w)
			return nil, false
		}
		traces = ddmodel.TracesFromSpans(spans)
	case v02:
		fallthrough
	case v03:
		fallthrough
	case v04:
		if err := decodeReceiverPayload(req.Body, &traces, v, contentType); err != nil {
			// log.Errorf("cannot decode %s traces payload: %v", v, err)
			HTTPDecodingError(err, []string{tagTraceHandler, fmt.Sprintf("v:%s", v)}, w)
			return nil, false
		}
	default:
		HTTPEndpointNotSupported([]string{tagTraceHandler, fmt.Sprintf("v:%s", v)}, w)
		return nil, false
	}

	return traces, true
}

func decodeReceiverPayload(r io.Reader, dest msgp.Decodable, v APIVersion, contentType string) error {
	switch contentType {
	case "application/msgpack":
		return msgp.Decode(r, dest)

	case "application/json":
		fallthrough
	case "text/json":
		fallthrough
	case "":
		return json.NewDecoder(r).Decode(dest)

	default:
		panic(fmt.Sprintf("unhandled content type %q", contentType))
	}
}
