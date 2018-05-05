package httpreporter

import (
	"context"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"go.uber.org/zap"
)

// Reporter forwards received spans to central collector tier over HTTP.
type Reporter struct {
	jClient *jaeger.CollectorClient
	zClient *zipkincore.ZipkinCollectorClient
	logger  *zap.Logger

	// builder for this reporter
	builder *Builder

	// sync mechanism for jClient and zClient access
	rw sync.RWMutex
}

func (r *Reporter) watchTokenUpdates(ctx context.Context) {
	r.logger.Info("watching updates on " + r.builder.TokenFile)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.logger.Fatal("unable to watch for token updates")
		return
	}
	defer watcher.Close()

	err = watcher.Add(r.builder.TokenFile)
	if err != nil {
		r.logger.Fatal(err.Error())
		return
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				r.updateClients()
			}
		case err := <-watcher.Errors:
			r.logger.Error(err.Error())
		case <-ctx.Done():
			return
		}
	}
}

func (r *Reporter) updateClients() {
	protFactory := r.builder.getNewTProtocolFactory()
	trans, err := r.builder.getNewTTransport()
	if err != nil {
		r.logger.Error(err.Error())
		return
	}
	r.rw.Lock()
	r.jClient = jaeger.NewCollectorClientFactory(trans, protFactory)
	r.zClient = zipkincore.NewZipkinCollectorClientFactory(trans, protFactory)
	r.rw.Unlock()
}

// EmitZipkinBatch implements EmitZipkinBatch() of Reporter
func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) error {
	r.rw.RLock()
	_, err := r.zClient.SubmitZipkinBatch(spans)
	r.rw.RUnlock()
	return err
}

// EmitBatch implements EmitBatch() of Reporter
func (r *Reporter) EmitBatch(batch *jaeger.Batch) error {
	r.rw.RLock()
	_, err := r.jClient.SubmitBatches([]*jaeger.Batch{batch})
	r.rw.RUnlock()
	return err
}
