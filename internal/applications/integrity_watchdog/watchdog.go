package integrity_watchdog

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/ivan-ca97/life/pkg/api/http_errors"
	"github.com/ivan-ca97/life/pkg/auth"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/checks"
	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/handler"
	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/repository"
	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/scheduler"
	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/storage"
)

type WatchdogApplication struct {
	scheduler       *scheduler.Scheduler
	watchdogHandler handler.WatchdogHandler
	errorHandler    http_errors.HttpErrorHandler
}

func NewWatchdogApplication(
	db *gorm.DB,
	period time.Duration,
	r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL string,
	authorizer auth.AuthorizationService,
	errorHandler http_errors.HttpErrorHandler,
) *WatchdogApplication {
	repo := repository.NewWatchdogRepository(db)
	dbCheck := checks.NewDBIntegrityCheck(repo)

	var r2Check *checks.R2OrphanCheck
	if r2AccountId != "" {
		lister := storage.NewR2Lister(r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket)
		r2Check = checks.NewR2OrphanCheck(lister, repo, r2PublicURL)
	}

	sched := scheduler.New(period, dbCheck, r2Check)
	h := handler.NewWatchdogHandler(sched, authorizer)

	return &WatchdogApplication{
		scheduler:       sched,
		watchdogHandler: h,
		errorHandler:    errorHandler,
	}
}

// Start begins the background periodic loop. Blocks until ctx is cancelled — call in a goroutine.
func (a *WatchdogApplication) Start(ctx context.Context) {
	a.scheduler.Start(ctx)
}
