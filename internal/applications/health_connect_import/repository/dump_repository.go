package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ivan-ca97/life/internal/applications/health_connect_import/ports"
)

type hcSyncDump struct {
	Id         int64     `gorm:"primaryKey;autoIncrement"`
	UserId     uuid.UUID `gorm:"not null"`
	AppVersion string    `gorm:"not null;default:''"`
	ReceivedAt time.Time `gorm:"not null;autoCreateTime"`
	Payload    []byte    `gorm:"type:jsonb;not null"`
}

func (hcSyncDump) TableName() string { return "hc_sync_dumps" }

type dumpRepository struct {
	db *gorm.DB
}

func NewDumpRepository(db *gorm.DB) ports.DumpStore {
	return &dumpRepository{db: db}
}

func (r *dumpRepository) Save(userId uuid.UUID, appVersion string, payload []byte) error {
	dump := &hcSyncDump{
		UserId:     userId,
		AppVersion: appVersion,
		Payload:    payload,
	}
	return r.db.Create(dump).Error
}
