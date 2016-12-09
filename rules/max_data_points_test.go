package rules

import (
	"testing"

	"github.com/influxdata/influxdb/influxql"
)

const queryWith86400Buckets = "select sum(value) from \"foo\" where time > now() - 24h group by time(1s)"
const queryWith144Buckets = "select sum(value) from \"foo\" where time > now() - 24h group by time(10m)"
const queryWithoutGroupBy = "select sum(value) from \"foo\" where time > now() - 24h"

var options = &Options{
	MaxBuckets: 200,
}

func TestValidateDataPointsMax(t *testing.T) {
	query := parseQuery(t, queryWith86400Buckets)
	err := validateDataPoints(query, options)

	if err == nil {
		t.Error("max-select-buckets should have been exceeded")
	}
}

func TestValidateDataPoints(t *testing.T) {
	query := parseQuery(t, queryWith144Buckets)
	err := validateDataPoints(query, options)

	if err != nil {
		t.Error("max-select-buckets should not have been exceeded")
	}
}

func TestValidateDataPointsWithNonSelectQuery(t *testing.T) {
	err := validateDataPoints(&influxql.ShowSeriesStatement{}, options)

	if err != nil {
		t.Error("should not try to validate non-select query")
	}
}

func TestValidateDataPointsWithNonGroupByQuery(t *testing.T) {
	query := parseQuery(t, queryWithoutGroupBy)

	err := validateDataPoints(query, options)

	if err != nil {
		t.Error("should not try to validate non-select query")
	}
}
