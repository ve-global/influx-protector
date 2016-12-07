package rules

import "testing"

const queryWithoutAggregation = "select * from foo where time > now() - 30m"
const queryWithAggregation = "select sum(*) from foo where time > now() - 30m"
const queryWithoutWhereClause = "select sum(*) from foo"

func TestValidateRawQuery(t *testing.T) {
	query := parseQuery(t, queryWithoutAggregation)

	err := validateRawQuery(query)

	if err == nil {
		t.Errorf("raw queries should not be allowed: %s", query.String())
	}
}

func TestValidateNonRawQuery(t *testing.T) {
	query := parseQuery(t, queryWithAggregation)

	err := validateRawQuery(query)

	if err != nil {
		t.Errorf("non-raw queries should be allowed: %s", err)
	}
}

func TestValidateMissingWhereClause(t *testing.T) {
	query := parseQuery(t, queryWithoutWhereClause)

	err := validateRawQuery(query)

	if err == nil {
		t.Errorf("queries without a where clause should not be allowed: %s", query.String())
	}
}
