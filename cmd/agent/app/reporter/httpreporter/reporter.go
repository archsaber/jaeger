package httpreporter

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	"github.com/jaegertracing/jaeger/thrift-gen/zipkincore"
	"go.uber.org/zap"
)

const maxPayloadLength = 1000

// Reporter forwards received spans to central collector tier over HTTP.
type Reporter struct {
	jClient *jaeger.CollectorClient
	zClient *zipkincore.ZipkinCollectorClient
	logger  *zap.Logger

	// builder for this reporter
	builder *Builder

	jBatches chan *jaeger.Batch
	// payload buffer for jaeger batches
	jPayload []*jaeger.Batch

	zBatches chan []*zipkincore.Span
	// payload buffer for zipkin spans
	zPayload []*zipkincore.Span

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
		case <-ctx.Done():
			return
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {
				r.updateClients()
			}
		case err := <-watcher.Errors:
			r.logger.Error(err.Error())
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

func (r *Reporter) flushJBatchesPeriodic(ctx context.Context) {
	tick := time.NewTicker(4 * time.Second)
	defer tick.Stop()

	for range tick.C {
		finished := false
		for {
			select {
			case <-ctx.Done():
				return
			case batch := <-r.jBatches:
				r.jPayload = append(r.jPayload, batch)
			default:
				finished = true
			}
			if finished || len(r.jPayload) >= maxPayloadLength {
				if len(r.jPayload) == 0 {
					break
				}
				r.rw.RLock()
				r.jClient.SubmitBatches(r.jPayload)
				r.rw.RUnlock()

				// reset payload buffer
				r.jPayload = r.jPayload[:0]
				break
			}
		}
	}
}

func (r *Reporter) flushZBatchesPeriodic(ctx context.Context) {
	tick := time.NewTicker(4 * time.Second)
	defer tick.Stop()

	for range tick.C {
		finished := false
		for {
			select {
			case <-ctx.Done():
				return
			case batch := <-r.zBatches:
				r.zPayload = append(r.zPayload, batch...)
			default:
				finished = true
			}
			if finished || len(r.zPayload) >= maxPayloadLength {
				if len(r.zPayload) == 0 {
					break
				}
				r.rw.RLock()
				r.zClient.SubmitZipkinBatch(r.zPayload)
				r.rw.RUnlock()

				// reset payload buffer
				r.zPayload = r.zPayload[:0]
			}
		}
	}
}

// EmitZipkinBatch implements EmitZipkinBatch() of Reporter
func (r *Reporter) EmitZipkinBatch(spans []*zipkincore.Span) error {
	select {
	case r.zBatches <- spans:
		return nil
	default:
		return errors.New("Zipkin spans dropped due to congestion")
	}
}

// EmitBatch implements EmitBatch() of Reporter
func (r *Reporter) EmitBatch(batch *jaeger.Batch) error {
	select {
	case r.jBatches <- batch:
		return nil
	default:
		return errors.New("Jaeger batch dropped due to congestion")
	}
}
