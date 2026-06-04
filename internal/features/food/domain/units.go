package domain

import "fmt"

type UnitDimension string

const (
	DimensionMass   UnitDimension = "mass"
	DimensionVolume UnitDimension = "volume"
	DimensionUnit   UnitDimension = "unit"
)

var StandardUnit = map[UnitDimension]string{
	DimensionMass:   "g",
	DimensionVolume: "ml",
	DimensionUnit:   "u",
}

type metricUnitInfo struct {
	Dimension  UnitDimension
	ToStandard float64
}

var metricUnits = map[string]metricUnitInfo{
	"mg": {DimensionMass, 0.001},
	"g":  {DimensionMass, 1},
	"kg": {DimensionMass, 1000},
	"ml": {DimensionVolume, 1},
	"cl": {DimensionVolume, 10},
	"dl": {DimensionVolume, 100},
	"l":  {DimensionVolume, 1000},
	"u":  {DimensionUnit, 1},
}

func IsMetricUnit(unit string) bool {
	_, ok := metricUnits[unit]
	return ok
}

func GetUnitDimension(unit string) (UnitDimension, bool) {
	info, ok := metricUnits[unit]
	if !ok {
		return "", false
	}
	return info.Dimension, true
}

func ConvertToStandard(qty float64, unit string) (float64, string, error) {
	info, ok := metricUnits[unit]
	if !ok {
		return 0, "", fmt.Errorf("unit '%s' is not a metric unit", unit)
	}
	standardUnit := StandardUnit[info.Dimension]
	return qty * info.ToStandard, standardUnit, nil
}

func IsValidMeasurementType(mt string) bool {
	return mt == string(DimensionMass) || mt == string(DimensionVolume) || mt == string(DimensionUnit)
}

func MetricUnitsForDimension(dim UnitDimension) []string {
	var units []string
	for unit, info := range metricUnits {
		if info.Dimension == dim {
			units = append(units, unit)
		}
	}
	return units
}
