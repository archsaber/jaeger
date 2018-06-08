package model

import (
	"time"

	"github.com/gocql/gocql"
)

type AlertGroupKey struct {
	DomainID      gocql.UUID `json:"domainID"`
	Environment   string     `json:"environment"`
	ServiceName   string     `json:"serviceName"`
	OperationName string     `json:"operationName"`
}

type AlertsQueryParams struct {
	AlertGroupKey
	AllOperations bool      `json:"allOperations"`
	Start         time.Time `json:"startTime"`
	End           time.Time `json:"endTime"`
}

type Alert struct {
	AlertGroupKey
	OpenTime    time.Time
	CloseTime   time.Time
	Measure     string
	Submeasure  string
	ActualValue float64
	Limit       float64
	Upper       bool
	Type        string
	Function    string
	Duration    time.Duration
}
