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
	PhotoUrl            string    `gorm:"not null;default:''"`
	DefaultCalories     *float64
	DefaultProteinGrams *float64
	DefaultCarbsGrams   *float64
	DefaultFatGrams     *float64
	DefaultFiberGrams   *float64
	MeasurementType     string    `gorm:"not null;default:'mass'"`
	BaseQuantity        float64   `gorm:"not null;default:1"`
	BaseUnit            string    `gorm:"not null;default:''"`
	Public              bool      `gorm:"not null;default:false"`
	GramsPerMl          *float64
	VolumeNote          *string
	UnitBaseEquivalent  *float64
	UnitNote            *string
	CreatedAt           time.Time     `gorm:"not null;autoCreateTime"`
	UpdatedAt           time.Time     `gorm:"not null;autoUpdateTime"`
	Tags                []foodTag     `gorm:"foreignKey:FoodId"`
	Ingredients         []foodIngredient `gorm:"foreignKey:FoodId"`
	Portions            []foodPortion `gorm:"foreignKey:FoodId"`
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

type foodPortion struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	FoodId         uuid.UUID `gorm:"type:uuid;not null"`
	Name           string    `gorm:"not null"`
	BaseEquivalent float64   `gorm:"not null"`
}

func (foodPortion) TableName() string {
	return "food_portions"
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
	portions := make([]domain.Portion, len(f.Portions))
	for i, p := range f.Portions {
		portions[i] = domain.Portion{
			Id:             p.Id,
			Name:           p.Name,
			BaseEquivalent: p.BaseEquivalent,
		}
	}

	var volumeConv *domain.VolumeConversion
	if f.GramsPerMl != nil {
		note := ""
		if f.VolumeNote != nil {
			note = *f.VolumeNote
		}
		volumeConv = &domain.VolumeConversion{
			GramsPerMl: *f.GramsPerMl,
			Note:       note,
		}
	}

	var unitConv *domain.UnitConversion
	if f.UnitBaseEquivalent != nil {
		note := ""
		if f.UnitNote != nil {
			note = *f.UnitNote
		}
		unitConv = &domain.UnitConversion{
			BaseEquivalent: *f.UnitBaseEquivalent,
			Note:           note,
		}
	}

	return &domain.Food{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		PhotoUrl:            f.PhotoUrl,
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
		VolumeConversion:    volumeConv,
		UnitConversion:      unitConv,
		Portions:            portions,
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
}

func foodFromDomain(f *domain.Food) *food {
	model := &food{
		Id:                  f.Id,
		UserId:              f.UserId,
		Name:                f.Name,
		PhotoUrl:            f.PhotoUrl,
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
	if f.VolumeConversion != nil {
		model.GramsPerMl = &f.VolumeConversion.GramsPerMl
		if f.VolumeConversion.Note != "" {
			model.VolumeNote = &f.VolumeConversion.Note
		}
	}
	if f.UnitConversion != nil {
		model.UnitBaseEquivalent = &f.UnitConversion.BaseEquivalent
		if f.UnitConversion.Note != "" {
			model.UnitNote = &f.UnitConversion.Note
		}
	}
	return model
}
