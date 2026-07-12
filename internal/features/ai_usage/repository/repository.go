package repository

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"
	"github.com/ivan-ca97/life/pkg/types"

	"github.com/ivan-ca97/life/internal/features/ai_usage/domain"
	"github.com/ivan-ca97/life/internal/features/ai_usage/ports"
)

type repository struct {
	db *gorm.DB
}

var _ ports.Repository = (*repository)(nil)

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

// FindPrice returns the price effective for the provider/model at `at`, or nil
// when none is on record.
func (r *repository) FindPrice(provider, model string, at time.Time) (*domain.ModelPrice, error) {
	var price aiModelPrice
	err := r.db.
		Where("provider = ? AND model = ? AND effective_from <= ?", provider, model, at).
		Order("effective_from DESC").
		First(&price).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("finding model price", err)
	}
	return price.toDomain(), nil
}

func (r *repository) ListTiers() ([]domain.Tier, error) {
	var models []aiTier
	err := r.db.Order("name ASC").Find(&models).Error
	if err != nil {
		return nil, cerr.NewInternalError("listing ai tiers", err)
	}
	tiers := make([]domain.Tier, len(models))
	for i, m := range models {
		tiers[i] = m.toDomain()
	}
	return tiers, nil
}

func (r *repository) CreateTier(tier *domain.Tier) error {
	var count int64
	err := r.db.Model(&aiTier{}).Where("name = ?", tier.Name).Count(&count).Error
	if err != nil {
		return cerr.NewInternalError("checking ai tier name", err)
	}
	if count > 0 {
		return domain.ErrTierNameTaken
	}
	model := &aiTier{
		Id:              tier.Id,
		Name:            tier.Name,
		MonthlyLimitUsd: tier.MonthlyLimitUsd,
		IsDefault:       tier.IsDefault,
		Enabled:         tier.Enabled,
	}
	err = r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("creating ai tier", err)
	}
	*tier = model.toDomain()
	return nil
}

func (r *repository) UpdateTier(id uuid.UUID, params ports.UpdateTierParams) (*domain.Tier, error) {
	updates := map[string]any{}
	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.MonthlyLimitUsd != nil {
		updates["monthly_limit_usd"] = *params.MonthlyLimitUsd // may be nil -> NULL (unlimited)
	}
	if params.Enabled != nil {
		updates["enabled"] = *params.Enabled
	}

	if len(updates) > 0 {
		result := r.db.Model(&aiTier{}).Where("id = ?", id).Updates(updates)
		if result.Error != nil {
			return nil, cerr.NewInternalError("updating ai tier", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil, domain.ErrTierNotFound
		}
	}

	var model aiTier
	err := r.db.Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTierNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading ai tier", err)
	}
	tier := model.toDomain()
	return &tier, nil
}

func (r *repository) GetDefaultTier() (*domain.Tier, error) {
	var model aiTier
	err := r.db.Where("is_default = ?", true).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTierNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading default ai tier", err)
	}
	tier := model.toDomain()
	return &tier, nil
}

func (r *repository) GetAllocation(userId uuid.UUID) (*domain.Allocation, error) {
	var userTier aiUserTier
	err := r.db.Where("user_id = ?", userId).First(&userTier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No explicit assignment: fall back to the default tier.
		tier, err := r.GetDefaultTier()
		if err != nil {
			return nil, err
		}
		result := &domain.Allocation{
			Tier: *tier,
		}
		return result, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading user ai tier", err)
	}

	var tierModel aiTier
	err = r.db.Where("id = ?", userTier.TierId).First(&tierModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTierNotFound
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading ai tier", err)
	}
	result := &domain.Allocation{
		Tier:         tierModel.toDomain(),
		SelfLimitUsd: userTier.SelfLimitUsd,
	}
	return result, nil
}

func (r *repository) AssignTier(userId, tierId uuid.UUID) error {
	var tier aiTier
	err := r.db.Where("id = ?", tierId).First(&tier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrTierNotFound
	}
	if err != nil {
		return cerr.NewInternalError("loading ai tier", err)
	}
	if !tier.Enabled {
		return domain.ErrTierDisabled
	}
	model := &aiUserTier{
		UserId: userId,
		TierId: tierId,
	}
	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tier_id", "updated_at"}),
	}
	err = r.db.
		Clauses(onConflict).
		Create(model).Error
	if err != nil {
		return cerr.NewInternalError("assigning ai tier", err)
	}
	return nil
}

func (r *repository) DeleteTier(id uuid.UUID) error {
	var tier aiTier
	err := r.db.Where("id = ?", id).First(&tier).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrTierNotFound
	}
	if err != nil {
		return cerr.NewInternalError("loading ai tier", err)
	}
	if tier.IsDefault {
		return cerr.NewConflictError("cannot delete the default tier")
	}

	var assigned int64
	err = r.db.Model(&aiUserTier{}).Where("tier_id = ?", id).Count(&assigned).Error
	if err != nil {
		return cerr.NewInternalError("checking ai tier usage", err)
	}
	if assigned > 0 {
		return domain.ErrTierInUse
	}

	err = r.db.Where("id = ?", id).Delete(&aiTier{}).Error
	if err != nil {
		return cerr.NewInternalError("deleting ai tier", err)
	}
	return nil
}

