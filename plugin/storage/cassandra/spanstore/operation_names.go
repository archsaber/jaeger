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

package spanstore

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/pkg/cache"
	"github.com/jaegertracing/jaeger/pkg/cassandra"
	casMetrics "github.com/jaegertracing/jaeger/pkg/cassandra/metrics"
)

const (
	insertOperationName = `INSERT INTO operation_names(service_name, operation_name, domain_id) VALUES (?, ?, ?)`
	queryOperationNames = `SELECT operation_name FROM operation_names WHERE service_name = ? AND domain_id = ?`
)

// OperationNamesStorage stores known operation names by service.
type OperationNamesStorage struct {
	// CQL statements are public so that Cassandra2 storage can override them
	InsertStmt     string
	QueryStmt      string
	session        cassandra.Session
	writeCacheTTL  time.Duration
	metrics        *casMetrics.Table
	operationNames cache.Cache
	logger         *zap.Logger
}

// NewOperationNamesStorage returns a new OperationNamesStorage
func NewOperationNamesStorage(
	session cassandra.Session,
	writeCacheTTL time.Duration,
	metricsFactory metrics.Factory,
	logger *zap.Logger,
) *OperationNamesStorage {
	return &OperationNamesStorage{
		session:       session,
		InsertStmt:    insertOperationName,
		QueryStmt:     queryOperationNames,
		metrics:       casMetrics.NewTable(metricsFactory, "OperationNames"),
		writeCacheTTL: writeCacheTTL,
		logger:        logger,
		operationNames: cache.NewLRUWithOptions(
			100000,
			&cache.Options{
				TTL:             writeCacheTTL,
				InitialCapacity: 10000,
			}),
	}
}

// Write saves Operation and Service name tuples
func (s *OperationNamesStorage) Write(serviceName string, operationName string, domainID gocql.UUID) error {
	var err error
	query := s.session.Query(s.InsertStmt)
	if inCache := checkWriteCache(domainID.String()+"|"+serviceName+"|"+operationName, s.operationNames, s.writeCacheTTL); !inCache {
		q := query.Bind(serviceName, operationName, domainID)
		err2 := s.metrics.Exec(q, s.logger)
		if err2 != nil {
			err = err2
		}
	}
	return err
}

// GetOperations returns all operations for a specific service traced by Jaeger
func (s *OperationNamesStorage) GetOperations(service string, domainID gocql.UUID) ([]string, error) {
	iter := s.session.Query(s.QueryStmt, service, domainID).Iter()

	var operation string
	var operations []string
	for iter.Scan(&operation) {
		operations = append(operations, operation)
	}
	if err := iter.Close(); err != nil {
		err = errors.Wrap(err, "Error reading operation_names from storage")
		return nil, err
	}
	return operations, nil
}
