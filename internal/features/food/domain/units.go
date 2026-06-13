package domain

import "github.com/ivan-ca97/life/pkg/units"

// UnitDimension and related constants are aliases over pkg/units so that
// food service and handler code continues to work without import changes.
type UnitDimension = units.Dimension

const (
	DimensionMass   = units.DimensionMass
	DimensionVolume = units.DimensionVolume
	DimensionUnit   = units.DimensionUnit
)

var StandardUnit = units.StandardUnit

func IsMetricUnit(unit string) bool               { return units.IsMetricUnit(unit) }
func GetUnitDimension(unit string) (UnitDimension, bool) { return units.GetDimension(unit) }
func ConvertToStandard(qty float64, unit string) (float64, string, error) {
	return units.ConvertToStandard(qty, unit)
}
func IsValidMeasurementType(mt string) bool        { return units.IsValidDimension(mt) }
func MetricUnitsForDimension(dim UnitDimension) []string { return units.ForDimension(dim) }
