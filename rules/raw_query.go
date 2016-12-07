package rules

import (
	"fmt"

	"github.com/influxdata/influxdb/influxql"
)

func validateRawQuery(query influxql.Statement) error {
	if stmt, ok := query.(*influxql.SelectStatement); ok {

		if stmt.IsRawQuery {
			return fmt.Errorf("raw queries are not allowed")
		}

		if stmt.Condition == nil {
			return fmt.Errorf("queries without a where clause are not allowed")
		}
	}

	return nil
}
