package rules

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/influxql"
)

func validateMaxDuration(query influxql.Statement, options *Options) error {
	now := time.Now().UTC()

	if stmt, ok := query.(*influxql.SelectStatement); ok {

		nowValuer := influxql.NowValuer{Now: now}
		stmt.Condition = influxql.Reduce(stmt.Condition, &nowValuer)

		for _, d := range stmt.Dimensions {
			d.Expr = influxql.Reduce(d.Expr, &nowValuer)
		}

		// Replace time(..) references in the query
		stmt.RewriteTimeFields()

		minTime, maxTime, err := influxql.TimeRange(stmt.Condition)

		if maxTime.IsZero() {
			maxTime = now
		}

		if err != nil {
			return err
		}

		duration := maxTime.Sub(minTime)

		if duration > options.MaxDuration {
			return fmt.Errorf("max-duration limit exceeded: (%d/%d)", duration, options.MaxDuration)
		}

		return nil
	}

	return nil
}
