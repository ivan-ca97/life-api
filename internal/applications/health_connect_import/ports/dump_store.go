package ports

import "github.com/google/uuid"

type DumpStore interface {
	Save(userId uuid.UUID, appVersion string, payload []byte) error
}
