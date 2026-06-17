package repository

import (
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/internal/features/meal/domain"
)

type meal struct {
	Id           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	UserId       uuid.UUID   `gorm:"type:uuid;not null"`
	Date         time.Time   `gorm:"type:date;not null"`
	Type         string      `gorm:"not null"`
	Name         string      `gorm:"default:''"`
	EatenAt      *time.Time
	Calories     *float64
	ProteinGrams *float64
	CarbsGrams   *float64
	FatGrams     *float64
	FiberGrams   *float64
	Notes        string      `gorm:"not null;default:''"`
	CreatedAt    time.Time   `gorm:"not null;autoCreateTime"`
	UpdatedAt    time.Time   `gorm:"not null;autoUpdateTime"`
	Tags         []mealTag   `gorm:"foreignKey:MealId"`
	Items        []mealItem  `gorm:"foreignKey:MealId"`
	Photos       []mealPhoto `gorm:"foreignKey:MealId"`
}

type mealTag struct {
	MealId uuid.UUID `gorm:"type:uuid;primaryKey"`
	Tag    string    `gorm:"primaryKey"`
}

type mealItem struct {
	Id                     uuid.UUID  `gorm:"type:uuid;primaryKey"`
	MealId                 uuid.UUID  `gorm:"type:uuid;not null"`
	FoodId                 uuid.UUID  `gorm:"type:uuid;not null"`
	InputQuantity          float64    `gorm:"column:input_quantity;not null;default:1"`
	InputUnit              *string    `gorm:"column:input_unit"`
	InputPortionId         *uuid.UUID `gorm:"column:input_portion_id;type:uuid"`
	InputPortionEquivalent *float64   `gorm:"column:input_portion_equivalent"`
	NormalizedQuantity     *float64   `gorm:"column:normalized_quantity"`
	NormalizedUnit         *string    `gorm:"column:normalized_unit"`
	Calories               *float64
	ProteinGrams           *float64
	CarbsGrams             *float64
	FatGrams               *float64
	FiberGrams             *float64
	Notes                  string       `gorm:"not null;default:''"`
	MeasurementMethod      *string      `gorm:"column:measurement_method"`
	Food                   mealItemFood `gorm:"foreignKey:FoodId;references:Id"`
}

type mealItemFood struct {
	Id   uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name string
}

func (mealItemFood) TableName() string {
	return "foods"
}

type mealPhoto struct {
	Id         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	MealId     uuid.UUID  `gorm:"type:uuid;not null"`
	MealItemId *uuid.UUID `gorm:"type:uuid"`
	Url        string     `gorm:"not null"`
	IsPrimary  bool       `gorm:"not null;default:false"`
	CreatedAt  time.Time  `gorm:"not null;autoCreateTime"`
}

func (mealPhoto) TableName() string {
	return "meal_photos"
}

func (m *meal) toDomain() *domain.Meal {
	tags := make([]string, len(m.Tags))
	for i, t := range m.Tags {
		tags[i] = t.Tag
	}
	items := make([]domain.MealItem, len(m.Items))
	for i, item := range m.Items {
		inputUnit := ""
		if item.InputUnit != nil {
			inputUnit = *item.InputUnit
		}
		normalizedQty := 0.0
		if item.NormalizedQuantity != nil {
			normalizedQty = *item.NormalizedQuantity
		}
		normalizedUnit := ""
		if item.NormalizedUnit != nil {
			normalizedUnit = *item.NormalizedUnit
		}
		var method domain.MeasurementMethod
		if item.MeasurementMethod != nil {
			method = domain.MeasurementMethod(*item.MeasurementMethod)
		}
		items[i] = domain.MealItem{
			Id:                     item.Id,
			MealId:                 item.MealId,
			FoodId:                 item.FoodId,
			FoodName:               item.Food.Name,
			InputQuantity:          item.InputQuantity,
			InputUnit:              inputUnit,
			InputPortionId:         item.InputPortionId,
			InputPortionEquivalent: item.InputPortionEquivalent,
			NormalizedQuantity:     normalizedQty,
			NormalizedUnit:         normalizedUnit,
			Calories:               item.Calories,
			ProteinGrams:           item.ProteinGrams,
			CarbsGrams:             item.CarbsGrams,
			FatGrams:               item.FatGrams,
			FiberGrams:             item.FiberGrams,
			Notes:                  item.Notes,
			MeasurementMethod:      method,
		}
	}
	photos := make([]domain.MealPhoto, len(m.Photos))
	for i, p := range m.Photos {
		photos[i] = domain.MealPhoto{
			Id:         p.Id,
			MealItemId: p.MealItemId,
			Url:        p.Url,
			IsPrimary:  p.IsPrimary,
		}
	}
	return &domain.Meal{
		Id:           m.Id,
		UserId:       m.UserId,
		Date:         m.Date,
		Type:         m.Type,
		Name:         m.Name,
		Photos:       photos,
		EatenAt:      m.EatenAt,
		Calories:     m.Calories,
		ProteinGrams: m.ProteinGrams,
		CarbsGrams:   m.CarbsGrams,
		FatGrams:     m.FatGrams,
		FiberGrams:   m.FiberGrams,
		Tags:         tags,
		Items:        items,
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func mealFromDomain(m *domain.Meal) *meal {
	return &meal{
		Id:           m.Id,
		UserId:       m.UserId,
		Date:         m.Date,
		Type:         m.Type,
		Name:         m.Name,
		EatenAt:      m.EatenAt,
		Calories:     m.Calories,
		ProteinGrams: m.ProteinGrams,
		CarbsGrams:   m.CarbsGrams,
		FatGrams:     m.FatGrams,
		FiberGrams:   m.FiberGrams,
		Notes:        m.Notes,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