func (r *repository) SetSelfLimit(userId uuid.UUID, selfLimitUsd *float64) error {
	// Upsert: the row may not exist yet for a user on the default tier. We need a
	// tier_id to satisfy NOT NULL, so default-tier users get the default assigned.
	var existing aiUserTier
	err := r.db.Where("user_id = ?", userId).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tier, err := r.GetDefaultTier()
		if err != nil {
			return err
		}
		row := &aiUserTier{
			UserId:       userId,
			TierId:       tier.Id,
			SelfLimitUsd: selfLimitUsd,
		}
		err = r.db.Create(row).Error
		if err != nil {
			return cerr.NewInternalError("setting ai self limit", err)
		}
		return nil
	}
	if err != nil {
		return cerr.NewInternalError("loading user ai tier", err)
	}
	updates := map[string]any{
		"self_limit_usd": selfLimitUsd,
	}
	result := r.db.Model(&aiUserTier{}).
		Where("user_id = ?", userId).
		Updates(updates)
	if result.Error != nil {
		return cerr.NewInternalError("setting ai self limit", result.Error)
	}
	return nil
}

func (r *repository) GetUsage(userId uuid.UUID, periodStart time.Time) (*domain.Usage, error) {
	var model aiUsage
	err := r.db.Where("user_id = ? AND period_start = ?", userId, periodStart).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No usage yet this period is a zero, not an error.
		result := &domain.Usage{
			UserId:      userId,
			PeriodStart: periodStart,
		}
		return result, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading ai usage", err)
	}
	return model.toDomain(), nil
}

func (r *repository) AddUsage(userId uuid.UUID, periodStart time.Time, delta ports.UsageDelta) error {
	micros := int64(math.Round(delta.CostUsd * 1_000_000))
	model := &aiUsage{
		UserId:        userId,
		PeriodStart:   periodStart,
		Requests:      delta.Requests,
		InputTokens:   delta.InputTokens,
		OutputTokens:  delta.OutputTokens,
		CostUsdMicros: micros,
	}
	assignments := map[string]any{
		"requests":        gorm.Expr("ai_usage.requests + ?", delta.Requests),
		"input_tokens":    gorm.Expr("ai_usage.input_tokens + ?", delta.InputTokens),
		"output_tokens":   gorm.Expr("ai_usage.output_tokens + ?", delta.OutputTokens),
		"cost_usd_micros": gorm.Expr("ai_usage.cost_usd_micros + ?", micros),
		"updated_at":      gorm.Expr("NOW()"),
	}
	onConflict := clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "period_start"}},
		DoUpdates: clause.Assignments(assignments),
	}
	err := r.db.
		Clauses(onConflict).
		Create(model).Error
	if err != nil {
		return cerr.NewInternalError("recording ai usage", err)
	}
	return nil
}

func (r *repository) InsertInteraction(entry ports.InteractionEntry) error {
	metadata := []byte("{}")
	if entry.Metadata != nil {
		encoded, err := json.Marshal(entry.Metadata)
		if err == nil {
			metadata = encoded
		}
	}
	model := &aiInteraction{
		Id:            uuid.New(),
		UserId:        entry.UserId,
		Operation:     entry.Operation,
		Provider:      entry.Provider,
		Model:         entry.Model,
		Status:        entry.Status,
		ErrorType:     entry.ErrorType,
		InputTokens:   entry.InputTokens,
		OutputTokens:  entry.OutputTokens,
		CostUsdMicros: int64(math.Round(entry.CostUsd * 1_000_000)),
		LatencyMs:     entry.LatencyMs,
		ProviderCalls: entry.ProviderCalls,
		CorrelationId: entry.CorrelationId,
		InputSummary:  entry.InputSummary,
		Metadata:      metadata,
	}
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("recording ai interaction", err)
	}
	return nil
}

func (r *repository) ListInteractions(filter ports.InteractionFilter) (types.Page[domain.Interaction], error) {
	q := r.db.Model(&aiInteraction{})
	if filter.UserId != nil {
		q = q.Where("user_id = ?", *filter.UserId)
	}

	var total int64
	err := q.Count(&total).Error
	if err != nil {
		return types.Page[domain.Interaction]{}, cerr.NewInternalError("counting ai interactions", err)
	}

	var models []aiInteraction
	err = q.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&models).Error
	if err != nil {
		return types.Page[domain.Interaction]{}, cerr.NewInternalError("listing ai interactions", err)
	}

	items := make([]domain.Interaction, len(models))
	for i, m := range models {
		items[i] = m.toDomain()
	}
	page := types.Page[domain.Interaction]{
		Items:  items,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}
	return page, nil
}
