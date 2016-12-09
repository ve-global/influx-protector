package rules

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb/influxql"
)

func validateDataPoints(query influxql.Statement, options *Options) error {
	now := time.Now().UTC()

	if stmt, ok := query.(*influxql.SelectStatement); ok {
		// Replace instances of "now()" with the current time
		nowValuer := influxql.NowValuer{Now: now}
		stmt.Condition = influxql.Reduce(stmt.Condition, &nowValuer)

		for _, d := range stmt.Dimensions {
			d.Expr = influxql.Reduce(d.Expr, &nowValuer)
		}

		// Replace time(..) references in the query
		stmt.RewriteTimeFields()

		buckets, err := getNumberOfBuckets(stmt, now)

		if err != nil {
			return err
		}

		if buckets > options.MaxBuckets {
			return fmt.Errorf("max-select-buckets limit exceeded: (%d/%d)", buckets, options.MaxBuckets)
		}

		return nil
	}

	return nil
}

func getNumberOfBuckets(stmt *influxql.SelectStatement, now time.Time) (int, error) {
	minTime, maxTime, err := influxql.TimeRange(stmt.Condition)

	if maxTime.IsZero() {
		maxTime = now
	}

	if err != nil {
		return 0, err
	}

	interval, err := stmt.GroupByInterval()

	if err != nil {
		return 0, err
	}

	if interval > 0 {
		// Determine the start and end time matched to the interval (may not match the actual times).
		min := minTime.Truncate(interval)
		max := maxTime.Truncate(interval).Add(interval)

		// Determine the number of buckets by finding the time span and dividing by the interval.
		return int(int64(max.Sub(min)) / int64(interval)), nil
	}

	return 0, nil
}
