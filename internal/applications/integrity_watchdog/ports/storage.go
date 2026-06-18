package ports

// ObjectLister lists object keys from object storage.
type ObjectLister interface {
	ListAllKeys(prefix string) ([]string, error)
}

// ObjectDeleter deletes objects from object storage by key.
type ObjectDeleter interface {
	DeleteKeys(keys []string) error
}
