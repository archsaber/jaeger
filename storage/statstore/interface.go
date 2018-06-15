package statstore

import (
	"github.com/jaegertracing/jaeger/model"
)

// Reader defines read behaviour from the stats store
type Reader interface {
	GetStats(params *model.StatsSeriesParams) ([]*model.StatSeries, error)
	GetAlertRules(params *model.AlertRuleQueryParams) ([]*model.AlertRule, error)
	GetAlerts(params *model.AlertsQueryParams) ([]*model.Alert, error)
}
