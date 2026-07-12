package repository

import (
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
)

type syncLogRepository struct {
	db *gorm.DB
}

var _ ports.SyncLogStore = (*syncLogRepository)(nil)

func NewSyncLogRepository(db *gorm.DB) *syncLogRepository {
	return &syncLogRepository{
		db: db,
	}
}

func (r *syncLogRepository) Create(log *ports.SyncLog) error {
	model := &syncLog{
		UserId:     log.UserId,
		AppVersion: log.AppVersion,
		SyncedAt:   log.SyncedAt,
		Result:     log.Result,
	}
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting sync log", err)
	}
	return nil
}
