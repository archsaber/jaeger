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
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/gorilla/mux"
	tchanThrift "github.com/uber/tchannel-go/thrift"

	jwt "github.com/dgrijalva/jwt-go"
	tJaeger "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

const (
	formatParam = "format"
	// UnableToReadBodyErrFormat is an error message for invalid requests
	UnableToReadBodyErrFormat = "Unable to process request body: %v"
)

type apiProcessor struct {
	jaegerBatchesHandler JaegerBatchesHandler
	tokenClaims          map[string]interface{}
}

func (aP *apiProcessor) SubmitBatches(batches []*tJaeger.Batch) ([]*tJaeger.BatchSubmitResponse, error) {
	headers := map[string]string{
		"nodeuuid": aP.tokenClaims["geneuuid"].(string),
	}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Minute))
	defer cancel()

	return aP.jaegerBatchesHandler.SubmitBatches(tchanThrift.WithHeaders(ctx, headers), batches)
}

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

// RegisterRoutes registers routes for this handler on the given router
func (aH *APIHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/traces", aH.saveSpan).Methods(http.MethodPost)
}

func (aH *APIHandler) saveSpan(w http.ResponseWriter, r *http.Request) {
	format := r.FormValue(formatParam)
	switch strings.ToLower(format) {
	case "jaeger.thrift":
		claims, err := checkTokenValidity(r.Header.Get("Authorization"), os.Getenv("SECRET_KEY"))
		if err != nil {
			http.Error(w, "Could not validate Authorization header", http.StatusUnauthorized)
			return
		}
		collectorProcessor := tJaeger.NewCollectorProcessor(aH.newAPIProcessor(claims))
		transport := thrift.NewStreamTransport(r.Body, w)
		protFactory := thrift.NewTBinaryProtocolFactoryDefault()
		collectorProcessor.Process(protFactory.GetProtocol(transport),
			protFactory.GetProtocol(transport))

	default:
		http.Error(w, fmt.Sprintf("Unsupported format type: %v", format), http.StatusBadRequest)
		return
	}
}

func (aH *APIHandler) newAPIProcessor(tokenClaims map[string]interface{}) *apiProcessor {
	return &apiProcessor{
		tokenClaims:          tokenClaims,
		jaegerBatchesHandler: aH.jaegerBatchesHandler,
	}
}

func checkTokenValidity(token, secret string) (map[string]interface{}, error) {
	if len(token) < 6 || strings.ToUpper(token[0:6]) != "BEARER" {
		return nil, errors.New("Unuthorizated access, this event will be logged and reported")
	}
	token = token[7:]
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		_, ok := claims["geneuuid"].(string)
		if !ok {
			return nil, errors.New("Incomplete token")
		}
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}
