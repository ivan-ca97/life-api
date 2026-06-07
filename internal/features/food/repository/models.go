package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/food/domain"
)

type food struct {
	Id                  uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId              uuid.UUID `gorm:"type:uuid;not null"`
	Name                string    `gorm:"not null"`
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string           `gorm:"not null;default:'mass'"`
	BaseQuantity        float64          `gorm:"not null;default:1"`
	BaseUnit            string           `gorm:"not null;default:''"`
	Public              bool             `gorm:"not null;default:false"`
	CreatedAt           time.Time        `gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time        `gorm:"not null;autoUpdateTime"`
	Tags                []foodTag        `gorm:"foreignKey:FoodId"`
	Ingredients         []foodIngredient `gorm:"foreignKey:FoodId"`
	Conversions         []foodConversion `gorm:"foreignKey:FoodId"`
}

type foodTag struct {
	FoodId uuid.UUID `gorm:"type:uuid;primaryKey"`
	Tag    string    `gorm:"primaryKey"`
}

type ingredient struct {
	Id     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId uuid.UUID `gorm:"type:uuid;not null"`
	Name   string    `gorm:"not null"`
}

func (ingredient) TableName() string {
	return "ingredients"
}

type foodIngredient struct {
	FoodId       uuid.UUID  `gorm:"type:uuid;primaryKey"`
	IngredientId uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Ingredient   ingredient `gorm:"foreignKey:IngredientId;references:Id"`
}

type foodConversion struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	FoodId         uuid.UUID `gorm:"type:uuid;not null"`
	Unit           string    `gorm:"not null"`
	BaseEquivalent float64   `gorm:"not null"`
	Inverse        bool      `gorm:"not null;default:false"`
	Note           *string
}

func (foodConversion) TableName() string {
	return "food_conversions"
}

func (f *food) toDomain() *domain.Food {
	tags := make([]string, len(f.Tags))
	for i, t := range f.Tags {
		tags[i] = t.Tag
	}
	ingredients := make([]domain.Ingredient, len(f.Ingredients))
	for i, fi := range f.Ingredients {
		ingredients[i] = domain.Ingredient{
			Id:   fi.Ingredient.Id,
			Name: fi.Ingredient.Name,
		}
	}
	conversions := make([]domain.Conversion, len(f.Conversions))
	for i, c := range f.Conversions {
		note := ""
		if c.Note != nil {
			note = *c.Note
		}
		conversions[i] = domain.Conversion{
			Id:             c.Id,
			Unit:           c.Unit,
			BaseEquivalent: c.BaseEquivalent,
			Inverse:        c.Inverse,
			Note:           note,
		}
	}
	return &domain.Food{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		DefaultCalories:     f.DefaultCalories,
		DefaultProteinGrams: f.DefaultProteinGrams,
		DefaultCarbsGrams:   f.DefaultCarbsGrams,
		DefaultFatGrams:     f.DefaultFatGrams,
		DefaultFiberGrams:   f.DefaultFiberGrams,
		MeasurementType:     f.MeasurementType,
		BaseQuantity:        f.BaseQuantity,
		BaseUnit:            f.BaseUnit,
		Public:              f.Public,
		Tags:                tags,
		Ingredients:         ingredients,
		Conversions:         conversions,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}

func foodFromDomain(f *domain.Food) *food {
	return &food{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		DefaultCalories:     f.DefaultCalories,
		DefaultProteinGrams: f.DefaultProteinGrams,
		DefaultCarbsGrams:   f.DefaultCarbsGrams,
		DefaultFatGrams:     f.DefaultFatGrams,
		DefaultFiberGrams:   f.DefaultFiberGrams,
		MeasurementType:     f.MeasurementType,
		BaseQuantity:        f.BaseQuantity,
		BaseUnit:            f.BaseUnit,
		Public:              f.Public,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}
