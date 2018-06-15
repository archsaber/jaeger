package statstore

import (
	"errors"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/cassandra"
	"go.uber.org/zap"
)

const (
	serviceBreakupStatQuery = `
		SELECT operation_name, time_stamp, value
		FROM service_breakup
		WHERE domain_id = ? AND env = ? AND service_name = ? AND measure = ? AND time_stamp > ? AND time_stamp < ?`
	serviceStatQuery = `
		SELECT time_stamp, value
		FROM service_stats
		WHERE domain_id = ? AND env = ? AND service_name = ? AND measure = ? AND time_stamp > ? AND time_stamp < ?`
	operationStatQuery = `
		SELECT time_stamp, value
		FROM operation_stats
		WHERE domain_id = ? AND env = ? AND service_name = ? AND operation_name = ? AND measure = ? AND time_stamp > ? AND time_stamp < ?`
	alertRulesQueryByOperation = `
		SELECT operation_name, measure, submeasure, creation_time, disabled, duration, function, threshold, type, upper
		FROM alert_rules
		WHERE domain_id = ? AND env = ? AND service_name = ? AND operation_name = ?`
	alertRulesQueryByService = `
		SELECT operation_name, measure, submeasure, creation_time, disabled, duration, function, threshold, type, upper
		FROM alert_rules
		WHERE domain_id = ? AND env = ? AND service_name = ?`
	alertsQuery = `
		SELECT operation_name, measure, submeasure, start_time, end_time, actual_value, threshold, type, function, upper, duration
		FROM alerts
		WHERE domain_id = ? AND env = ? AND service_name = ? AND start_time >= ? AND start_time <= ?`
)

var (
	// ErrMalformedRequestObject occurs when a request object is nil
	ErrMalformedRequestObject = errors.New("Malformed request object")

	// ErrEnvironmentNotSet occurs when attempting to query with an empty service name
	ErrEnvironmentNotSet = errors.New("Environment Name must be set")

	// ErrServiceNameNotSet occurs when attempting to query with an empty service name
	ErrServiceNameNotSet = errors.New("Service Name must be set")

	// ErrStartAndEndTimeNotSet occurs when start time and end time are not set
	ErrStartAndEndTimeNotSet = errors.New("Start and End Time must be set")

	// ErrStartTimeMinGreaterThanMax occurs when start time min is above start time max
	ErrStartTimeMinGreaterThanMax = errors.New("Start Time Minimum is above Maximum")
)

// StatReader is a reader for stats
type StatReader struct {
	session     cassandra.Session
	consistency cassandra.Consistency
	logger      *zap.Logger
}

// NewStatReader returns a new reader for stats
func NewStatReader(session cassandra.Session, logger *zap.Logger) *StatReader {
	return &StatReader{
		session:     session,
		consistency: cassandra.One,
		logger:      logger,
	}
}

func (s *StatReader) GetAlerts(params *model.AlertsQueryParams) ([]*model.Alert, error) {
	if err := validateAlertsQueryParams(params); err != nil {
		return nil, err
	}

	var q cassandra.Query
	q = s.session.Query(alertsQuery, params.DomainID, params.Environment, params.ServiceName,
		model.TimeAsEpochMicroseconds(params.Start),
		model.TimeAsEpochMicroseconds(params.End))

	i := q.Consistency(s.consistency).Iter()

	var operation, measure, submeasure, function, alertType string
	var actualValue, threshold float64
	var upper bool
	var startTime, endTime, duration int64

	alerts := make([]*model.Alert, 0)
	for i.Scan(&operation, &measure, &submeasure, &startTime, &endTime, &actualValue, &threshold,
		&alertType, &function, &upper, &duration) {

		if (params.AllOperations && operation == "") ||
			(!params.AllOperations && operation != params.OperationName) {
			continue
		}

		groupKey := params.AlertGroupKey
		groupKey.OperationName = operation

		var closeTime time.Time
		if endTime > 0 {
			closeTime = model.EpochMicrosecondsAsTime(uint64(endTime))
		}

		alerts = append(alerts, &model.Alert{
			AlertGroupKey: groupKey,
			Measure:       measure,
			Submeasure:    submeasure,
			OpenTime:      model.EpochMicrosecondsAsTime(uint64(startTime)),
			CloseTime:     closeTime,
			ActualValue:   actualValue,
			Limit:         threshold,
			Upper:         upper,
			Type:          alertType,
			Function:      function,
			Duration:      model.MicrosecondsAsDuration(uint64(duration)),
		})
	}

	return alerts, nil
}

