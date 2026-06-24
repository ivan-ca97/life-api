package ports

import (
	"context"

	"github.com/google/uuid"
)

// Payload is the JSON contract sent by the Android companion app. Each data
// type is a separate array, omitted when empty. Times are ISO-8601 UTC strings.
type Payload struct {
	SyncedAt         string             `json:"synced_at"`
	AppVersion       string             `json:"app_version"`
	Weight           []WeightRecord     `json:"weight"`
	ExerciseSessions []ExerciseRecord   `json:"exercise_sessions"`
	Steps            []StepsRecord      `json:"steps"`       // deprecated: ignored, kept for backwards compat
	StepsDaily       []DailyStepsRecord `json:"steps_daily"` // aggregated daily totals (HC dedup)
	Sleep            []SleepRecord      `json:"sleep"`
	HeartRate        []HeartRateRecord  `json:"heart_rate"`
}

// Id on every record is the Health Connect metadata id (a stable UUID), used to
// deduplicate across re-syncs.

type WeightRecord struct {
	Id        string  `json:"id"`
	Kilograms float64 `json:"kilograms"`
	Time      string  `json:"time"`
}

type ExerciseRecord struct {
	Id              string   `json:"id"`
	Type            string   `json:"type"`
	StartTime       string   `json:"start_time"`
	EndTime         string   `json:"end_time"`
	DurationSeconds *int     `json:"duration_seconds"`
	DistanceMeters  *float64 `json:"distance_meters"`
	Steps           *int     `json:"steps"`        // HC-aggregated steps for the session window (null if 0)
	Title           string   `json:"title"`
	DataOrigin      string   `json:"data_origin"`  // package name of the originating app (e.g. "hevy")
}

type StepsRecord struct {
	Id        string `json:"id"`
	Count     int    `json:"count"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type DailyStepsRecord struct {
	Date  string `json:"date"`  // YYYY-MM-DD (local date)
	Count int    `json:"count"` // HC-deduplicated total for the day
}

type SleepStage struct {
	Stage     string `json:"stage"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type SleepRecord struct {
	Id        string       `json:"id"`
	StartTime string       `json:"start_time"`
	EndTime   string       `json:"end_time"`
	Stages    []SleepStage `json:"stages"`
}

type HeartRateSample struct {
	Time string `json:"time"`
	Bpm  int    `json:"bpm"`
}

type HeartRateRecord struct {
	Id      string            `json:"id"`
	Samples []HeartRateSample `json:"samples"`
}

// TypeResult counts the outcome of importing one data type.
type TypeResult struct {
	Created int
	Skipped int
	Blocked int
}

type ImportResult struct {
	Weight    TypeResult
	Exercise  TypeResult
	Steps     TypeResult
	Sleep     TypeResult
	HeartRate TypeResult
}

type HealthConnectImportUseCase interface {
	Import(ctx context.Context, userId uuid.UUID, payload *Payload) (*ImportResult, error)
}
