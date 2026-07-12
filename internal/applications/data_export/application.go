package data_export

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api"
	"github.com/ivan-ca97/life/pkg/api/endpoint"
	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/permissions"
)

type DataExportApplication struct {
	db           *gorm.DB
	authorizer   auth.AuthorizationService
	errorHandler http_errors.HttpErrorHandler
}

func NewDataExportApplication(db *gorm.DB, authorizer auth.AuthorizationService, errorHandler http_errors.HttpErrorHandler) *DataExportApplication {
	return &DataExportApplication{
		db:           db,
		authorizer:   authorizer,
		errorHandler: errorHandler,
	}
}

func (a *DataExportApplication) ProtectedRoutes(r chi.Router) {
	r.Get("/export", endpoint.JSON(a.errorHandler, a.Export))
}

// --- response types ---

type userExport struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Username  *string   `json:"username,omitempty"`
	HeightCm  *int      `json:"height_cm,omitempty"`
	BirthDate *string   `json:"birth_date,omitempty"`
	Sex       *string   `json:"sex,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type weightEntryExport struct {
	Id                string    `json:"id"`
	Date              string    `json:"date"`
	WeightKg          float64   `json:"weight_kg"`
	BodyFatPercentage *float64  `json:"body_fat_percentage,omitempty"`
	Notes             string    `json:"notes,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type exerciseExport struct {
	Id                      string    `json:"id"`
	Date                    string    `json:"date"`
	Type                    string    `json:"type"`
	Name                    string    `json:"name"`
	DurationSeconds         *int      `json:"duration_seconds,omitempty"`
	EstimatedCaloriesBurned *float64  `json:"estimated_calories_burned,omitempty"`
	Steps                   *int      `json:"steps,omitempty"`
	DistanceMeters          *float64  `json:"distance_meters,omitempty"`
	ElevationGainMeters     *float64  `json:"elevation_gain_meters,omitempty"`
	AverageHeartRate        *int      `json:"average_heart_rate,omitempty"`
	MaxHeartRate            *int      `json:"max_heart_rate,omitempty"`
	TotalVolumeKg           *float64  `json:"total_volume_kg,omitempty"`
	TotalSets               *int      `json:"total_sets,omitempty"`
	Notes                   string    `json:"notes,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
}

type mealExport struct {
	Id           string    `json:"id"`
	Date         string    `json:"date"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Calories     *float64  `json:"calories,omitempty"`
	ProteinGrams *float64  `json:"protein_grams,omitempty"`
	CarbsGrams   *float64  `json:"carbs_grams,omitempty"`
	FatGrams     *float64  `json:"fat_grams,omitempty"`
	FiberGrams   *float64  `json:"fiber_grams,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type foodExport struct {
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	MeasurementType     string    `json:"measurement_type"`
	BaseUnit            string    `json:"base_unit"`
	BaseQuantity        *float64  `json:"base_quantity,omitempty"`
	DefaultCalories     *float64  `json:"default_calories,omitempty"`
	DefaultProteinGrams *float64  `json:"default_protein_grams,omitempty"`
	DefaultCarbsGrams   *float64  `json:"default_carbs_grams,omitempty"`
	DefaultFatGrams     *float64  `json:"default_fat_grams,omitempty"`
	DefaultFiberGrams   *float64  `json:"default_fiber_grams,omitempty"`
	Public              bool      `json:"public"`
	CreatedAt           time.Time `json:"created_at"`
}

type bodyMeasurementExport struct {
	Date      string    `json:"date"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Notes     string    `json:"notes,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type goalExport struct {
	DailyCalories        *float64  `json:"daily_calories,omitempty"`
	DailyProteinGrams    *float64  `json:"daily_protein_grams,omitempty"`
	DailyCarbsGrams      *float64  `json:"daily_carbs_grams,omitempty"`
	DailyFatGrams        *float64  `json:"daily_fat_grams,omitempty"`
	DailyFiberGrams      *float64  `json:"daily_fiber_grams,omitempty"`
	DailySteps           *int      `json:"daily_steps,omitempty"`
	DailyExerciseMinutes *int      `json:"daily_exercise_minutes,omitempty"`
	TargetWeightKg       *float64  `json:"target_weight_kg,omitempty"`
	StartedAt            time.Time `json:"started_at"`
}

type exportResponse struct {
	ExportedAt       string                  `json:"exported_at"`
	User             *userExport             `json:"user"`
	WeightEntries    []weightEntryExport     `json:"weight_entries"`
	Exercises        []exerciseExport        `json:"exercises"`
	Meals            []mealExport            `json:"meals"`
	Foods            []foodExport            `json:"foods"`
	BodyMeasurements []bodyMeasurementExport `json:"body_measurements"`
	Goal             *goalExport             `json:"goal"`
}

// --- handler ---

func (a *DataExportApplication) Export(r *http.Request) (*exportResponse, int, error) {
	userId, err := api.PathParamUUID(r, "userId")
	if err != nil {
		return nil, 0, err
	}
	err = a.authorizer.Authorize(r.Context(), userId, permissions.UsersRead)
	if err != nil {
		return nil, 0, err
	}
	data, err := a.buildExport(userId)
	if err != nil {
		return nil, 0, err
	}
	return data, http.StatusOK, nil
}

func (a *DataExportApplication) buildExport(userId uuid.UUID) (*exportResponse, error) {
	response := &exportResponse{
		ExportedAt:       time.Now().UTC().Format(time.RFC3339),
		WeightEntries:    []weightEntryExport{},
		Exercises:        []exerciseExport{},
		Meals:            []mealExport{},
		Foods:            []foodExport{},
		BodyMeasurements: []bodyMeasurementExport{},
	}

	// user profile
	var userRow struct {
		Id        uuid.UUID  `gorm:"column:id"`
		Email     string     `gorm:"column:email"`
		Username  *string    `gorm:"column:username"`
		HeightCm  *int       `gorm:"column:height_cm"`
		BirthDate *time.Time `gorm:"column:birth_date"`
		Sex       *string    `gorm:"column:sex"`
		CreatedAt time.Time  `gorm:"column:created_at"`
	}
	err := a.db.Raw(
		`SELECT id, email, username, height_cm, birth_date, sex, created_at FROM users WHERE id = ?`,
		userId,
	).Scan(&userRow).Error
	if err != nil {
		return nil, err
	}
	response.User = &userExport{
		Id:        userRow.Id.String(),
		Email:     userRow.Email,
		Username:  userRow.Username,
		HeightCm:  userRow.HeightCm,
		Sex:       userRow.Sex,
		CreatedAt: userRow.CreatedAt,
	}
	if userRow.BirthDate != nil {
		s := userRow.BirthDate.Format("2006-01-02")
		response.User.BirthDate = &s
	}

	// weight entries
	var weightRows []struct {
		Id                uuid.UUID `gorm:"column:id"`
		Date              time.Time `gorm:"column:date"`
		WeightKg          float64   `gorm:"column:weight_kg"`
		BodyFatPercentage *float64  `gorm:"column:body_fat_percentage"`
		Notes             string    `gorm:"column:notes"`
		CreatedAt         time.Time `gorm:"column:created_at"`
	}
	err = a.db.Raw(
		`SELECT id, date, weight_kg, body_fat_percentage, notes, created_at FROM weight_entries WHERE user_id = ? ORDER BY date`,
		userId,
	).Scan(&weightRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range weightRows {
		item := weightEntryExport{
			Id:                row.Id.String(),
			Date:              row.Date.Format("2006-01-02"),
			WeightKg:          row.WeightKg,
			BodyFatPercentage: row.BodyFatPercentage,
			Notes:             row.Notes,
			CreatedAt:         row.CreatedAt,
		}
		response.WeightEntries = append(response.WeightEntries, item)
	}

	// exercises
	var exerciseRows []struct {
		Id                      uuid.UUID `gorm:"column:id"`
		Date                    time.Time `gorm:"column:date"`
		Type                    string    `gorm:"column:type"`
		Name                    string    `gorm:"column:name"`
		DurationSeconds         *int      `gorm:"column:duration_seconds"`
		EstimatedCaloriesBurned *float64  `gorm:"column:estimated_calories_burned"`
		Steps                   *int      `gorm:"column:steps"`
		DistanceMeters          *float64  `gorm:"column:distance_meters"`
		ElevationGainMeters     *float64  `gorm:"column:elevation_gain_meters"`
		AverageHeartRate        *int      `gorm:"column:average_heart_rate"`
		MaxHeartRate            *int      `gorm:"column:max_heart_rate"`
		TotalVolumeKg           *float64  `gorm:"column:total_volume_kg"`
		TotalSets               *int      `gorm:"column:total_sets"`
		Notes                   string    `gorm:"column:notes"`
		CreatedAt               time.Time `gorm:"column:created_at"`
	}
	err = a.db.Raw(`
		SELECT id, date, type, name, duration_seconds, estimated_calories_burned, steps,
		       distance_meters, elevation_gain_meters, average_heart_rate, max_heart_rate,
		       total_volume_kg, total_sets, notes, created_at
		FROM exercises WHERE user_id = ? ORDER BY date
	`, userId).Scan(&exerciseRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range exerciseRows {
		item := exerciseExport{
			Id:                      row.Id.String(),
			Date:                    row.Date.Format("2006-01-02"),
			Type:                    row.Type,
			Name:                    row.Name,
			DurationSeconds:         row.DurationSeconds,
			EstimatedCaloriesBurned: row.EstimatedCaloriesBurned,
			Steps:                   row.Steps,
			DistanceMeters:          row.DistanceMeters,
			ElevationGainMeters:     row.ElevationGainMeters,
			AverageHeartRate:        row.AverageHeartRate,
			MaxHeartRate:            row.MaxHeartRate,
			TotalVolumeKg:           row.TotalVolumeKg,
			TotalSets:               row.TotalSets,
			Notes:                   row.Notes,
			CreatedAt:               row.CreatedAt,
		}
		response.Exercises = append(response.Exercises, item)
	}

	// meals
	var mealRows []struct {
		Id           uuid.UUID `gorm:"column:id"`
		Date         time.Time `gorm:"column:date"`
		Type         string    `gorm:"column:type"`
		Name         string    `gorm:"column:name"`
		Calories     *float64  `gorm:"column:calories"`
		ProteinGrams *float64  `gorm:"column:protein_grams"`
		CarbsGrams   *float64  `gorm:"column:carbs_grams"`
		FatGrams     *float64  `gorm:"column:fat_grams"`
		FiberGrams   *float64  `gorm:"column:fiber_grams"`
		Notes        string    `gorm:"column:notes"`
		CreatedAt    time.Time `gorm:"column:created_at"`
	}
	err = a.db.Raw(`
		SELECT id, date, type, name, calories, protein_grams, carbs_grams, fat_grams, fiber_grams, notes, created_at
		FROM meals WHERE user_id = ? ORDER BY date
	`, userId).Scan(&mealRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range mealRows {
		item := mealExport{
			Id:           row.Id.String(),
			Date:         row.Date.Format("2006-01-02"),
			Type:         row.Type,
			Name:         row.Name,
			Calories:     row.Calories,
			ProteinGrams: row.ProteinGrams,
			CarbsGrams:   row.CarbsGrams,
			FatGrams:     row.FatGrams,
			FiberGrams:   row.FiberGrams,
			Notes:        row.Notes,
			CreatedAt:    row.CreatedAt,
		}
		response.Meals = append(response.Meals, item)
	}

	// foods (user's own)
	var foodRows []struct {
		Id                  uuid.UUID `gorm:"column:id"`
		Name                string    `gorm:"column:name"`
		MeasurementType     string    `gorm:"column:measurement_type"`
		BaseUnit            string    `gorm:"column:base_unit"`
		BaseQuantity        *float64  `gorm:"column:base_quantity"`
		DefaultCalories     *float64  `gorm:"column:default_calories"`
		DefaultProteinGrams *float64  `gorm:"column:default_protein_grams"`
		DefaultCarbsGrams   *float64  `gorm:"column:default_carbs_grams"`
		DefaultFatGrams     *float64  `gorm:"column:default_fat_grams"`
		DefaultFiberGrams   *float64  `gorm:"column:default_fiber_grams"`
		Public              bool      `gorm:"column:public"`
		CreatedAt           time.Time `gorm:"column:created_at"`
	}
	err = a.db.Raw(`
		SELECT id, name, measurement_type, base_unit, base_quantity,
		       default_calories, default_protein_grams, default_carbs_grams,
		       default_fat_grams, default_fiber_grams, public, created_at
		FROM foods WHERE user_id = ? ORDER BY name
	`, userId).Scan(&foodRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range foodRows {
		item := foodExport{
			Id:                  row.Id.String(),
			Name:                row.Name,
			MeasurementType:     row.MeasurementType,
			BaseUnit:            row.BaseUnit,
			BaseQuantity:        row.BaseQuantity,
			DefaultCalories:     row.DefaultCalories,
			DefaultProteinGrams: row.DefaultProteinGrams,
			DefaultCarbsGrams:   row.DefaultCarbsGrams,
			DefaultFatGrams:     row.DefaultFatGrams,
			DefaultFiberGrams:   row.DefaultFiberGrams,
			Public:              row.Public,
			CreatedAt:           row.CreatedAt,
		}
		response.Foods = append(response.Foods, item)
	}

	// body measurements
	var measureRows []struct {
		Date      time.Time `gorm:"column:date"`
		Type      string    `gorm:"column:type"`
		Value     float64   `gorm:"column:value"`
		Notes     string    `gorm:"column:notes"`
		UpdatedAt time.Time `gorm:"column:updated_at"`
	}
	err = a.db.Raw(
		`SELECT date, type, value, notes, updated_at FROM body_measurements WHERE user_id = ? ORDER BY date, type`,
		userId,
	).Scan(&measureRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range measureRows {
		item := bodyMeasurementExport{
			Date:      row.Date.Format("2006-01-02"),
			Type:      row.Type,
			Value:     row.Value,
			Notes:     row.Notes,
			UpdatedAt: row.UpdatedAt,
		}
		response.BodyMeasurements = append(response.BodyMeasurements, item)
	}

	// goal
	var goalRow struct {
		DailyCalories        *float64  `gorm:"column:daily_calories"`
		DailyProteinGrams    *float64  `gorm:"column:daily_protein_grams"`
		DailyCarbsGrams      *float64  `gorm:"column:daily_carbs_grams"`
		DailyFatGrams        *float64  `gorm:"column:daily_fat_grams"`
		DailyFiberGrams      *float64  `gorm:"column:daily_fiber_grams"`
		DailySteps           *int      `gorm:"column:daily_steps"`
		DailyExerciseMinutes *int      `gorm:"column:daily_exercise_minutes"`
		TargetWeightKg       *float64  `gorm:"column:target_weight_kg"`
		StartedAt            time.Time `gorm:"column:started_at"`
	}
	err = a.db.Raw(
		`SELECT daily_calories, daily_protein_grams, daily_carbs_grams, daily_fat_grams,
		        daily_fiber_grams, daily_steps, daily_exercise_minutes, target_weight_kg, started_at
		 FROM goals WHERE user_id = ?`,
		userId,
	).Scan(&goalRow).Error
	if err != nil {
		return nil, err
	}
	if goalRow.StartedAt != (time.Time{}) {
		response.Goal = &goalExport{
			DailyCalories:        goalRow.DailyCalories,
			DailyProteinGrams:    goalRow.DailyProteinGrams,
			DailyCarbsGrams:      goalRow.DailyCarbsGrams,
			DailyFatGrams:        goalRow.DailyFatGrams,
			DailyFiberGrams:      goalRow.DailyFiberGrams,
			DailySteps:           goalRow.DailySteps,
			DailyExerciseMinutes: goalRow.DailyExerciseMinutes,
			TargetWeightKg:       goalRow.TargetWeightKg,
			StartedAt:            goalRow.StartedAt,
		}
	}

	return response, nil
}
