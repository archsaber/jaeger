package model

import (
	"time"

	"github.com/gocql/gocql"
)

// Different kinds of stats (measures) that are stored
const (
	HITS     = "hits"
	ERRORS   = "errors"
	DURATION = "duration"
)

// StatPoint represents a single data point
type StatPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// StatSeriesKey is the primary key for a stats time-series
type StatSeriesKey struct {
	DomainID      gocql.UUID `json:"domainID"`
	Environment   string     `json:"environment"`
	ServiceName   string     `json:"serviceName"`
	OperationName string     `json:"operationName"`
	Measure       string     `json:"measure"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       time.Time  `json:"endTime"`
}

// StatSeries represents a single stats time-series
type StatSeries struct {
	StatSeriesKey
	Values []*StatPoint `json:"values"`
}
