package ports

import "github.com/google/uuid"

type CrossContextPhoto struct {
	PhotoId    uuid.UUID
	MealId     uuid.UUID
	MealItemId uuid.UUID
}

type ItemGroup struct {
	MealId     uuid.UUID
	MealItemId uuid.UUID
}

type InvalidFoodUnit struct {
	FoodId   uuid.UUID
	BaseUnit string
}

type WatchdogRepository interface {
	AllPhotoURLs() ([]string, error)
	CrossContextPhotos() ([]CrossContextPhoto, error)
	MealGroupsMissingPrimary() ([]uuid.UUID, error)
	ItemGroupsMissingPrimary() ([]ItemGroup, error)
	InvalidFoodBaseUnits() ([]InvalidFoodUnit, error)
}
