package rules

import (
	"fmt"
	"strings"

	"github.com/influxdata/influxdb/influxql"
)

const minRegexLength = 5

func validateSeries(query influxql.Statement) error {
	if stmt, ok := query.(*influxql.SelectStatement); ok {

		sources := influxql.Sources(stmt.Sources)

		if sources.HasRegex() {
			regex := sources.String()

			regex = strings.TrimLeft(regex, "/")
			regex = strings.TrimRight(regex, "/")

			if len(regex) < minRegexLength || !strings.Contains(regex, ".") {
				return fmt.Errorf("regex is too short: %s", regex)
			}
		}
	}
	return nil
}
