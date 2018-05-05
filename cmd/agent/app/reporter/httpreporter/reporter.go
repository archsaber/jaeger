package httpreporter

import (
	"context"
	"io/ioutil"
	"sync"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/fsnotify/fsnotify"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"go.uber.org/zap"
)

// Reporter forwards received spans to central collector tier over HTTP.
type Reporter struct {
	thttpClient *thrift.THttpClient

	jClient *jaeger.CollectorClient
	zClient *zipkincore.ZipkinCollectorClient
	logger  *zap.Logger

	// authorization token
	tokenFile string
	token     string

	// sync mechanism for token access
	rw sync.RWMutex
}

func (r *Reporter) getToken() string {
	r.rw.RLock()
	token := r.token
	r.rw.RUnlock()
	return token
}

func (r *Reporter) updateToken() {
	tokenBytes, err := ioutil.ReadFile(r.tokenFile)
	if err != nil {
		r.logger.Error("unable to set token from file")
		return
	}
	r.rw.Lock()
	r.token = string(tokenBytes)
	r.rw.Unlock()
}

func (r *Reporter) watchTokenUpdates(ctx context.Context) {
	r.logger.Info("watching updates on " + r.tokenFile)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.logger.Fatal("unable to watch for token updates")
		return
	}
	defer watcher.Close()

	err = watcher.Add(r.tokenFile)
	if err != nil {
		r.logger.Fatal(err.Error())
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				r.updateToken()
			}
		case err := <-watcher.Errors:
			r.logger.Error(err.Error())
		case <-ctx.Done():
			return
		}
	}
}

// New returns a new Reporter
func newReporter(trans thrift.TTransport, protFactory thrift.TProtocolFactory, tokenFile string,
	zlogger *zap.Logger) *Reporter {

	httpTransport, _ := (trans).(*thrift.THttpClient)
	r := &Reporter{
		thttpClient: httpTransport,
		jClient:     jaeger.NewCollectorClientFactory(httpTransport, protFactory),
		zClient:     zipkincore.NewZipkinCollectorClientFactory(httpTransport, protFactory),
		tokenFile:   tokenFile,
		logger:      zlogger,
	}
	r.updateToken()
	go r.watchTokenUpdates(context.Background())
	return r
}

// EmitZipkinBatch implements EmitZipkinBatch() of Reporter
func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) error {
	r.thttpClient.SetHeader("Authorization", r.getToken())
	_, err := r.zClient.SubmitZipkinBatch(spans)
	return err
}

// EmitBatch implements EmitBatch() of Reporter
func (r *Reporter) EmitBatch(batch *jaeger.Batch) error {
	r.thttpClient.SetHeader("Authorization", r.getToken())
	_, err := r.jClient.SubmitBatches([]*jaeger.Batch{batch})
	return err
}
