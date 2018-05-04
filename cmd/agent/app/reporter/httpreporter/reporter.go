package httpreporter

import (
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"go.uber.org/zap"
)

// Reporter forwards received spans to central collector tier over HTTP.
type Reporter struct {
	jClient *jaeger.CollectorClient
	zClient *zipkincore.ZipkinCollectorClient
	logger  *zap.Logger
}

// New returns a new Reporter
func newReporter(trans thrift.TTransport, protFactory thrift.TProtocolFactory,
	zlogger *zap.Logger) *Reporter {
	return &Reporter{
		jClient: jaeger.NewCollectorClientFactory(trans, protFactory),
		zClient: zipkincore.NewZipkinCollectorClientFactory(trans, protFactory),
		logger:  zlogger,
	}
}

// EmitZipkinBatch implements EmitZipkinBatch() of Reporter
func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) error {
	_, err := r.zClient.SubmitZipkinBatch(spans)
	return err
}

// EmitBatch implements EmitBatch() of Reporter
func (r *Reporter) EmitBatch(batch *jaeger.Batch) error {
	_, err := r.jClient.SubmitBatches([]*jaeger.Batch{batch})
	return err
}
