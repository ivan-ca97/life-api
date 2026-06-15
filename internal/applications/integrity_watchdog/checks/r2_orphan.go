package checks

import (
	"fmt"
	"strings"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/ports"
)

type R2OrphanResult struct {
	// Keys in R2 not referenced by any DB record.
	OrphanKeys []string
	// URLs stored in DB that no longer exist in R2.
	BrokenRefs []string
}

func (r *R2OrphanResult) IsClean() bool {
	return len(r.OrphanKeys) == 0 && len(r.BrokenRefs) == 0
}

type R2OrphanCheck struct {
	lister     ports.ObjectLister
	repository ports.WatchdogRepository
	publicURL  string
}

func NewR2OrphanCheck(lister ports.ObjectLister, repository ports.WatchdogRepository, publicURL string) *R2OrphanCheck {
	return &R2OrphanCheck{lister: lister, repository: repository, publicURL: publicURL}
}

func (c *R2OrphanCheck) Run() (*R2OrphanResult, error) {
	keys, err := c.lister.ListAllKeys("")
	if err != nil {
		return nil, fmt.Errorf("listing R2: %w", err)
	}

	dbURLs, err := c.repository.AllPhotoURLs()
	if err != nil {
		return nil, fmt.Errorf("querying DB URLs: %w", err)
	}

	base := strings.TrimRight(c.publicURL, "/") + "/"

	dbURLSet := make(map[string]bool, len(dbURLs))
	for _, u := range dbURLs {
		dbURLSet[u] = true
	}

	r2KeySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		r2KeySet[k] = true
	}

	var orphans []string
	for _, k := range keys {
		if !dbURLSet[base+k] {
			orphans = append(orphans, k)
		}
	}

	var broken []string
	for _, u := range dbURLs {
		key := strings.TrimPrefix(u, base)
		if !r2KeySet[key] {
			broken = append(broken, u)
		}
	}

	return &R2OrphanResult{OrphanKeys: orphans, BrokenRefs: broken}, nil
}
