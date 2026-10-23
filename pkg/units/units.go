package units

import "fmt"

type Dimension string

const (
	DimensionMass   Dimension = "mass"
	DimensionVolume Dimension = "volume"
)

var StandardUnit = map[Dimension]string{
	DimensionMass:   "g",
	DimensionVolume: "ml",
}

// A metric unit. Values come from ParseUnit, so a Unit is always one of the
// metric units.
type Unit struct {
	name       string
	dimension  Dimension
	toStandard float64
}

var metricUnits = indexByName([]Unit{
	{
		name:       "mg",
		dimension:  DimensionMass,
		toStandard: 0.001,
	},
	{
		name:       "g",
		dimension:  DimensionMass,
		toStandard: 1,
	},
	{
		name:       "kg",
		dimension:  DimensionMass,
		toStandard: 1000,
	},
	{
		name:       "ml",
		dimension:  DimensionVolume,
		toStandard: 1,
	},
	{
		name:       "cl",
		dimension:  DimensionVolume,
		toStandard: 10,
	},
	{
		name:       "dl",
		dimension:  DimensionVolume,
		toStandard: 100,
	},
	{
		name:       "l",
		dimension:  DimensionVolume,
		toStandard: 1000,
	},
})

func indexByName(all []Unit) map[string]Unit {
	byName := make(map[string]Unit, len(all))
	for _, unit := range all {
		byName[unit.name] = unit
	}
	return byName
}

func ParseUnit(name string) (*Unit, error) {
	unit, ok := metricUnits[name]
	if !ok {
		return nil, fmt.Errorf("unit '%s' is not a metric unit", name)
	}
	return &unit, nil
}

func (u Unit) String() string {
	return u.name
}

func (u Unit) Dimension() Dimension {
	return u.dimension
}

func (u Unit) ToStandard(quantity float64) float64 {
	return quantity * u.toStandard
}

func IsMetricUnit(name string) bool {
	_, ok := metricUnits[name]
	return ok
}

func GetDimension(name string) (Dimension, bool) {
	unit, ok := metricUnits[name]
	if !ok {
		return "", false
	}
	return unit.dimension, true
}

func ConvertToStandard(quantity float64, name string) (float64, string, error) {
	unit, ok := metricUnits[name]
	if !ok {
		return 0, "", fmt.Errorf("unit '%s' is not a metric unit", name)
	}
	return unit.ToStandard(quantity), StandardUnit[unit.dimension], nil
}

func IsValidDimension(dimension string) bool {
	return dimension == string(DimensionMass) || dimension == string(DimensionVolume)
}

func ForDimension(dimension Dimension) []string {
	var result []string
	for name, unit := range metricUnits {
		if unit.dimension == dimension {
			result = append(result, name)
		}
	}
	return result
}
