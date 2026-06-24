package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	cerr "github.com/ivan-ca97/life/pkg/custom_error"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
)

type externalHealthRecordRepository struct {
	db *gorm.DB
}

var _ ports.RawRecordStore = (*externalHealthRecordRepository)(nil)

func NewExternalHealthRecordRepository(db *gorm.DB) *externalHealthRecordRepository {
	return &externalHealthRecordRepository{
		db: db,
	}
}

func (r *externalHealthRecordRepository) ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error) {
	var count int64
	err := r.db.
		Model(&externalHealthRecord{}).
		Where("user_id = ? AND external_id = ?", userId, externalId).
		Count(&count).
		Error
	if err != nil {
		return false, cerr.NewInternalError("checking external health record existence", err)
	}
	return count > 0, nil
}

func (r *externalHealthRecordRepository) Create(rec *ports.RawRecord) error {
	model := &externalHealthRecord{
		Id:         uuid.New(),
		UserId:     rec.UserId,
		Source:     rec.Source,
		Type:       rec.Type,
		ExternalId: rec.ExternalId,
		RecordedAt: rec.RecordedAt,
		Payload:    rec.Payload,
	}
	err := r.db.Create(model).Error
	if err != nil {
		return cerr.NewInternalError("inserting external health record", err)
	}
	return nil
}
