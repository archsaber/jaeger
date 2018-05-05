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
	CollectorHostPorts []string `yaml:"collectorHTTPHostPorts"`
	TokenFile          string   `yaml:"tokenFile"`
}

// NewBuilder creates a new reporter builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// CreateReporter creates the HTTP-based Reporter
func (b *Builder) CreateReporter(mFactory metrics.Factory, logger *zap.Logger) (*Reporter, error) {
	if len(b.CollectorHostPorts) == 0 {
		return nil, errors.New("collector address string not specified")
	}
	// connect to the first collector for now
	transport, err := thrift.NewTHttpPostClient("https://" + b.CollectorHostPorts[0])
	if err != nil {
		return nil, err
	}
	return newReporter(transport, thrift.NewTBinaryProtocolFactoryDefault(), b.TokenFile, logger), nil
}
