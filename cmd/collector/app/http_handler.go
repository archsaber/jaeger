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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DataDog/datadog-trace-agent/model"
	"github.com/DataDog/datadog-trace-agent/writer"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/mux"
	avroKafka "github.com/jaegertracing/jaeger/avro-gen/kafka"
	jModel "github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/jwt"
	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
	goKafka "github.com/segmentio/kafka-go"
	tchanThrift "github.com/uber/tchannel-go/thrift"
)

const (
	formatParam = "format"
	// UnableToReadBodyErrFormat is an error message for invalid requests
	UnableToReadBodyErrFormat = "Unable to process request body: %v"
)

type apiProcessor struct {
	jaegerBatchesHandler JaegerBatchesHandler
	nodeuuid             string
	domainid             string
}

func (aP *apiProcessor) SubmitBatches(batches []*tJaeger.Batch) ([]*tJaeger.BatchSubmitResponse, error) {
	headers := map[string]string{
		"nodeuuid": aP.nodeuuid,
		"domainid": aP.domainid,
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()

	return aP.jaegerBatchesHandler.SubmitBatches(tchanThrift.WithHeaders(ctx, headers), batches)
}

// APIHandler handles all HTTP calls to the collector
type APIHandler struct {
	jaegerBatchesHandler JaegerBatchesHandler
	statsProducer        *goKafka.Writer
}

// NewAPIHandler returns a new APIHandler
func NewAPIHandler(
	jaegerBatchesHandler JaegerBatchesHandler,
	producer *goKafka.Writer,
) *APIHandler {
	return &APIHandler{
		jaegerBatchesHandler: jaegerBatchesHandler,
		statsProducer:        producer,
	}
}

// RegisterRoutes registers routes for this handler on the given router
func (aH *APIHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/traces", aH.saveSpan).Methods(http.MethodPost)
	router.HandleFunc("/api/v0.2/ddtraces", aH.handleDDTraces)
	router.HandleFunc("/api/v0.2/ddstats", aH.handleDDStats)
}

func (aH *APIHandler) handleDDTraces(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.CheckTokenValidity(r.Header.Get(writer.APIHTTPHeaderKey),
		os.Getenv("SECRET_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	apiProcessor, err := aH.newAPIProcessor(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rBodyReader, err := gzip.NewReader(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rBodyUncompressed, err := ioutil.ReadAll(rBodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tracePayload model.TracePayload
	err = proto.Unmarshal(rBodyUncompressed, &tracePayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var jBatches []*tJaeger.Batch
	for _, ddTrace := range tracePayload.Traces {
		var jSpans []*tJaeger.Span
		var t model.Trace = ddTrace.Spans
		for _, ddSpan := range ddTrace.Spans {
			jSpans = append(jSpans, convertDDSpanToJaeger(ddSpan))
		}
		jBatches = append(jBatches, &tJaeger.Batch{
			Process: &tJaeger.Process{
				ServiceName: t.GetRoot().Service,
				Tags: []*tJaeger.Tag{
					&tJaeger.Tag{
						Key:   "env",
						VType: tJaeger.TagType_STRING,
						VStr:  &tracePayload.Env,
					},
				},
			},
			Spans: jSpans,
		})
	}
	_, err = apiProcessor.SubmitBatches(jBatches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func (aH *APIHandler) handleDDStats(w http.ResponseWriter, r *http.Request) {
	tokenClaims, err := jwt.CheckTokenValidity(r.Header.Get(writer.APIHTTPHeaderKey),
		os.Getenv("SECRET_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	geneuuid, ok1 := tokenClaims["geneuuid"]
	domain, ok2 := tokenClaims["domain"]
	if !ok1 || !ok2 {
		http.Error(w, "Incomplete token", http.StatusUnauthorized)
		return
	}
	nodeuuid, ok1 := geneuuid.(string)
	domainid, ok2 := domain.(string)
	if !ok1 || !ok2 {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	rBodyReader, err := gzip.NewReader(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rBodyUncompressed, err := ioutil.ReadAll(rBodyReader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var statsPayload model.StatsPayload
	err = json.Unmarshal(rBodyUncompressed, &statsPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var kafkaMessages []goKafka.Message
	var rawMessage bytes.Buffer
	for _, statBucket := range statsPayload.Stats {
		for _, count := range statBucket.Counts {
			measure := count.Measure
			service := count.TagSet.Get("service").Value
			resource := count.TagSet.Get("resource").Value
			env := count.TagSet.Get("env").Value
			httpstatuscode := count.TagSet.Get("http.status_code").Value
			sublayerService := count.TagSet.Get("sublayer_service").Value
			sublayerType := count.TagSet.Get("sublayer_type").Value
			submeasure := measure
			if measure == "_sublayers.duration.by_service" {
				measure = jModel.DURATION_BY_SERVICE
				submeasure = sublayerService
			} else if measure == "_sublayers.duration.by_type" {
				measure = jModel.DURATION_BY_TYPE
				submeasure = sublayerType
			}

			avro := &avroKafka.StatsPoint{
				Start:      statBucket.Start,
				Duration:   statBucket.Duration,
				Measure:    measure,
				Submeasure: submeasure,
				Value:      count.Value,
				DomainID:   domainid,
				NodeUUID:   nodeuuid,
				Service:    service,
				// Stats are sent only for top level spans
				Operation:      getJaegerOperationName(count.Name, resource, true),
				Env:            env,
				HTTPStatusCode: httpstatuscode,
			}
			rawMessage.Reset()
			err := json.NewEncoder(&rawMessage).Encode(avro)
			if err != nil {
				continue
			}
			kafkaMessages = append(kafkaMessages, goKafka.Message{
				Key:   []byte(avro.DomainID),
				Value: rawMessage.Bytes(),
				Time:  time.Unix(0, avro.Start+avro.Duration),
			})
		}
	}
	aH.statsProducer.WriteMessages(context.Background(), kafkaMessages...)
}

func (aH *APIHandler) saveSpan(w http.ResponseWriter, r *http.Request) {
	format := r.FormValue(formatParam)
	switch strings.ToLower(format) {
	case "jaeger.thrift":
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) < 6 || strings.ToUpper(authHeader[0:6]) != "BEARER" {
			http.Error(w, "Unuthorizated access, this event will be logged and reported",
				http.StatusUnauthorized)
			return
		}
		claims, err := jwt.CheckTokenValidity(authHeader[7:], os.Getenv("SECRET_KEY"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		apiProcessor, err := aH.newAPIProcessor(claims)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		collectorProcessor := tJaeger.NewCollectorProcessor(apiProcessor)
		transport := thrift.NewStreamTransport(r.Body, w)
		protFactory := thrift.NewTBinaryProtocolFactoryDefault()
		collectorProcessor.Process(protFactory.GetProtocol(transport),
			protFactory.GetProtocol(transport))

	default:
		http.Error(w, fmt.Sprintf("Unsupported format type: %v", format), http.StatusBadRequest)
		return
	}
}

func (aH *APIHandler) newAPIProcessor(tokenClaims map[string]interface{}) (*apiProcessor, error) {
	geneuuid, ok1 := tokenClaims["geneuuid"]
	domain, ok2 := tokenClaims["domain"]
	if !ok1 || !ok2 {
		return nil, errors.New("Incomplete token")
	}
	nodeuuid, ok1 := geneuuid.(string)
	domainid, ok2 := domain.(string)
	if !ok1 || !ok2 {
		return nil, errors.New("Invalid token")
	}
	return &apiProcessor{
		jaegerBatchesHandler: aH.jaegerBatchesHandler,
		nodeuuid:             nodeuuid,
		domainid:             domainid,
	}, nil
}

func convertDDSpanToJaeger(s *model.Span) *tJaeger.Span {
	jaegerTags := ddMetaInfoToTags(s.GetMeta())
	ddMetrics := s.GetMetrics()
	for metric := range ddMetrics {
		if strings.HasPrefix(metric, "_") {
			continue
		}
		val := ddMetrics[metric]
		jaegerTags = append(jaegerTags, &tJaeger.Tag{
			Key:     metric,
			VType:   tJaeger.TagType_DOUBLE,
			VDouble: &val,
		})
	}

	jaegerTags = append(jaegerTags, &tJaeger.Tag{
		Key:   "resource",
		VType: tJaeger.TagType_STRING,
		VStr:  &s.Resource,
	})

	jaegerSpan := &tJaeger.Span{
		TraceIdHigh:   0,
		TraceIdLow:    int64(s.TraceID),
		SpanId:        int64(s.SpanID),
		ParentSpanId:  int64(s.ParentID),
		OperationName: getJaegerOperationName(s.Name, s.Resource, s.TopLevel()),
		Flags:         1,
		StartTime:     s.Start / 1e3,
		Duration:      s.Duration / 1e3,
		Tags:          jaegerTags,
	}
	return jaegerSpan
}

func ddMetaInfoToTags(meta map[string]string) []*tJaeger.Tag {
	jaegerTags := make([]*tJaeger.Tag, 0, len(meta)+1)
	for k := range meta {
		val := meta[k]
		jaegerTags = append(jaegerTags, &tJaeger.Tag{
			Key:   k,
			VType: tJaeger.TagType_STRING,
			VStr:  &val,
		})
	}
	return jaegerTags
}

func getJaegerOperationName(ddName, ddResource string, isTopLevel bool) string {
	jaegerOperation := ddName
	if isTopLevel {
		if normResource, ok := model.NormMetricNameParse(ddResource); ok && normResource != ddName {
			jaegerOperation += " - " + ddResource
		}
	}
	return jaegerOperation
}
