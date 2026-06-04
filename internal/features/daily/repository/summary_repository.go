package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/daily/domain"
	"github.com/ivan-ca97/life/internal/features/daily/ports"
)

type summaryRepository struct {
	db *gorm.DB
}

var _ ports.SummaryRepository = (*summaryRepository)(nil)

func NewSummaryRepository(db *gorm.DB) *summaryRepository {
	return &summaryRepository{
		db: db,
	}
}

func (r *summaryRepository) GetDailySummary(userId uuid.UUID, date time.Time) (*domain.DailySummary, error) {
	mealsSummary, err := r.getMealsSummary(userId, date)
	if err != nil {
		return nil, err
	}

	exerciseSummary, err := r.getExerciseSummary(userId, date)
	if err != nil {
		return nil, err
	}

	weightEntry, err := r.getWeightEntry(userId, date)
	if err != nil {
		return nil, err
	}

	goals, err := r.getGoals(userId)
	if err != nil {
		return nil, err
	}

	summary := &domain.DailySummary{
		Date:            date,
		MealsSummary:    *mealsSummary,
		ExerciseSummary: *exerciseSummary,
		WeightEntry:     weightEntry,
		Goals:           goals,
	}
	return summary, nil
}

func (r *summaryRepository) getMealsSummary(userId uuid.UUID, date time.Time) (*domain.MealsSummary, error) {
	var result domain.MealsSummary
	err := r.db.
		Table("meals").
		Select(`
			COALESCE(SUM(calories), 0) as total_calories,
			COALESCE(SUM(protein_grams), 0) as total_protein_grams,
			COALESCE(SUM(carbs_grams), 0) as total_carbs_grams,
			COALESCE(SUM(fat_grams), 0) as total_fat_grams,
			COALESCE(SUM(fiber_grams), 0) as total_fiber_grams,
			COUNT(*) as count
		`).
		Where("user_id = ? AND date = ?", userId, date).
		Scan(&result).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("aggregating meals summary", err)
	}
	return &result, nil
}

func (r *summaryRepository) getExerciseSummary(userId uuid.UUID, date time.Time) (*domain.ExerciseSummary, error) {
	var result domain.ExerciseSummary
	err := r.db.
		Table("exercises").
		Select(`
			COALESCE(SUM(estimated_calories_burned), 0) as total_calories_burned,
			COALESCE(SUM(steps), 0) as total_steps,
			COALESCE(SUM(duration_seconds), 0) as total_duration_seconds,
			COALESCE(SUM(distance_meters), 0) as total_distance_meters,
			COUNT(*) as count
		`).
		Where("user_id = ? AND date = ?", userId, date).
		Scan(&result).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("aggregating exercise summary", err)
	}
	return &result, nil
}

func (r *summaryRepository) getWeightEntry(userId uuid.UUID, date time.Time) (*domain.WeightEntrySummary, error) {
	var result struct {
		WeightKg          float64
		BodyFatPercentage *float64
	}
	err := r.db.
		Table("weight_entries").
		Select("weight_kg, body_fat_percentage").
		Where("user_id = ? AND date = ?", userId, date).
		Take(&result).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("fetching weight entry for summary", err)
	}
	return &domain.WeightEntrySummary{
		WeightKg:          result.WeightKg,
		BodyFatPercentage: result.BodyFatPercentage,
	}, nil
}

func (r *summaryRepository) getGoals(userId uuid.UUID) (*domain.GoalsSummary, error) {
	var result struct {
		DailyCalories        *float64
		DailyProteinGrams    *float64
		DailyCarbsGrams      *float64
		DailyFatGrams        *float64
		DailyFiberGrams      *float64
		DailySteps           *int
		DailyExerciseMinutes *int
		TargetWeightKg       *float64
	}
	err := r.db.
		Table("goals").
		Select("daily_calories, daily_protein_grams, daily_carbs_grams, daily_fat_grams, daily_fiber_grams, daily_steps, daily_exercise_minutes, target_weight_kg").
		Where("user_id = ?", userId).
		Take(&result).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("fetching goals for summary", err)
	}
	return &domain.GoalsSummary{
		DailyCalories:        result.DailyCalories,
		DailyProteinGrams:    result.DailyProteinGrams,
		DailyCarbsGrams:      result.DailyCarbsGrams,
		DailyFatGrams:        result.DailyFatGrams,
		DailyFiberGrams:      result.DailyFiberGrams,
		DailySteps:           result.DailySteps,
		DailyExerciseMinutes: result.DailyExerciseMinutes,
		TargetWeightKg:       result.TargetWeightKg,
	}, nil
}
