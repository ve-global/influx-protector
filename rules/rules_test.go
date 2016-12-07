package rules

import (
	"strings"
	"testing"

	"github.com/influxdata/influxdb/influxql"
)

func parseQuery(t *testing.T, query string) influxql.Statement {
	stmt, qerr := influxql.NewParser(strings.NewReader(query)).ParseStatement()
	if qerr != nil {
		t.Errorf("failed parsing query: %s", qerr)
	}

	return stmt
}