func (s *StatReader) GetAlertRules(params *model.AlertRuleQueryParams) ([]*model.AlertRule, error) {
	if err := validateAlertRuleQueryParams(params); err != nil {
		return nil, err
	}

	var q cassandra.Query
	if params.AllOperations {
		q = s.session.Query(alertRulesQueryByService, params.DomainID, params.Environment,
			params.ServiceName)
	} else {
		q = s.session.Query(alertRulesQueryByOperation, params.DomainID, params.Environment,
			params.ServiceName, params.OperationName)
	}
	i := q.Consistency(s.consistency).Iter()

	var operation, measure, submeasure, function, alertType string
	var threshold float64
	var disabled, upper bool
	var creationTime, duration int64

	rules := make([]*model.AlertRule, 0)
	for i.Scan(&operation, &measure, &submeasure, &creationTime, &disabled, &duration, &function,
		&threshold, &alertType, &upper) {

		groupKey := params.AlertRuleGroupKey
		if params.AllOperations {
			if operation == "" {
				// don't want service level alert rules when query asks for all operation level rules
				continue
			}
			groupKey.OperationName = operation
		}

		rules = append(rules, &model.AlertRule{
			AlertRuleGroupKey: groupKey,
			Measure:           measure,
			Submeasure:        submeasure,
			CreationTime:      model.EpochMicrosecondsAsTime(uint64(creationTime)),
			Disabled:          disabled,
			Duration:          model.MicrosecondsAsDuration(uint64(duration)),
			Function:          function,
			Threshold:         threshold,
			Type:              alertType,
			Upper:             upper,
		})
	}
	return rules, nil
}

// GetStats fetches stats data
func (s *StatReader) GetStats(params *model.StatsSeriesParams) ([]*model.StatSeries, error) {
	key := &params.StatSeriesKey
	if err := validateStatSeriesKey(key); err != nil {
		return nil, err
	}
	measures := []string{}
	if key.Measure == "" {
		measures = append(measures, model.HITS, model.DURATION, model.ERRORS,
			model.DURATION_BY_SERVICE, model.DURATION_BY_TYPE)
	} else {
		measures = []string{key.Measure}
	}

	var retMe []*model.StatSeries

	for _, measure := range measures {
		var q cassandra.Query
		if params.AllOperations {
			q = s.session.Query(serviceBreakupStatQuery, key.DomainID, key.Environment, key.ServiceName,
				measure,
				model.TimeAsEpochMicroseconds(key.StartTime),
				model.TimeAsEpochMicroseconds(key.EndTime))

			i := q.Consistency(s.consistency).Iter()

			var operation string
			var timeStamp int64
			var valuesByOperation = make(map[string][]*model.StatPoint)
			var value map[string]float64
			for i.Scan(&operation, &timeStamp, &value) {
				valuesForOperation, ok := valuesByOperation[operation]
				if !ok {
					valuesForOperation = make([]*model.StatPoint, 0, 200)
				}
				valuesForOperation = append(valuesForOperation, &model.StatPoint{
					Timestamp: timeStamp,
					Value:     value,
				})
				valuesByOperation[operation] = valuesForOperation
			}

			for op := range valuesByOperation {
				keyWithMeasure := *key
				keyWithMeasure.Measure = measure
				keyWithMeasure.OperationName = op
				retMe = append(retMe, &model.StatSeries{
					StatSeriesKey: keyWithMeasure,
					Values:        valuesByOperation[op],
				})
			}
			continue
		}

		if key.OperationName != "" {
			q = s.session.Query(operationStatQuery, key.DomainID, key.Environment, key.ServiceName,
				key.OperationName, measure, model.TimeAsEpochMicroseconds(key.StartTime),
				model.TimeAsEpochMicroseconds(key.EndTime))
		} else {
			q = s.session.Query(serviceStatQuery, key.DomainID, key.Environment, key.ServiceName,
				measure,
				model.TimeAsEpochMicroseconds(key.StartTime),
				model.TimeAsEpochMicroseconds(key.EndTime))
		}
		i := q.Consistency(s.consistency).Iter()

		var values []*model.StatPoint
		var timeStamp int64
		var value map[string]float64
		for i.Scan(&timeStamp, &value) {
			values = append(values, &model.StatPoint{
				Timestamp: timeStamp,
				Value:     value,
			})
		}
		keyWithMeasure := *key
		keyWithMeasure.Measure = measure
		retMe = append(retMe, &model.StatSeries{
			StatSeriesKey: keyWithMeasure,
			Values:        values,
		})
	}

	return retMe, nil
}

func validateStatSeriesKey(k *model.StatSeriesKey) error {
	if k == nil {
		return ErrMalformedRequestObject
	}
	if k.Environment == "" {
		return ErrEnvironmentNotSet
	}
	if k.ServiceName == "" {
		return ErrServiceNameNotSet
	}
	if k.StartTime.IsZero() || k.EndTime.IsZero() {
		return ErrStartAndEndTimeNotSet
	}
	if !k.StartTime.IsZero() && !k.EndTime.IsZero() && k.EndTime.Before(k.StartTime) {
		return ErrStartTimeMinGreaterThanMax
	}
	return nil
}

func validateAlertRuleQueryParams(k *model.AlertRuleQueryParams) error {
	if k == nil {
		return ErrMalformedRequestObject
	}
	if k.Environment == "" {
		return ErrEnvironmentNotSet
	}
	if k.ServiceName == "" {
		return ErrServiceNameNotSet
	}
	return nil
}

func validateAlertsQueryParams(k *model.AlertsQueryParams) error {
	if k == nil {
		return ErrMalformedRequestObject
	}
	if k.Environment == "" {
		return ErrEnvironmentNotSet
	}
	if k.ServiceName == "" {
		return ErrServiceNameNotSet
	}
	if k.Start.IsZero() || k.End.IsZero() {
		return ErrStartAndEndTimeNotSet
	}
	if !k.Start.IsZero() && !k.End.IsZero() && k.End.Before(k.Start) {
		return ErrStartTimeMinGreaterThanMax
	}
	return nil
}
