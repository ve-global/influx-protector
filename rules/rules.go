package rules

import "github.com/influxdata/influxdb/influxql"

// RunRules runs the validation rules against the statement
func RunRules(query influxql.Statement) error {
	if err := validateQueryType(query); err != nil {
		return err
	}

	if err := validateDataPoints(query); err != nil {
		return err
	}

	if err := validateRawQuery(query); err != nil {
		return err
	}

	if err := validateSeries(query); err != nil {
		return err
	}

	return nil
}
