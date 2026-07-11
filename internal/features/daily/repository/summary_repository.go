package repository

import (
	"errors"
	"math"
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

	profile, err := r.getUserProfile(userId)
	if err != nil {
		return nil, err
	}

	closed, err := r.isDateClosed(userId, date)
	if err != nil {
		return nil, err
	}

	steps, err := r.getSteps(userId, date)
	if err != nil {
		return nil, err
	}

	var weightKg *float64
	if weightEntry != nil {
		weightKg = &weightEntry.WeightKg
	} else {
		weightKg = r.getLatestWeightKg(userId, date)
	}

	bmr := calculateBMR(profile, weightEntry, date)

	summary := &domain.DailySummary{
		Date:                date,
		Closed:              closed,
		MealsSummary:        *mealsSummary,
		ExerciseSummary:     *exerciseSummary,
		WeightEntry:         weightEntry,
		Goals:               goals,
		EstimatedBMR:        bmr,
		StepsCaloriesBurned: calculateStepCalories(steps, weightKg),
	}
	applyCaloricBalance(summary)
	return summary, nil
}

func (r *summaryRepository) GetDailySummaryRange(userId uuid.UUID, from, to time.Time) ([]domain.DailySummary, error) {
	mealsMap, err := r.getMealsSummaryRange(userId, from, to)
	if err != nil {
		return nil, err
	}

	exerciseMap, err := r.getExerciseSummaryRange(userId, from, to)
	if err != nil {
		return nil, err
	}

	weightMap, err := r.getWeightEntriesRange(userId, from, to)
	if err != nil {
		return nil, err
	}

	goals, err := r.getGoals(userId)
	if err != nil {
		return nil, err
	}

	profile, err := r.getUserProfile(userId)
	if err != nil {
		return nil, err
	}

	closedDates, err := r.getClosedDatesRange(userId, from, to)
	if err != nil {
		return nil, err
	}

	stepsMap, err := r.getStepsRange(userId, from, to)
	if err != nil {
		return nil, err
	}
	baselineWeight := r.getLatestWeightKg(userId, from)

	var summaries []domain.DailySummary
	for d := from; !d.After(to); d = d.AddDate(0, 0, 1) {
		key := d.Format("2006-01-02")
		summary := domain.DailySummary{
			Date:   d,
			Closed: closedDates[key],
			Goals:  goals,
		}
		m, ok := mealsMap[key]
		if ok {
			summary.MealsSummary = m
		}
		e, ok := exerciseMap[key]
		if ok {
			summary.ExerciseSummary = e
		}
		var weightEntry *domain.WeightEntrySummary
		var weightKg *float64
		w, ok := weightMap[key]
		if ok {
			summary.WeightEntry = &w
			weightEntry = &w
			weightKg = &w.WeightKg
		} else {
			weightKg = baselineWeight
		}
		summary.EstimatedBMR = calculateBMR(profile, weightEntry, d)
		summary.StepsCaloriesBurned = calculateStepCalories(stepsMap[key], weightKg)
		applyCaloricBalance(&summary)
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (r *summaryRepository) getMealsSummaryRange(userId uuid.UUID, from, to time.Time) (map[string]domain.MealsSummary, error) {
	var results []struct {
		Date              string
		TotalCalories     float64
		TotalProteinGrams float64
		TotalCarbsGrams   float64
		TotalFatGrams     float64
		TotalFiberGrams   float64
		Count             int
	}
	err := r.db.
		Table("meals").
		Select(`
			date::text as date,
			COALESCE(SUM(calories), 0) as total_calories,
			COALESCE(SUM(protein_grams), 0) as total_protein_grams,
			COALESCE(SUM(carbs_grams), 0) as total_carbs_grams,
			COALESCE(SUM(fat_grams), 0) as total_fat_grams,
			COALESCE(SUM(fiber_grams), 0) as total_fiber_grams,
			COUNT(*) as count
		`).
		Where("user_id = ? AND date >= ? AND date <= ?", userId, from, to).
		Group("date").
		Scan(&results).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("aggregating meals summary range", err)
	}
	m := make(map[string]domain.MealsSummary, len(results))
	for _, r := range results {
		m[r.Date] = domain.MealsSummary{
			TotalCalories:     r.TotalCalories,
			TotalProteinGrams: r.TotalProteinGrams,
			TotalCarbsGrams:   r.TotalCarbsGrams,
			TotalFatGrams:     r.TotalFatGrams,
			TotalFiberGrams:   r.TotalFiberGrams,
			Count:             r.Count,
		}
	}
	return m, nil
}

func (r *summaryRepository) getExerciseSummaryRange(userId uuid.UUID, from, to time.Time) (map[string]domain.ExerciseSummary, error) {
	var results []struct {
		Date                 string
		TotalCaloriesBurned  float64
		TotalSteps           int
		TotalDurationSeconds int
		TotalDistanceMeters  float64
		Count                int
	}
	err := r.db.
		Table("exercises").
		Select(`
			date::text as date,
			COALESCE(SUM(estimated_calories_burned), 0) as total_calories_burned,
			COALESCE(SUM(steps), 0) as total_steps,
			COALESCE(SUM(duration_seconds), 0) as total_duration_seconds,
			COALESCE(SUM(distance_meters), 0) as total_distance_meters,
			COUNT(*) as count
		`).
		Where("user_id = ? AND date >= ? AND date <= ?", userId, from, to).
		Group("date").
		Scan(&results).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("aggregating exercise summary range", err)
	}
	m := make(map[string]domain.ExerciseSummary, len(results))
	for _, r := range results {
		m[r.Date] = domain.ExerciseSummary{
			TotalCaloriesBurned:  r.TotalCaloriesBurned,
			TotalSteps:           r.TotalSteps,
			TotalDurationSeconds: r.TotalDurationSeconds,
			TotalDistanceMeters:  r.TotalDistanceMeters,
			Count:                r.Count,
		}
	}
	return m, nil
}

