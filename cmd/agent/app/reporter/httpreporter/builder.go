package httpreporter

import (
	"context"
	"errors"
	"io/ioutil"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

// Builder Struct to hold configurations
type Builder struct {
	// CollectorHostPort are host:ports of a static list of Jaeger Collectors.
	CollectorHostPorts []string `yaml:"collectorHTTPHostPorts"`
	TokenFile          string   `yaml:"tokenFile"`

	transport thrift.TTransport
}

// NewBuilder creates a new reporter builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// CreateReporter creates the HTTP-based Reporter
func (b *Builder) CreateReporter(mFactory metrics.Factory, logger *zap.Logger) (*Reporter, error) {
	trans, err := b.getNewTTransport()
	if err != nil {
		return nil, err
	}
	protFactory := b.getNewTProtocolFactory()
	r := &Reporter{
		jClient: jaeger.NewCollectorClientFactory(trans, protFactory),
		zClient: zipkincore.NewZipkinCollectorClientFactory(trans, protFactory),
		logger:  logger,
		builder: b,
	}
	go r.watchTokenUpdates(context.Background())
	return r, nil
}

func (b *Builder) getNewTProtocolFactory() thrift.TProtocolFactory {
	return thrift.NewTBinaryProtocolFactoryDefault()
}

func (b *Builder) getNewTTransport() (thrift.TTransport, error) {
	if b.transport != nil {
		b.transport.Close()
		b.transport = nil
	}
	token := readTokenFromFile(b.TokenFile)
	if len(b.CollectorHostPorts) == 0 {
		return nil, errors.New("collector address string not specified")
	}

	// connect to the first collector for now
	transport, err := thrift.NewTHttpPostClient(b.CollectorHostPorts[0])
	if err != nil {
		return nil, err
	}
	httpTransport, _ := (transport).(*thrift.THttpClient)
	if token != "" {
		httpTransport.SetHeader("Authorization", token)
	}
	b.transport = httpTransport
	return httpTransport, nil
}

func readTokenFromFile(tokenFile string) string {
	tokenBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return ""
	}
	return string(tokenBytes)
}
