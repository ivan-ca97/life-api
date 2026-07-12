package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/meal/ports"
)

type foodLookupModel struct {
	Id                  uuid.UUID `gorm:"type:uuid;primaryKey"`
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string  `gorm:"not null;default:'mass'"`
	BaseQuantity        float64 `gorm:"not null;default:1"`
	BaseUnit            string  `gorm:"not null;default:''"`
	GramsPerMl          *float64
	UnitBaseEquivalent  *float64
	Portions            []foodPortionLookup `gorm:"foreignKey:FoodId"`
}

func (foodLookupModel) TableName() string { return "foods" }

type foodPortionLookup struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	FoodId         uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	BaseEquivalent float64   `gorm:"not null"`
}

func (foodPortionLookup) TableName() string { return "food_portions" }

type foodLookup struct {
	db *gorm.DB
}

var _ ports.FoodLookup = (*foodLookup)(nil)

func NewFoodLookup(db *gorm.DB) *foodLookup {
	return &foodLookup{db: db}
}

func (r *foodLookup) FindByIds(userId uuid.UUID, ids []uuid.UUID) (map[uuid.UUID]ports.FoodNutrition, error) {
	if len(ids) == 0 {
		return map[uuid.UUID]ports.FoodNutrition{}, nil
	}

	var models []foodLookupModel
	err := r.db.
		Preload("Portions").
		Where("id IN ? AND (user_id = ? OR public = true)", ids, userId).
		Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("looking up foods for nutrition calculation", err)
	}

	result := make(map[uuid.UUID]ports.FoodNutrition, len(models))
	for _, m := range models {
		nutrition := ports.FoodNutrition{
			Id:                  m.Id,
			DefaultCalories:     m.DefaultCalories,
			DefaultProteinGrams: m.DefaultProteinGrams,
			DefaultCarbsGrams:   m.DefaultCarbsGrams,
			DefaultFatGrams:     m.DefaultFatGrams,
			DefaultFiberGrams:   m.DefaultFiberGrams,
			MeasurementType:     m.MeasurementType,
			BaseQuantity:        m.BaseQuantity,
			BaseUnit:            m.BaseUnit,
		}
		if m.GramsPerMl != nil {
			nutrition.VolumeConversion = &ports.FoodVolumeConversion{
				GramsPerMl: *m.GramsPerMl,
			}
		}
		if m.UnitBaseEquivalent != nil {
			nutrition.UnitConversion = &ports.FoodUnitConversion{
				BaseEquivalent: *m.UnitBaseEquivalent,
			}
		}
		portions := make([]ports.FoodPortion, len(m.Portions))
		for i, p := range m.Portions {
			portions[i] = ports.FoodPortion{
				Id:             p.Id,
				Name:           p.Name,
				BaseEquivalent: p.BaseEquivalent,
			}
		}
		nutrition.Portions = portions
		result[m.Id] = nutrition
	}
	return result, nil
}
