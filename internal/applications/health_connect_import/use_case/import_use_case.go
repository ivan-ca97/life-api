package use_case

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
	exerciseDomain "github.com/ivan-ca97/life/internal/features/exercise/domain"
	exercisePorts "github.com/ivan-ca97/life/internal/features/exercise/ports"
	weightDomain "github.com/ivan-ca97/life/internal/features/weight/domain"
	weightPorts "github.com/ivan-ca97/life/internal/features/weight/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

const (
	source            = "hc"
	dailyWalkMinSteps = 200 // minimum residual steps to create a daily-walk exercise
	// dailyWalkName is the exercise name for HC-imported daily walks. It doubles
	// as the lookup key for re-sync dedup, so it must match existing rows exactly
	// (see migration 000035, which renamed legacy "Caminata cotidiana" rows).
	dailyWalkName = "Daily walk"
)

type healthConnectImportUseCase struct {
	weightService      weightPorts.WeightEntryService
	weightRepository   weightPorts.WeightEntryRepository
	exerciseService    exercisePorts.ExerciseService
	exerciseRepository exercisePorts.ExerciseRepository
	rawStore           ports.RawRecordStore
	syncLogs           ports.SyncLogStore
	authorizer         auth.AuthorizationService
}

var _ ports.HealthConnectImportUseCase = (*healthConnectImportUseCase)(nil)

func NewHealthConnectImportUseCase(
	weightService weightPorts.WeightEntryService,
	weightRepository weightPorts.WeightEntryRepository,
	exerciseService exercisePorts.ExerciseService,
	exerciseRepository exercisePorts.ExerciseRepository,
	rawStore ports.RawRecordStore,
	syncLogs ports.SyncLogStore,
	authorizer auth.AuthorizationService,
) *healthConnectImportUseCase {
	return &healthConnectImportUseCase{
		weightService:      weightService,
		weightRepository:   weightRepository,
		exerciseService:    exerciseService,
		exerciseRepository: exerciseRepository,
		rawStore:           rawStore,
		syncLogs:           syncLogs,
		authorizer:         authorizer,
	}
}

func (uc *healthConnectImportUseCase) Import(ctx context.Context, userId uuid.UUID, payload *ports.Payload) (*ports.ImportResult, error) {
	result := &ports.ImportResult{}

	if len(payload.Weight) > 0 {
		err := uc.authorizer.Authorize(ctx, userId, permissions.WeightCreate)
		if err != nil {
			return nil, err
		}
		err = uc.importWeight(userId, payload.Weight, &result.Weight)
		if err != nil {
			return nil, err
		}
	}

	hasActivity := len(payload.ExerciseSessions)+len(payload.StepsDaily)+len(payload.Sleep)+len(payload.HeartRate) > 0
	if hasActivity {
		err := uc.authorizer.Authorize(ctx, userId, permissions.ExercisesCreate)
		if err != nil {
			return nil, err
		}
	}

	err := uc.importActivity(userId, payload.ExerciseSessions, payload.StepsDaily, payload.HeartRate, &result.Exercise)
	if err != nil {
		return nil, err
	}
	err = uc.importSleep(userId, payload.Sleep, &result.Sleep)
	if err != nil {
		return nil, err
	}
	err = uc.importHeartRate(userId, payload.HeartRate, &result.HeartRate)
	if err != nil {
		return nil, err
	}

	uc.writeSyncLog(userId, payload, result)

	return result, nil
}

func (uc *healthConnectImportUseCase) writeSyncLog(userId uuid.UUID, payload *ports.Payload, result *ports.ImportResult) {
	syncedAt, err := parseInstant(payload.SyncedAt)
	if err != nil {
		syncedAt = time.Now().UTC()
	}
	summary := map[string]any{
		"weight":     map[string]int{"created": result.Weight.Created, "skipped": result.Weight.Skipped, "blocked": result.Weight.Blocked},
		"exercise":   map[string]int{"created": result.Exercise.Created, "skipped": result.Exercise.Skipped, "blocked": result.Exercise.Blocked},
		"sleep":      map[string]int{"created": result.Sleep.Created, "skipped": result.Sleep.Skipped, "blocked": result.Sleep.Blocked},
		"heart_rate": map[string]int{"created": result.HeartRate.Created, "skipped": result.HeartRate.Skipped, "blocked": result.HeartRate.Blocked},
	}
	resultBytes, err := json.Marshal(summary)
	if err != nil {
		return
	}
	syncLog := &ports.SyncLog{
		UserId:     userId,
		AppVersion: payload.AppVersion,
		SyncedAt:   syncedAt,
		Result:     resultBytes,
	}
	_ = uc.syncLogs.Create(syncLog)
}

// ─── Weight ───────────────────────────────────────────────────────────────────

