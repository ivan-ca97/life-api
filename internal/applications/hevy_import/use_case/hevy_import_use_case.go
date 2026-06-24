package use_case

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ivan-ca97/life/pkg/auth"
	"github.com/ivan-ca97/life/pkg/dayclosure"

	"github.com/ivan-ca97/life/internal/applications/hevy_import/ports"
	exerciseDomain "github.com/ivan-ca97/life/internal/features/exercise/domain"
	exercisePorts "github.com/ivan-ca97/life/internal/features/exercise/ports"
	"github.com/ivan-ca97/life/internal/permissions"
)

var spanishMonths = map[string]string{
	"ene": "Jan", "feb": "Feb", "mar": "Mar", "abr": "Apr",
	"may": "May", "jun": "Jun", "jul": "Jul", "ago": "Aug",
	"sep": "Sep", "oct": "Oct", "nov": "Nov", "dic": "Dec",
}

type hevyRow struct {
	sessionTitle  string
	startTime     time.Time
	endTime       time.Time
	exerciseTitle string
	setIndex      int
	setType       string
	weightKg      *float64
	reps          *int
}

type sessionExerciseKey struct {
	startTimeUnix int64
	exerciseName  string
}

type aggregatedExercise struct {
	sessionTitle string
	date         time.Time
	startedAt    time.Time
	durationSecs int
	exerciseName string
	sets         []hevyRow
}

type hevyImportUseCase struct {
	exerciseService    exercisePorts.ExerciseService
	exerciseRepository exercisePorts.ExerciseRepository
	authorizer         auth.AuthorizationService
}

var _ ports.HevyImportUseCase = (*hevyImportUseCase)(nil)

func NewHevyImportUseCase(
	exerciseService exercisePorts.ExerciseService,
	exerciseRepository exercisePorts.ExerciseRepository,
	authorizer auth.AuthorizationService,
) *hevyImportUseCase {
	return &hevyImportUseCase{
		exerciseService:    exerciseService,
		exerciseRepository: exerciseRepository,
		authorizer:         authorizer,
	}
}

func (uc *hevyImportUseCase) Import(ctx context.Context, userId uuid.UUID, csvReader io.Reader) (*ports.ImportResult, error) {
	err := uc.authorizer.Authorize(ctx, userId, permissions.ExercisesCreate)
	if err != nil {
		return nil, err
	}

	rows, err := parseCSV(csvReader)
	if err != nil {
		return nil, err
	}

	exercises := aggregate(rows)

	result := &ports.ImportResult{
		Results: make([]ports.ImportResultItem, 0, len(exercises)),
	}

	for _, ex := range exercises {
		dateString := ex.date.Format("2006-01-02")

		exists, err := uc.exerciseRepository.ExistsByDateAndName(userId, ex.date, ex.exerciseName)
		if err != nil {
			return nil, err
		}
		if exists {
			existing, err := uc.exerciseRepository.FindByDateAndName(userId, ex.date, ex.exerciseName)
			if err != nil {
				return nil, err
			}
			if _, err := uc.exerciseRepository.Update(existing.Id, userId, buildEnrichParams(ex)); err != nil {
				return nil, err
			}
			result.Enriched++
			result.Results = append(result.Results, ports.ImportResultItem{
				Date:   dateString,
				Name:   ex.exerciseName,
				Status: "enriched",
			})
			continue
		}

		params := buildCreateParams(ex)
		_, err = uc.exerciseService.Create(userId, params)
		if errors.Is(err, dayclosure.ErrDayClosed) {
			result.Blocked++
			result.Results = append(result.Results, ports.ImportResultItem{
				Date:   dateString,
				Name:   ex.exerciseName,
				Status: "blocked",
				Reason: "day is closed",
			})
			continue
		}
		if err != nil {
			return nil, err
		}

		result.Created++
		result.Results = append(result.Results, ports.ImportResultItem{
			Date:   dateString,
			Name:   ex.exerciseName,
			Status: "created",
		})
	}

	return result, nil
}

func parseCSV(reader io.Reader) ([]hevyRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true

	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("reading CSV headers: %w", err)
	}

	colIndex := map[string]int{}
	for i, h := range headers {
		colIndex[strings.TrimSpace(h)] = i
	}

	requiredColumns := []string{"title", "start_time", "end_time", "exercise_title", "set_index", "set_type"}
	for _, col := range requiredColumns {
		if _, ok := colIndex[col]; !ok {
			return nil, fmt.Errorf("missing required column: %s", col)
		}
	}

	var rows []hevyRow
	lineNumber := 1
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading CSV line %d: %w", lineNumber+1, err)
		}
		lineNumber++

		row, err := parseRow(record, colIndex, lineNumber)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func parseRow(record []string, colIndex map[string]int, line int) (hevyRow, error) {
	getField := func(name string) string {
		idx, ok := colIndex[name]
		if !ok || idx >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[idx])
	}

	startTime, err := parseSpanishDateTime(getField("start_time"))
	if err != nil {
		return hevyRow{}, fmt.Errorf("line %d: %w", line, err)
	}

	endTime, err := parseSpanishDateTime(getField("end_time"))
	if err != nil {
		return hevyRow{}, fmt.Errorf("line %d: %w", line, err)
	}

	setIndex, err := strconv.Atoi(getField("set_index"))
	if err != nil {
		return hevyRow{}, fmt.Errorf("line %d: parsing set_index: %w", line, err)
	}

	var weightKg *float64
	if raw := getField("weight_kg"); raw != "" {
		w, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return hevyRow{}, fmt.Errorf("line %d: parsing weight_kg: %w", line, err)
		}
		weightKg = &w
	}

	var reps *int
	if raw := getField("reps"); raw != "" {
		r, err := strconv.Atoi(raw)
		if err != nil {
			return hevyRow{}, fmt.Errorf("line %d: parsing reps: %w", line, err)
		}
		reps = &r
	}

	row := hevyRow{
		sessionTitle:  getField("title"),
		startTime:     startTime,
		endTime:       endTime,
		exerciseTitle: getField("exercise_title"),
		setIndex:      setIndex,
		setType:       getField("set_type"),
		weightKg:      weightKg,
		reps:          reps,
	}
	return row, nil
}

