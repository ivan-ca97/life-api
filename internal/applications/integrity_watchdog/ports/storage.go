package ports

// ObjectLister lists object keys from object storage.
type ObjectLister interface {
	ListAllKeys(prefix string) ([]string, error)
}