func (r *summaryRepository) getWeightEntriesRange(userId uuid.UUID, from, to time.Time) (map[string]domain.WeightEntrySummary, error) {
	var results []struct {
		Date              string
		WeightKg          float64
		BodyFatPercentage *float64
	}
	err := r.db.
		Table("weight_entries").
		Select("date::text as date, weight_kg, body_fat_percentage").
		Where("user_id = ? AND date >= ? AND date <= ?", userId, from, to).
		Scan(&results).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching weight entries range", err)
	}
	m := make(map[string]domain.WeightEntrySummary, len(results))
	for _, r := range results {
		m[r.Date] = domain.WeightEntrySummary{
			WeightKg:          r.WeightKg,
			BodyFatPercentage: r.BodyFatPercentage,
		}
	}
	return m, nil
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

type userProfile struct {
	HeightCm  *int
	BirthDate *time.Time
	Sex       *string
}

func (r *summaryRepository) getUserProfile(userId uuid.UUID) (*userProfile, error) {
	var profile userProfile
	err := r.db.
		Table("users").
		Select("height_cm, birth_date, sex").
		Where("id = ?", userId).
		Take(&profile).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching user profile for BMR", err)
	}
	return &profile, nil
}

func calculateBMR(profile *userProfile, weight *domain.WeightEntrySummary, date time.Time) *float64 {
	if weight == nil {
		return nil
	}

	// Katch-McArdle: 370 + 21.6 × lean_mass_kg
	if weight.BodyFatPercentage != nil {
		leanMass := weight.WeightKg * (1 - *weight.BodyFatPercentage/100)
		bmr := 370 + 21.6*leanMass
		bmr = math.Round(bmr*100) / 100
		return &bmr
	}

	// Mifflin-St Jeor fallback: 10×weight + 6.25×height - 5×age + offset
	if profile.HeightCm != nil && profile.BirthDate != nil && profile.Sex != nil {
		age := date.Year() - profile.BirthDate.Year()
		if date.YearDay() < profile.BirthDate.YearDay() {
			age--
		}
		bmr := 10*weight.WeightKg + 6.25*float64(*profile.HeightCm) - 5*float64(age)
		if *profile.Sex == "male" {
			bmr += 5
		} else {
			bmr -= 161
		}
		bmr = math.Round(bmr*100) / 100
		return &bmr
	}

	return nil
}

func applyCaloricBalance(summary *domain.DailySummary) {
	if summary.EstimatedBMR == nil {
		return
	}
	burned := *summary.EstimatedBMR + summary.ExerciseSummary.TotalCaloriesBurned
	if summary.StepsCaloriesBurned != nil {
		burned += *summary.StepsCaloriesBurned
	}
	balance := math.Round((summary.MealsSummary.TotalCalories-burned)*100) / 100
	summary.CaloricBalance = &balance
}

func (r *summaryRepository) getSteps(userId uuid.UUID, date time.Time) (int, error) {
	var steps int
	err := r.db.Raw(`
		SELECT COALESCE(SUM(steps), 0) FROM exercises WHERE user_id = ? AND date = ?
	`, userId, date).Scan(&steps).Error
	if err != nil {
		return 0, cerr.NewInternalError("fetching daily steps from exercises for summary", err)
	}
	return steps, nil
}

