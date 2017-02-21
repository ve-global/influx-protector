package rules

import (
	"testing"

	"github.com/influxdata/influxdb/influxql"
)

const queryWith36h = "select sum(value) from \"foo\" where time > now() - 36h group by time(10m)"
const queryWith12h = "select sum(value) from \"foo\" where time > now() - 12h group by time(10m)"
const queryWithoutTime = "select sum(value) from \"foo\" where foo = \"bar\""

func TestValidateDurationMax(t *testing.T) {
	query := parseQuery(t, queryWith36h)
	err := validateMaxDuration(query, options)

	if err == nil {
		t.Error("max-duration should have been exceeded")
	}
}

func TestValidateDuration(t *testing.T) {
	query := parseQuery(t, queryWith12h)
	err := validateMaxDuration(query, options)

	if err != nil {
		t.Error("max-duration should not have been exceeded")
	}
}

func TestValidateMaxDurationWithNonSelectQuery(t *testing.T) {
	err := validateMaxDuration(&influxql.ShowSeriesStatement{}, options)

	if err != nil {
		t.Error("should not try to validate non-select query")
	}
}

func TestValidateMaxDurationWithNonGroupByQuery(t *testing.T) {
	query := parseQuery(t, queryWithoutGroupBy)

	err := validateMaxDuration(query, options)

	if err != nil {
		t.Error("should not try to validate non-select query")
	}
}

func TestValidateMaxDurationWithNonTimeQuery(t *testing.T) {
	query := parseQuery(t, queryWithoutTime)

	err := validateMaxDuration(query, options)

	if err == nil {
		t.Error("max age should have been exceeded")
	}
}
