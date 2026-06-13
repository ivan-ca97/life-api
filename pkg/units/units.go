package units

import "fmt"

type Dimension string

const (
	DimensionMass   Dimension = "mass"
	DimensionVolume Dimension = "volume"
	DimensionUnit   Dimension = "unit"
)

var StandardUnit = map[Dimension]string{
	DimensionMass:   "g",
	DimensionVolume: "ml",
	DimensionUnit:   "u",
}

type metricUnitInfo struct {
	Dimension  Dimension
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

func GetDimension(unit string) (Dimension, bool) {
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
	return qty * info.ToStandard, StandardUnit[info.Dimension], nil
}

func IsValidDimension(d string) bool {
	return d == string(DimensionMass) || d == string(DimensionVolume) || d == string(DimensionUnit)
}

func ForDimension(dim Dimension) []string {
	var result []string
	for unit, info := range metricUnits {
		if info.Dimension == dim {
			result = append(result, unit)
		}
	}
	return result
}
