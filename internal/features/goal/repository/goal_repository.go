package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/features/goal/domain"
	"github.com/ivan-ca97/life/internal/features/goal/ports"
)

type goalRepository struct {
	db *gorm.DB
}

var _ ports.GoalRepository = (*goalRepository)(nil)

func NewGoalRepository(db *gorm.DB) *goalRepository {
	return &goalRepository{
		db: db,
	}
}

func (r *goalRepository) FindByUserId(userId uuid.UUID) (*domain.Goal, error) {
	var model goal
	err := r.db.
		Where("user_id = ?", userId).
		First(&model).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrGoalNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding goal by user id", err)
	}
	return model.toDomain(), nil
}

func (r *goalRepository) Upsert(g *domain.Goal) (*domain.Goal, error) {
	model := goalFromDomain(g)
	err := r.db.Save(model).Error
	if err != nil {
		return nil, cerr.NewInternalError("upserting goal", err)
	}
	return model.toDomain(), nil
}

func (r *goalRepository) GetProgress(userId uuid.UUID, from, to time.Time) (*domain.GoalProgress, error) {
	goal, err := r.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	type mealRow struct {
		Date         time.Time `gorm:"column:date"`
		Calories     float64   `gorm:"column:calories"`
		ProteinGrams float64   `gorm:"column:protein_grams"`
		CarbsGrams   float64   `gorm:"column:carbs_grams"`
		FatGrams     float64   `gorm:"column:fat_grams"`
		FiberGrams   float64   `gorm:"column:fiber_grams"`
	}
	var mealRows []mealRow
	if err := r.db.Raw(`
		SELECT date::date AS date,
			SUM(COALESCE(calories, 0))      AS calories,
			SUM(COALESCE(protein_grams, 0)) AS protein_grams,
			SUM(COALESCE(carbs_grams, 0))   AS carbs_grams,
			SUM(COALESCE(fat_grams, 0))     AS fat_grams,
			SUM(COALESCE(fiber_grams, 0))   AS fiber_grams
		FROM meals
		WHERE user_id = ? AND date::date >= ? AND date::date <= ?
		GROUP BY date::date
	`, userId, fromStr, toStr).Scan(&mealRows).Error; err != nil {
		return nil, cerr.NewInternalError("querying meal progress", err)
	}

	type stepsRow struct {
		Date  string `gorm:"column:date"`
		Steps int    `gorm:"column:steps"`
	}
	var stepsRows []stepsRow
	if err := r.db.Raw(`
		SELECT date::text AS date, COALESCE(SUM(steps), 0) AS steps
		FROM exercises
		WHERE user_id = ? AND date >= ? AND date <= ?
		GROUP BY date
	`, userId, fromStr, toStr).Scan(&stepsRows).Error; err != nil {
		return nil, cerr.NewInternalError("querying steps progress", err)
	}

	type exerciseRow struct {
		ExerciseMinutes float64 `gorm:"column:exercise_minutes"`
	}
	var exerciseRows []exerciseRow
	if err := r.db.Raw(`
		SELECT SUM(COALESCE(duration_seconds, 0)) / 60.0 AS exercise_minutes
		FROM exercises
		WHERE user_id = ? AND date::date >= ? AND date::date <= ?
		GROUP BY date::date
	`, userId, fromStr, toStr).Scan(&exerciseRows).Error; err != nil {
		return nil, cerr.NewInternalError("querying exercise progress", err)
	}

	var wRow struct {
		WeightKg float64 `gorm:"column:weight_kg"`
	}
	if err := r.db.Raw(`
		SELECT weight_kg FROM weight_entries WHERE user_id = ? ORDER BY date DESC LIMIT 1
	`, userId).Scan(&wRow).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, cerr.NewInternalError("querying current weight", err)
	}
	var currentWeight *float64
	if wRow.WeightKg > 0 {
		w := wRow.WeightKg
		currentWeight = &w
	}

	daysTotal := int(to.Sub(from).Hours()/24) + 1
	progress := &domain.GoalProgress{
		From:      from,
		To:        to,
		Goal:      goal,
		DaysTotal: daysTotal,
	}

	// Meal-based metrics (single pass over mealRows)
	if goal.DailyCalories != nil || goal.DailyProteinGrams != nil ||
		goal.DailyCarbsGrams != nil || goal.DailyFatGrams != nil || goal.DailyFiberGrams != nil {
		var calSum, protSum, carbSum, fatSum, fiberSum float64
		var calMet, protMet, carbMet, fatMet, fiberMet int
		for _, row := range mealRows {
			calSum += row.Calories
			protSum += row.ProteinGrams
			carbSum += row.CarbsGrams
			fatSum += row.FatGrams
			fiberSum += row.FiberGrams
			if goal.DailyCalories != nil && row.Calories >= *goal.DailyCalories {
				calMet++
			}
			if goal.DailyProteinGrams != nil && row.ProteinGrams >= *goal.DailyProteinGrams {
				protMet++
			}
			if goal.DailyCarbsGrams != nil && row.CarbsGrams >= *goal.DailyCarbsGrams {
				carbMet++
			}
			if goal.DailyFatGrams != nil && row.FatGrams >= *goal.DailyFatGrams {
				fatMet++
			}
			if goal.DailyFiberGrams != nil && row.FiberGrams >= *goal.DailyFiberGrams {
				fiberMet++
			}
		}
		tracked := len(mealRows)
		if goal.DailyCalories != nil {
			progress.DailyCalories = buildMetric(*goal.DailyCalories, calSum, calMet, tracked, daysTotal)
		}
		if goal.DailyProteinGrams != nil {
			progress.DailyProteinGrams = buildMetric(*goal.DailyProteinGrams, protSum, protMet, tracked, daysTotal)
		}
		if goal.DailyCarbsGrams != nil {
			progress.DailyCarbsGrams = buildMetric(*goal.DailyCarbsGrams, carbSum, carbMet, tracked, daysTotal)
		}
		if goal.DailyFatGrams != nil {
			progress.DailyFatGrams = buildMetric(*goal.DailyFatGrams, fatSum, fatMet, tracked, daysTotal)
		}
		if goal.DailyFiberGrams != nil {
			progress.DailyFiberGrams = buildMetric(*goal.DailyFiberGrams, fiberSum, fiberMet, tracked, daysTotal)
		}
	}

	if goal.DailySteps != nil {
		var sum float64
		var met int
		for _, row := range stepsRows {
			sum += float64(row.Steps)
			if row.Steps >= *goal.DailySteps {
				met++
			}
		}
		progress.DailySteps = buildMetric(float64(*goal.DailySteps), sum, met, len(stepsRows), daysTotal)
	}

	if goal.DailyExerciseMinutes != nil {
		var sum float64
		var met int
		for _, row := range exerciseRows {
			sum += row.ExerciseMinutes
			if row.ExerciseMinutes >= float64(*goal.DailyExerciseMinutes) {
				met++
			}
		}
		progress.DailyExerciseMinutes = buildMetric(float64(*goal.DailyExerciseMinutes), sum, met, len(exerciseRows), daysTotal)
	}

	if goal.TargetWeightKg != nil {
		progress.WeightProgress = &domain.WeightProgress{
			TargetKg:  *goal.TargetWeightKg,
			CurrentKg: currentWeight,
		}
	}

	return progress, nil
}

func buildMetric(target, sum float64, met, tracked, total int) *domain.GoalMetric {
	avg := 0.0
	if tracked > 0 {
		avg = sum / float64(tracked)
	}
	return &domain.GoalMetric{
		Target:      target,
		Average:     avg,
		DaysMet:     met,
		DaysTracked: tracked,
		DaysTotal:   total,
	}
}
