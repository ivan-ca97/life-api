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
	return &repository{db: db}
}

func (r *repository) ListTiers() ([]domain.Tier, error) {
	var models []aiTier
	if err := r.db.Order("name ASC").Find(&models).Error; err != nil {
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
	if err := r.db.Model(&aiTier{}).Where("name = ?", tier.Name).Count(&count).Error; err != nil {
		return cerr.NewInternalError("checking ai tier name", err)
	}
	if count > 0 {
		return domain.ErrTierNameTaken
	}
	model := &aiTier{
		Id:              tier.Id,
		Name:            tier.Name,
		MonthlyLimitUSD: tier.MonthlyLimitUSD,
		IsDefault:       tier.IsDefault,
		Enabled:         tier.Enabled,
	}
	if err := r.db.Create(model).Error; err != nil {
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
	if params.MonthlyLimitUSD != nil {
		updates["monthly_limit_usd"] = *params.MonthlyLimitUSD // may be nil -> NULL (unlimited)
	}
	if params.Enabled != nil {
		updates["enabled"] = *params.Enabled
	}

	if len(updates) > 0 {
		res := r.db.Model(&aiTier{}).Where("id = ?", id).Updates(updates)
		if res.Error != nil {
			return nil, cerr.NewInternalError("updating ai tier", res.Error)
		}
		if res.RowsAffected == 0 {
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
		return &domain.Allocation{Tier: *tier}, nil
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
	return &domain.Allocation{
		Tier:         tierModel.toDomain(),
		SelfLimitUSD: userTier.SelfLimitUSD,
	}, nil
}

func (r *repository) AssignTier(userId, tierId uuid.UUID) error {
	var count int64
	if err := r.db.Model(&aiTier{}).Where("id = ?", tierId).Count(&count).Error; err != nil {
		return cerr.NewInternalError("checking ai tier", err)
	}
	if count == 0 {
		return domain.ErrTierNotFound
	}
	model := &aiUserTier{UserId: userId, TierId: tierId}
	err := r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"tier_id", "updated_at"}),
		}).
		Create(model).Error
	if err != nil {
		return cerr.NewInternalError("assigning ai tier", err)
	}
	return nil
}

func (r *repository) SetSelfLimit(userId uuid.UUID, selfLimitUSD *float64) error {
	// Upsert: the row may not exist yet for a user on the default tier. We need a
	// tier_id to satisfy NOT NULL, so default-tier users get the default assigned.
	var existing aiUserTier
	err := r.db.Where("user_id = ?", userId).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		tier, err := r.GetDefaultTier()
		if err != nil {
			return err
		}
		row := &aiUserTier{UserId: userId, TierId: tier.Id, SelfLimitUSD: selfLimitUSD}
		if err := r.db.Create(row).Error; err != nil {
			return cerr.NewInternalError("setting ai self limit", err)
		}
		return nil
	}
	if err != nil {
		return cerr.NewInternalError("loading user ai tier", err)
	}
	res := r.db.Model(&aiUserTier{}).
		Where("user_id = ?", userId).
		Updates(map[string]any{"self_limit_usd": selfLimitUSD})
	if res.Error != nil {
		return cerr.NewInternalError("setting ai self limit", res.Error)
	}
	return nil
}

func (r *repository) GetUsage(userId uuid.UUID, periodStart time.Time) (*domain.Usage, error) {
	var model aiUsage
	err := r.db.Where("user_id = ? AND period_start = ?", userId, periodStart).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No usage yet this period is a zero, not an error.
		return &domain.Usage{UserId: userId, PeriodStart: periodStart}, nil
	}
	if err != nil {
		return nil, cerr.NewInternalError("loading ai usage", err)
	}
	return model.toDomain(), nil
}

func (r *repository) AddUsage(userId uuid.UUID, periodStart time.Time, delta ports.UsageDelta) error {
	micros := int64(math.Round(delta.CostUSD * 1_000_000))
	model := &aiUsage{
		UserId:        userId,
		PeriodStart:   periodStart,
		Requests:      delta.Requests,
		InputTokens:   delta.InputTokens,
		OutputTokens:  delta.OutputTokens,
		CostUsdMicros: micros,
	}
	err := r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "period_start"}},
			DoUpdates: clause.Assignments(map[string]any{
				"requests":        gorm.Expr("ai_usage.requests + ?", delta.Requests),
				"input_tokens":    gorm.Expr("ai_usage.input_tokens + ?", delta.InputTokens),
				"output_tokens":   gorm.Expr("ai_usage.output_tokens + ?", delta.OutputTokens),
				"cost_usd_micros": gorm.Expr("ai_usage.cost_usd_micros + ?", micros),
				"updated_at":      gorm.Expr("NOW()"),
			}),
		}).
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
		CostUsdMicros: int64(math.Round(entry.CostUSD * 1_000_000)),
		LatencyMs:     entry.LatencyMs,
		ProviderCalls: entry.ProviderCalls,
		CorrelationId: entry.CorrelationId,
		InputSummary:  entry.InputSummary,
		Metadata:      metadata,
	}
	if err := r.db.Create(model).Error; err != nil {
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
	if err := q.Count(&total).Error; err != nil {
		return types.Page[domain.Interaction]{}, cerr.NewInternalError("counting ai interactions", err)
	}

	var models []aiInteraction
	err := q.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&models).Error
	if err != nil {
		return types.Page[domain.Interaction]{}, cerr.NewInternalError("listing ai interactions", err)
	}

	items := make([]domain.Interaction, len(models))
	for i, m := range models {
		items[i] = m.toDomain()
	}
	return types.Page[domain.Interaction]{
		Items:  items,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}
