package httpreporter

import (
	"errors"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Builder Struct to hold configurations
type Builder struct {
	// CollectorHostPort are host:ports of a static list of Jaeger Collectors.
	CollectorHostPort []string `yaml:"collectorHTTPHostPort"`
	// Token is an authorization token
	Token string
}

// NewBuilder creates a new reporter builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// CreateReporter creates the TChannel-based Reporter
func (b *Builder) CreateReporter(mFactory metrics.Factory, logger *zap.Logger) (*Reporter, error) {
	if len(b.CollectorHostPort) == 0 {
		return nil, errors.New("collector address string not specified")
	}
	// connect to the first collector for now
	transport, err := thrift.NewTHttpPostClient("https://" + b.CollectorHostPort[0])
	if err != nil {
		return nil, err
	}
	httpTransport, _ := (transport).(*thrift.THttpClient)
	if b.Token != "" {
		httpTransport.SetHeader("Authorization", b.Token)
	}
	return newReporter(httpTransport, thrift.NewTBinaryProtocolFactoryDefault(), logger), nil
}
