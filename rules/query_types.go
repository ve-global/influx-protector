package rules

import (
	"fmt"
	"reflect"

	"github.com/influxdata/influxdb/influxql"
)

// ValidateQueryType determines if the given query should be allowed
func validateQueryType(stmt influxql.Statement) error {
	if _, ok := stmt.(*influxql.SelectStatement); ok {
		return nil
	}

	if _, ok := stmt.(*influxql.ShowSeriesStatement); ok {
		return nil
	}

	return fmt.Errorf("query type not allowed: %s", reflect.TypeOf(stmt).String())
}
