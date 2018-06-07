package statstore

import (
	"github.com/jaegertracing/jaeger/model"
)

// Reader defines read behaviour from the stats store
type Reader interface {
	GetStats(params *model.StatSeriesKey) ([]*model.StatSeries, error)
	GetAlertRules(params *model.AlertRuleQueryParams) ([]*model.AlertRule, error)
}
