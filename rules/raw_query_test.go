package rules

import "testing"

const queryWithWhereClause = "select * from foo where time > now() - 30m"
const queryWithoutWhereClause = "select sum(value) from foo"

func TestValidateQueryWithoutWhereClause(t *testing.T) {
	query := parseQuery(t, queryWithoutWhereClause)

	err := validateWhereClause(query)

	if err == nil {
		t.Errorf("queries without a where clause should not be allowed: %s", query.String())
	}
}

func TestValidateQueryWithWhereClause(t *testing.T) {
	query := parseQuery(t, queryWithWhereClause)

	err := validateWhereClause(query)

	if err != nil {
		t.Errorf("queries with a where clause should be allowed: %s", err)
	}
}
