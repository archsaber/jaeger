package ddtrace

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/DataDog/datadog-trace-agent/cmd/ddtrace"
	ddconfig "github.com/DataDog/datadog-trace-agent/config"
	ddmodel "github.com/DataDog/datadog-trace-agent/model"
	"github.com/jaegertracing/jaeger/cmd/agent/app/reporter"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
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
	receiver *ddtrace.HTTPReceiver
	Traces   chan ddmodel.Trace
	Services chan ddmodel.ServicesMetadata

	conf     ProcessorConfig
	reporter reporter.Reporter

	processing sync.WaitGroup
	debug      bool
}

// NewProcessor returns a pointer to a new DDServer
func (c ProcessorConfig) NewProcessor(
	reporter reporter.Reporter) (*Processor, error) {

	traceChan := make(chan ddmodel.Trace, 5000)
	servicesChan := make(chan ddmodel.ServicesMetadata, 50)

	conf := ddconfig.NewDefaultAgentConfig()
	conf.ConnectionLimit = c.ConnectionLimit
	conf.ReceiverTimeout = c.ReceiverTimeout
	hostPortSplit := strings.Split(c.HostPort, ":")
	if len(hostPortSplit) != 2 {
		return nil, errors.New("Invalid host-port string")
	}
	conf.ReceiverHost = hostPortSplit[0]
	receiverPort, err := strconv.ParseInt(hostPortSplit[1], 10, 32)
	if err != nil {
		return nil, err
	}
	conf.ReceiverPort = int(receiverPort)

	// use buffered channels so that handlers are not waiting on downstream processing
	ddtraceProcessor := Processor{
		receiver: ddtrace.NewHTTPReceiver(conf, ddconfig.NewDynamicConfig(), traceChan, servicesChan),
		Traces:   traceChan,
		Services: servicesChan,

		conf:     c,
		reporter: reporter,
		debug:    false,
	}

	return &ddtraceProcessor, nil
}

// Serve starts doing the HTTP server and is ready to receive traces
func (r *Processor) Serve() {
	r.processing.Add(r.conf.NumProcessors)
	for i := 0; i < r.conf.NumProcessors; i++ {
		go r.processDDTraces()
	}
	r.receiver.Run()
}

// Stop stops the processor
func (r *Processor) Stop() {
	// close the workers by closing the channels from which they read
	close(r.Traces)
	close(r.Services)

	// close the server
	r.receiver.Stop()

	r.processing.Wait()
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
					ServiceName: t.GetRoot().Service,
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
	ddName := s.Name
	jaegerTags = append(jaegerTags, &jaeger.Tag{
		Key:   "operation",
		VType: jaeger.TagType_STRING,
		VStr:  &ddName,
	})
	jaegerSpan := &jaeger.Span{
		TraceIdHigh:   0,
		TraceIdLow:    int64(s.TraceID),
		SpanId:        int64(s.SpanID),
		ParentSpanId:  int64(s.ParentID),
		OperationName: s.Resource,
		Flags:         1,
		StartTime:     s.Start / 1e3,
		Duration:      s.Duration / 1e3,
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
