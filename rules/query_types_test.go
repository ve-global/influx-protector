package rules

import (
	"testing"

	"github.com/influxdata/influxdb/influxql"
)

func TestValidateQueryTypeSelectStmt(t *testing.T) {
	err := validateQueryType(&influxql.SelectStatement{})

	if err != nil {
		t.Errorf("SelectStatement should be allowed: %s", err)
	}
}

func TestValidateQueryTypeShowSeriesStmt(t *testing.T) {
	err := validateQueryType(&influxql.ShowSeriesStatement{})

	if err != nil {
		t.Errorf("ShowSeriesStatement should be allowed: %s", err)
	}
}

func TestValidateQueryTypeDropMeasurmentStmt(t *testing.T) {
	err := validateQueryType(&influxql.DropMeasurementStatement{})

	if err == nil {
		t.Errorf("DropMeasurementStatement should not be allowed: %s", err)
	}
}

func TestValidateQueryTypeShowMeasurmentsStmt(t *testing.T) {
	err := validateQueryType(&influxql.ShowMeasurementsStatement{})

	if err == nil {
		t.Errorf("ShowMeasurementsStatement should not be allowed: %s", err)
	}
}

func TestValidateQueryTypeDropDatabaseStmt(t *testing.T) {
	err := validateQueryType(&influxql.DropDatabaseStatement{})

	if err == nil {
		t.Errorf("DropDatabaseStatement should not be allowed: %s", err)
	}
}