func (r *summaryRepository) getStepsRange(userId uuid.UUID, from, to time.Time) (map[string]int, error) {
	var results []struct {
		Date  string
		Steps int
	}
	err := r.db.Raw(`
		SELECT date::text as date, COALESCE(SUM(steps), 0) as steps
		FROM exercises
		WHERE user_id = ? AND date >= ? AND date <= ?
		GROUP BY date
	`, userId, from, to).Scan(&results).Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching daily steps range from exercises for summary", err)
	}
	m := make(map[string]int, len(results))
	for _, row := range results {
		m[row.Date] = row.Steps
	}
	return m, nil
}

func (r *summaryRepository) getLatestWeightKg(userId uuid.UUID, upToDate time.Time) *float64 {
	var result struct {
		WeightKg float64
	}
	err := r.db.Raw(`
		SELECT weight_kg FROM weight_entries
		WHERE user_id = ? AND date <= ?
		ORDER BY date DESC LIMIT 1
	`, userId, upToDate).Scan(&result).Error
	if err != nil || result.WeightKg == 0 {
		return nil
	}
	return &result.WeightKg
}

func calculateStepCalories(steps int, weightKg *float64) *float64 {
	if steps == 0 || weightKg == nil {
		return nil
	}
	cal := math.Round(float64(steps)*0.0005**weightKg*100) / 100
	return &cal
}

func (r *summaryRepository) GetDailyCheck(userId uuid.UUID, date time.Time) (*domain.DailyCheck, error) {
	var missingMeasurements int64
	err := r.db.Raw(`
		SELECT COUNT(mi.id)
		FROM meal_items mi
		JOIN meals m ON m.id = mi.meal_id
		WHERE m.user_id = ? AND m.date = ?
		  AND (mi.measurement_method IS NULL OR mi.measurement_method = '')
	`, userId, date).Scan(&missingMeasurements).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking missing measurements", err)
	}

	var mealsWithoutPhoto int64
	err = r.db.Raw(`
		SELECT COUNT(m.id)
		FROM meals m
		WHERE m.user_id = ? AND m.date = ?
		  AND NOT EXISTS (SELECT 1 FROM meal_photos mp WHERE mp.meal_id = m.id)
	`, userId, date).Scan(&mealsWithoutPhoto).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking meals without photo", err)
	}

	var dailyPhotoCount int64
	err = r.db.Raw(`
		SELECT COUNT(*) FROM daily_photos WHERE user_id = ? AND date = ?
	`, userId, date).Scan(&dailyPhotoCount).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking daily photo", err)
	}

	var stepsCount int64
	err = r.db.Raw(`
		SELECT COUNT(*) FROM exercises WHERE user_id = ? AND date = ? AND steps IS NOT NULL AND steps > 0
	`, userId, date).Scan(&stepsCount).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking daily steps", err)
	}

	var exerciseCount int64
	err = r.db.Raw(`
		SELECT COUNT(*) FROM exercises WHERE user_id = ? AND date = ?
	`, userId, date).Scan(&exerciseCount).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking exercises", err)
	}

	weekAgo := date.AddDate(0, 0, -6)
	var recentWeightCount int64
	err = r.db.Raw(`
		SELECT COUNT(*) FROM weight_entries
		WHERE user_id = ? AND date >= ? AND date <= ?
	`, userId, weekAgo, date).Scan(&recentWeightCount).Error
	if err != nil {
		return nil, cerr.NewInternalError("checking recent weight", err)
	}

	check := &domain.DailyCheck{
		Date:                date,
		MissingMeasurements: int(missingMeasurements),
		MealsWithoutPhoto:   int(mealsWithoutPhoto),
		HasDailyPhoto:       dailyPhotoCount > 0,
		HasSteps:            stepsCount > 0,
		HasExercise:         exerciseCount > 0,
		HasRecentWeight:     recentWeightCount > 0,
	}
	check.Complete = check.MissingMeasurements == 0 &&
		check.MealsWithoutPhoto == 0 &&
		check.HasDailyPhoto &&
		check.HasSteps &&
		check.HasRecentWeight
	return check, nil
}

func (r *summaryRepository) isDateClosed(userId uuid.UUID, date time.Time) (bool, error) {
	var count int64
	err := r.db.
		Table("day_closures").
		Where("user_id = ? AND date = ?", userId, date).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking day closure for summary", err)
	}
	return count > 0, nil
}

func (r *summaryRepository) getClosedDatesRange(userId uuid.UUID, from, to time.Time) (map[string]bool, error) {
	var results []struct {
		Date string
	}
	err := r.db.
		Table("day_closures").
		Select("date::text as date").
		Where("user_id = ? AND date >= ? AND date <= ?", userId, from, to).
		Scan(&results).
		Error
	if err != nil {
		return nil, cerr.NewInternalError("fetching closed dates for summary", err)
	}
	m := make(map[string]bool, len(results))
	for _, r := range results {
		m[r.Date] = true
	}
	return m, nil
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
