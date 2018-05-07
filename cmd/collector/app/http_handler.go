// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gorilla/mux"
	tchanThrift "github.com/uber/tchannel-go/thrift"

	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

const (
	formatParam = "format"
	// UnableToReadBodyErrFormat is an error message for invalid requests
	UnableToReadBodyErrFormat = "Unable to process request body: %v"
)

// APIHandler handles all HTTP calls to the collector
type APIHandler struct {
	jaegerBatchesHandler JaegerBatchesHandler
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(
	jaegerBatchesHandler JaegerBatchesHandler,
) *APIHandler {
	return &APIHandler{
		jaegerBatchesHandler: jaegerBatchesHandler,
	}
}

func (aH *APIHandler) SubmitBatches(batches []*tJaeger.Batch) ([]*tJaeger.BatchSubmitResponse, error) {
	ctx, cancel := tchanThrift.NewContext(time.Minute)
	defer cancel()
	return aH.jaegerBatchesHandler.SubmitBatches(ctx, batches)
}

// RegisterRoutes registers routes for this handler on the given router
func (aH *APIHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/traces", aH.saveSpan).Methods(http.MethodPost)
}

func (aH *APIHandler) saveSpan(w http.ResponseWriter, r *http.Request) {
	format := r.FormValue(formatParam)
	switch strings.ToLower(format) {
	case "jaeger.thrift":
		collectorProcessor := tJaeger.NewCollectorProcessor(aH)
		transport := thrift.NewStreamTransport(r.Body, w)
		protFactory := thrift.NewTBinaryProtocolFactoryDefault()
		collectorProcessor.Process(protFactory.GetProtocol(transport),
			protFactory.GetProtocol(transport))

	default:
		http.Error(w, fmt.Sprintf("Unsupported format type: %v", format), http.StatusBadRequest)
		return
	}
}