func parseSpanishDateTime(s string) (time.Time, error) {
	normalized := s
	for spanish, english := range spanishMonths {
		normalized = strings.ReplaceAll(normalized, " "+spanish+" ", " "+english+" ")
	}
	t, err := time.Parse("2 Jan 2006, 15:04", normalized)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing date %q: %w", s, err)
	}
	return t, nil
}

func aggregate(rows []hevyRow) []aggregatedExercise {
	exerciseMap := map[sessionExerciseKey]*aggregatedExercise{}
	var order []sessionExerciseKey

	for _, row := range rows {
		key := sessionExerciseKey{
			startTimeUnix: row.startTime.Unix(),
			exerciseName:  row.exerciseTitle,
		}

		agg, exists := exerciseMap[key]
		if !exists {
			durationSecs := int(row.endTime.Sub(row.startTime).Seconds())
			date := time.Date(
				row.startTime.Year(),
				row.startTime.Month(),
				row.startTime.Day(),
				0, 0, 0, 0,
				time.UTC,
			)
			agg = &aggregatedExercise{
				sessionTitle: row.sessionTitle,
				date:         date,
				startedAt:    row.startTime,
				durationSecs: durationSecs,
				exerciseName: row.exerciseTitle,
			}
			exerciseMap[key] = agg
			order = append(order, key)
		}
		agg.sets = append(agg.sets, row)
	}

	result := make([]aggregatedExercise, len(order))
	for i, key := range order {
		result[i] = *exerciseMap[key]
	}
	return result
}

func buildEnrichParams(ex aggregatedExercise) exercisePorts.UpdateParams {
	totalSets := len(ex.sets)
	notes := buildNotes(ex.sets)
	tags := []string{ex.sessionTitle}
	importSource := exerciseDomain.ImportSourceHealthConnectHevy

	params := exercisePorts.UpdateParams{
		TotalSets:    &totalSets,
		Notes:        &notes,
		Tags:         &tags,
		ImportSource: &importSource,
	}

	var totalVolume float64
	for _, set := range ex.sets {
		if set.weightKg != nil && set.reps != nil {
			totalVolume += *set.weightKg * float64(*set.reps)
			params.TotalVolumeKg = &totalVolume
		}
	}
	return params
}

func buildCreateParams(ex aggregatedExercise) exercisePorts.CreateParams {
	importSource := exerciseDomain.ImportSourceHevy
	startedAt := ex.startedAt
	durationSecs := ex.durationSecs
	totalSets := len(ex.sets)

	var totalVolume float64
	hasVolume := false
	for _, set := range ex.sets {
		if set.weightKg != nil && set.reps != nil {
			totalVolume += *set.weightKg * float64(*set.reps)
			hasVolume = true
		}
	}

	var totalVolumePtr *float64
	if hasVolume {
		totalVolumePtr = &totalVolume
	}

	notes := buildNotes(ex.sets)

	params := exercisePorts.CreateParams{
		Date:            ex.date,
		Type:            exerciseDomain.ExerciseTypeWeightlifting,
		Name:            ex.exerciseName,
		StartedAt:       &startedAt,
		DurationSeconds: &durationSecs,
		TotalVolumeKg:   totalVolumePtr,
		TotalSets:       &totalSets,
		Tags:            []string{ex.sessionTitle},
		Notes:           notes,
		ImportSource:    &importSource,
	}
	return params
}

func buildNotes(sets []hevyRow) string {
	var builder strings.Builder
	for i, set := range sets {
		if i > 0 {
			builder.WriteString("\n")
		}

		setNum := set.setIndex + 1
		if set.setType == "normal" || set.setType == "" {
			builder.WriteString(fmt.Sprintf("Set %d: ", setNum))
		} else {
			builder.WriteString(fmt.Sprintf("Set %d (%s): ", setNum, set.setType))
		}

		if set.weightKg != nil {
			builder.WriteString(fmt.Sprintf("%.1f kg", *set.weightKg))
		} else {
			builder.WriteString("BW")
		}

		if set.reps != nil {
			builder.WriteString(fmt.Sprintf(" × %d", *set.reps))
		}
	}
	return builder.String()
}
