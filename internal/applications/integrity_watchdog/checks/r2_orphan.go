package checks

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type R2OrphanResult struct {
	// Keys in R2 not referenced by any DB record.
	OrphanKeys []string
	// Keys that were found as orphans and successfully deleted from R2.
	DeletedKeys []string
	// URLs stored in DB that no longer exist in R2.
	BrokenRefs []string
}

func (r *R2OrphanResult) IsClean() bool {
	return len(r.OrphanKeys) == 0 && len(r.DeletedKeys) == 0 && len(r.BrokenRefs) == 0
}

type R2OrphanCheck struct {
	lister     ports.ObjectLister
	deleter    ports.ObjectDeleter // nil = report only, non-nil = auto-delete
	repository ports.WatchdogRepository
	publicUrl  string
}

func NewR2OrphanCheck(lister ports.ObjectLister, deleter ports.ObjectDeleter, repository ports.WatchdogRepository, publicUrl string) *R2OrphanCheck {
	return &R2OrphanCheck{
		lister:     lister,
		deleter:    deleter,
		repository: repository,
		publicUrl:  publicUrl,
	}
}

func (c *R2OrphanCheck) Run() (*R2OrphanResult, error) {
	keys, err := c.lister.ListAllKeys("")
	if err != nil {
		return nil, fmt.Errorf("listing R2: %w", err)
	}

	dbUrls, err := c.repository.AllPhotoUrls()
	if err != nil {
		return nil, fmt.Errorf("querying DB URLs: %w", err)
	}

	base := strings.TrimRight(c.publicUrl, "/") + "/"

	dbUrlSet := make(map[string]bool, len(dbUrls))
	for _, u := range dbUrls {
		dbUrlSet[u] = true
	}

	r2KeySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		r2KeySet[k] = true
	}

	var orphans []string
	for _, k := range keys {
		if !dbUrlSet[base+k] {
			orphans = append(orphans, k)
		}
	}

	var broken []string
	for _, u := range dbUrls {
		key := strings.TrimPrefix(u, base)
		if !r2KeySet[key] {
			broken = append(broken, u)
		}
	}

	result := &R2OrphanResult{
		BrokenRefs: broken,
	}

	if c.deleter != nil && len(orphans) > 0 {
		err := c.deleter.DeleteKeys(orphans)
		if err != nil {
			slog.Error("r2 orphan cleanup: delete failed", "error", err, "count", len(orphans))
			result.OrphanKeys = orphans
		} else {
			slog.Info("r2 orphan cleanup: deleted orphans", "count", len(orphans))
			result.DeletedKeys = orphans
		}
	} else {
		result.OrphanKeys = orphans
	}

	return result, nil
}
