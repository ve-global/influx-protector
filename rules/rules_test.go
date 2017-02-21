package rules

import (
	"strings"
	"testing"
	"time"

	"github.com/influxdata/influxdb/influxql"
)

var options = &Options{
	MaxBuckets:  200,
	MaxDuration: time.Hour * 24,
}

func parseQuery(t *testing.T, query string) influxql.Statement {
	stmt, qerr := influxql.NewParser(strings.NewReader(query)).ParseStatement()
	if qerr != nil {
		t.Errorf("failed parsing query: %s", qerr)
	}

	return stmt
}

func TestRunRules(t *testing.T) {
	RunRules(&influxql.ShowSeriesStatement{}, options)
}
