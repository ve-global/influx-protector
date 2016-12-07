package rules

import "testing"

const queryWithRegex = "select sum(*) from /foo.bar/ where time > now() - 20m group by time(60s)"
const queryWithShortRegex = "select sum(*) from /foo/ where time > now() - 20m group by time(60s)"
const queryWithoutRegex = "select sum(*) from \"foo\" where time > now() - 20m group by time(60s)"

func TestValidateSeriesWithShortRegex(t *testing.T) {
	query := parseQuery(t, queryWithShortRegex)
	err := validateSeries(query)

	if err == nil {
		t.Error("short regexes should not be allowed")
	}
}

func TestValidateSeriesWithRegex(t *testing.T) {
	query := parseQuery(t, queryWithRegex)
	err := validateSeries(query)

	if err != nil {
		t.Error("regexes should be allowed")
	}
}

func TestValidateSeriesWithoutRegex(t *testing.T) {
	query := parseQuery(t, queryWithoutRegex)
	err := validateSeries(query)

	if err != nil {
		t.Error("non-regex queries should be allowed")
	}
}