func (uc *healthConnectImportUseCase) importWeight(userId uuid.UUID, records []ports.WeightRecord, out *ports.TypeResult) error {
	for _, rec := range records {
		externalId := externalIdFor("weight", rec.Id)
		if externalId != nil {
			exists, err := uc.weightRepository.ExistsByExternalId(userId, *externalId)
			if err != nil {
				return err
			}
			if exists {
				out.Skipped++
				continue
			}
		}

		t, err := parseInstant(rec.Time)
		if err != nil {
			out.Skipped++
			continue
		}

		createParams := weightPorts.CreateParams{
			Date:       dateOf(t),
			WeightKg:   rec.Kilograms,
			ExternalId: externalId,
		}
		_, err = uc.weightService.Create(userId, createParams)
		switch {
		case errors.Is(err, dayclosure.ErrDayClosed):
			out.Blocked++
		case errors.Is(err, weightDomain.ErrWeightEntryConflict):
			out.Skipped++
		case err != nil:
			return err
		default:
			out.Created++
		}
	}
	return nil
}

// ─── Activity (exercises + daily steps) ──────────────────────────────────────

func (uc *healthConnectImportUseCase) importActivity(
	userId uuid.UUID,
	sessions []ports.ExerciseRecord,
	stepsDaily []ports.DailyStepsRecord,
	hrRecords []ports.HeartRateRecord,
	exOut *ports.TypeResult,
) error {
	// 1. Import explicit exercise sessions from Health Connect.
	for _, session := range sessions {
		start, err := parseInstant(session.StartTime)
		if err != nil {
			exOut.Skipped++
			continue
		}
		end, err := parseInstant(session.EndTime)
		if err != nil {
			exOut.Skipped++
			continue
		}

		externalId := externalIdFor("exercise", session.Id)
		if externalId != nil {
			exists, err := uc.exerciseRepository.ExistsByExternalId(userId, *externalId)
			if err != nil {
				return err
			}
			if exists {
				exOut.Skipped++
				continue
			}
		}

		avgHR, maxHR := computeHR(hrRecords, start, end)

		exerciseType := mapExerciseType(session.Type)
		name := strings.TrimSpace(session.Title)
		if name == "" {
			name = exerciseType
		}

		importSource := exerciseDomain.ImportSourceHealthConnect
		// Steps are only meaningful for walking/running; cycling and weightlifting
		// sessions can carry incidental step counts from HC that would corrupt the
		// daily total. Guard here so old app versions don't pollute the count.
		var sessionSteps *int
		if exerciseType == exerciseDomain.ExerciseTypeWalking || exerciseType == exerciseDomain.ExerciseTypeRunning {
			sessionSteps = session.Steps
		}
		params := exercisePorts.CreateParams{
			Date:             dateOf(start),
			Type:             exerciseType,
			Name:             name,
			StartedAt:        &start,
			DurationSeconds:  session.DurationSeconds,
			DistanceMeters:   session.DistanceMeters,
			Steps:            sessionSteps,
			AverageHeartRate: avgHR,
			MaxHeartRate:     maxHR,
			ExternalId:       externalId,
			ImportSource:     &importSource,
		}
		if params.DurationSeconds == nil {
			d := int(end.Sub(start).Seconds())
			if d > 0 {
				params.DurationSeconds = &d
			}
		}

		_, err = uc.exerciseService.Create(userId, params)
		switch {
		case errors.Is(err, dayclosure.ErrDayClosed):
			exOut.Blocked++
		case err != nil:
			return err
		default:
			exOut.Created++
		}
	}

	// 2. Create/update the daily-walk exercise from HC-aggregated daily step totals.
	err := uc.importDailyWalk(userId, stepsDaily, exOut)
	if err != nil {
		return err
	}
	return nil
}

// importDailyWalk creates or updates one daily-walk exercise per day using the
// pre-aggregated daily step total sent by the mobile app. On re-sync the step
// count is overwritten so it always reflects the latest HC aggregate.
func (uc *healthConnectImportUseCase) importDailyWalk(userId uuid.UUID, stepsDaily []ports.DailyStepsRecord, out *ports.TypeResult) error {
	for _, sd := range stepsDaily {
		if sd.Count < dailyWalkMinSteps {
			continue
		}

		// sd.Date is already the correct local calendar date — parse it as UTC
		// midnight so it stores literally. dateOf() must NOT be applied here:
		// it would shift midnight UTC to the previous day in AR time.
		date, err := time.Parse("2006-01-02", sd.Date)
		if err != nil {
			out.Skipped++
			continue
		}

		// If a daily-walk exercise already exists for this day, update its step
		// count. Re-syncs overwrite rather than accumulate.
		existing, err := uc.exerciseRepository.FindByDateAndName(userId, date, dailyWalkName)
		if err != nil && !errors.Is(err, exerciseDomain.ErrExerciseNotFound) {
			return err
		}
		if existing != nil {
			steps := sd.Count
			updateParams := exercisePorts.UpdateParams{
				Steps: &steps,
			}
			_, err = uc.exerciseService.Update(existing.Id, userId, updateParams)
			switch {
			case errors.Is(err, dayclosure.ErrDayClosed):
				out.Blocked++
			case err != nil:
				return err
			default:
				out.Skipped++ // updated in place, not a new creation
			}
			continue
		}

		externalId := externalIdFor("daily_walk", sd.Date)
		importSource := exerciseDomain.ImportSourceHealthConnect
		startedAt := date.Add(12 * time.Hour) // noon UTC = 09:00 ARG, unambiguously the same local date
		params := exercisePorts.CreateParams{
			Date:         date,
			StartedAt:    &startedAt,
			Type:         exerciseDomain.ExerciseTypeWalking,
			Name:         dailyWalkName,
			Steps:        &sd.Count,
			ExternalId:   externalId,
			ImportSource: &importSource,
		}

		_, err = uc.exerciseService.Create(userId, params)
		switch {
		case errors.Is(err, dayclosure.ErrDayClosed):
			out.Blocked++
		case err != nil:
			return err
		default:
			out.Created++
		}
	}
	return nil
}

