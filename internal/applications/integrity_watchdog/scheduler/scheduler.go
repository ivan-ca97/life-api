package scheduler

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/ivan-ca97/life/internal/applications/integrity_watchdog/checks"
)

type RunResult struct {
	StartedAt  time.Time
	FinishedAt time.Time
	DB         *checks.DBIntegrityResult
	R2         *checks.R2OrphanResult
	Err        error
}

type Scheduler struct {
	periodNs   atomic.Int64
	trigger    chan struct{}
	reconfig   chan time.Duration
	running    atomic.Bool
	lastResult atomic.Pointer[RunResult]
	dbCheck    *checks.DBIntegrityCheck
	r2Check    *checks.R2OrphanCheck // nil when R2 is not configured
}

func New(period time.Duration, dbCheck *checks.DBIntegrityCheck, r2Check *checks.R2OrphanCheck) *Scheduler {
	s := &Scheduler{
		trigger:  make(chan struct{}, 1),
		reconfig: make(chan time.Duration, 1),
		dbCheck:  dbCheck,
		r2Check:  r2Check,
	}
	s.periodNs.Store(int64(period))
	return s
}

// Start runs the periodic loop. Blocks until ctx is cancelled — call in a goroutine.
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.Period())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runAsync(ctx)
		case <-s.trigger:
			ticker.Reset(s.Period())
			s.runAsync(ctx)
		case d := <-s.reconfig:
			s.periodNs.Store(int64(d))
			ticker.Reset(d)
		case <-ctx.Done():
			return
		}
	}
}

// Trigger requests an immediate run. Returns false if a trigger is already pending.
func (s *Scheduler) Trigger() bool {
	select {
	case s.trigger <- struct{}{}:
		return true
	default:
		return false
	}
}

// SetPeriod reconfigures the ticker period for the next tick.
func (s *Scheduler) SetPeriod(d time.Duration) {
	s.periodNs.Store(int64(d))
	// Replace any pending reconfig rather than queuing two.
	select {
	case s.reconfig <- d:
	default:
		select {
		case <-s.reconfig:
		default:
		}
		s.reconfig <- d
	}
}

func (s *Scheduler) Period() time.Duration {
	return time.Duration(s.periodNs.Load())
}

func (s *Scheduler) IsRunning() bool {
	return s.running.Load()
}

func (s *Scheduler) LastResult() *RunResult {
	return s.lastResult.Load()
}

func (s *Scheduler) runAsync(ctx context.Context) {
	if !s.running.CompareAndSwap(false, true) {
		return
	}
	go func() {
		defer s.running.Store(false)
		result := s.runAll(ctx)
		s.lastResult.Store(result)
	}()
}

func (s *Scheduler) runAll(ctx context.Context) *RunResult {
	result := &RunResult{
		StartedAt: time.Now(),
	}
	slog.Info("integrity watchdog: run started")

	dbResult, err := s.dbCheck.Run()
	if err != nil {
		slog.Error("integrity watchdog: DB check failed", "error", err)
		result.Err = err
		result.FinishedAt = time.Now()
		return result
	}
	result.DB = dbResult

	if s.r2Check != nil {
		r2Result, err := s.r2Check.Run()
		if err != nil {
			slog.Error("integrity watchdog: R2 check failed", "error", err)
			result.Err = err
			result.FinishedAt = time.Now()
			return result
		}
		result.R2 = r2Result
	}

	result.FinishedAt = time.Now()
	slog.Info("integrity watchdog: run finished",
		"duration_ms", result.FinishedAt.Sub(result.StartedAt).Milliseconds(),
		"db_clean", result.DB.IsClean(),
	)
	return result
}
