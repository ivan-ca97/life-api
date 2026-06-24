package ports

import (
	"time"

	"github.com/google/uuid"
)

// RawRecord is a health record persisted verbatim, for data types that don't
// have a dedicated feature yet (sleep, heart rate). Payload holds the original
// JSON of the record so nothing is lost before proper modeling.
type RawRecord struct {
	UserId     uuid.UUID
	Source     string
	Type       string
	ExternalId string
	RecordedAt time.Time
	Payload    []byte
}

type RawRecordStore interface {
	ExistsByExternalId(userId uuid.UUID, externalId string) (bool, error)
	Create(record *RawRecord) error
}

// SyncLog is a lightweight summary written after each successful HC import.
type SyncLog struct {
	UserId     uuid.UUID
	AppVersion string
	SyncedAt   time.Time
	Result     []byte // JSON-encoded ImportResult counts
}

type SyncLogStore interface {
	Create(log *SyncLog) error
}