// ─── Sleep & Heart rate (raw storage) ────────────────────────────────────────

func (uc *healthConnectImportUseCase) importSleep(userId uuid.UUID, records []ports.SleepRecord, out *ports.TypeResult) error {
	for _, rec := range records {
		t, err := parseInstant(rec.StartTime)
		if err != nil {
			out.Skipped++
			continue
		}
		err = uc.storeRaw(userId, "sleep", externalIdFor("sleep", rec.Id), t, rec, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *healthConnectImportUseCase) importHeartRate(userId uuid.UUID, records []ports.HeartRateRecord, out *ports.TypeResult) error {
	for _, rec := range records {
		if len(rec.Samples) == 0 {
			out.Skipped++
			continue
		}
		t, err := parseInstant(rec.Samples[0].Time)
		if err != nil {
			out.Skipped++
			continue
		}
		err = uc.storeRaw(userId, "heart_rate", externalIdFor("heart_rate", rec.Id), t, rec, out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *healthConnectImportUseCase) storeRaw(userId uuid.UUID, recordType string, externalId *string, recordedAt time.Time, raw any, out *ports.TypeResult) error {
	if externalId == nil {
		out.Skipped++
		return nil
	}
	exists, err := uc.rawStore.ExistsByExternalId(userId, *externalId)
	if err != nil {
		return err
	}
	if exists {
		out.Skipped++
		return nil
	}
	payloadBytes, err := json.Marshal(raw)
	if err != nil {
		return cerr.NewInternalError("marshaling raw health record", err)
	}
	rawRecord := &ports.RawRecord{
		UserId:     userId,
		Source:     source,
		Type:       recordType,
		ExternalId: *externalId,
		RecordedAt: recordedAt,
		Payload:    payloadBytes,
	}
	err = uc.rawStore.Create(rawRecord)
	if err != nil {
		return err
	}
	out.Created++
	return nil
}

// ─── Heart rate helpers ───────────────────────────────────────────────────────

func computeHR(hrRecords []ports.HeartRateRecord, start, end time.Time) (avg *int, max *int) {
	var total, count, maxBpm int
	for _, hr := range hrRecords {
		for _, sample := range hr.Samples {
			t, err := parseInstant(sample.Time)
			if err != nil {
				continue
			}
			if t.Before(start) || t.After(end) {
				continue
			}
			total += sample.Bpm
			count++
			if sample.Bpm > maxBpm {
				maxBpm = sample.Bpm
			}
		}
	}
	if count == 0 {
		return nil, nil
	}
	avgBpm := total / count
	return &avgBpm, &maxBpm
}

// ─── Exercise type mapping ────────────────────────────────────────────────────

func mapExerciseType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "running", "jogging", "running_treadmill":
		return exerciseDomain.ExerciseTypeRunning
	case "walking", "hiking":
		return exerciseDomain.ExerciseTypeWalking
	case "biking", "cycling", "biking_stationary":
		return exerciseDomain.ExerciseTypeCycling
	case "strength_training", "weightlifting":
		return exerciseDomain.ExerciseTypeWeightlifting
	default:
		return exerciseDomain.ExerciseTypeOther
	}
}

// ─── ID helpers ───────────────────────────────────────────────────────────────

func externalIdFor(recordType, id string) *string {
	if id == "" {
		return nil
	}
	value := source + ":" + recordType + ":" + id
	return &value
}

// ─── Time helpers ─────────────────────────────────────────────────────────────

// Argentina has no DST — UTC-3 year-round, no OS timezone data required.
var arTz = time.FixedZone("America/Argentina/Buenos_Aires", -3*60*60)

func parseInstant(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("empty time")
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return t.UTC(), nil
	}
	t, err = time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// dateOf returns the calendar date of t in Argentina local time, stored as a
// UTC midnight value. Using UTC date would give the wrong day for sessions that
// happen after 21:00 ARG (= 00:00+ UTC next day).
func dateOf(t time.Time) time.Time {
	local := t.In(arTz)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.UTC)
}
