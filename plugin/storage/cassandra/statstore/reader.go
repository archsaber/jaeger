package statstore

import (
	"errors"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/pkg/cassandra"
	"go.uber.org/zap"
)

const (
	serviceStatQuery = `
		SELECT time_stamp, value
		FROM service_stats
		WHERE domain_id = ? AND env = ? AND service_name = ? AND measure = ? AND time_stamp > ? AND time_stamp < ?`
	operationStatQuery = `
		SELECT time_stamp, value
		FROM operation_stats
		WHERE domain_id = ? AND env = ? AND service_name = ? AND operation_name = ? AND measure = ? AND time_stamp > ? AND time_stamp < ?`
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

// GetStats fetches stats data
func (s *StatReader) GetStats(key *model.StatSeriesKey) ([]*model.StatSeries, error) {
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
	return nil
}
