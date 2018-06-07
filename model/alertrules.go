package model

import (
	"time"

	"github.com/gocql/gocql"
)

type AlertRuleGroupKey struct {
	DomainID      gocql.UUID `json:"domainID"`
	Environment   string     `json:"environment"`
	ServiceName   string     `json:"serviceName"`
	OperationName string     `json:"operationName"`
}

type AlertRuleQueryParams struct {
	AlertRuleGroupKey
	AllOperations bool `json:"allOperations"`
}

type AlertRule struct {
	AlertRuleGroupKey
	Measure      string        `json:"measure"`
	Submeasure   string        `json:"submeasure"`
	CreationTime time.Time     `json:"creationTime"`
	Duration     time.Duration `json:"duration"`
	Disabled     bool          `json:"disabled"`
	Threshold    float64       `json:"threshold"`
	Upper        bool          `json:"upper"`
	Type         string        `json:"type"`
	Function     string        `json:"function"`
}
