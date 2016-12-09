package rules

import "github.com/influxdata/influxdb/influxql"

// Options for applying rules
type Options struct {
	MaxBuckets int
}

// RunRules runs the validation rules against the statement
func RunRules(query influxql.Statement, options *Options) error {
	if err := validateQueryType(query); err != nil {
		return err
	}

	if err := validateDataPoints(query, options); err != nil {
		return err
	}

	if err := validateWhereClause(query); err != nil {
		return err
	}

	if err := validateSeries(query); err != nil {
		return err
	}

	return nil
}
