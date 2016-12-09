package rules

import (
	"fmt"

	"github.com/influxdata/influxdb/influxql"
)

func validateWhereClause(query influxql.Statement) error {
	if stmt, ok := query.(*influxql.SelectStatement); ok {

		if stmt.Condition == nil {
			return fmt.Errorf("queries without a where clause are not allowed")
		}
	}

	return nil
}
