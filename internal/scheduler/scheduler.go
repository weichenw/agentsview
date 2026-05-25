package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Scheduler manages cron-driven job execution.
type Scheduler struct {
	mu       sync.Mutex
	store    *Store
	runner   *Runner
	cron     *cron.Cron
	entries  map[string]cron.EntryID // job ID -> cron entry ID
	stopCh   chan struct{}
	stopOnce sync.Once
}

// New creates a Scheduler with the given timezone. Call Start() to begin execution.
func New(store *Store, runner *Runner, loc *time.Location) *Scheduler {
	if loc == nil {
		loc = time.Local
	}
	return &Scheduler{
		store:   store,
		runner:  runner,
		cron:    cron.New(cron.WithLocation(loc)),
		entries: make(map[string]cron.EntryID),
		stopCh:  make(chan struct{}),
	}
}

// Start loads all enabled jobs and registers them with the cron runner.
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, job := range s.store.List() {
		if !job.Enabled {
			continue
		}
		s.addJobLocked(&job)
	}

	s.cron.Start()
	log.Printf("scheduler: started with timezone %s, %d active job(s)", s.cron.Location().String(), len(s.entries))

	go s.heartbeatLoop()
}

// Stop gracefully shuts down the cron runner, waiting for any
// running jobs to finish.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.cron.Stop() // returns context that's done when all jobs finish
	<-ctx.Done()
	s.stopOnce.Do(func() { close(s.stopCh) })
	log.Printf("scheduler: stopped")
}

// Reload removes the cron entry for a job and re-adds it if enabled.
// Called by the HTTP handlers after create/update/enable/disable/delete.
func (s *Scheduler) Reload(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing entry if any.
	if entryID, ok := s.entries[id]; ok {
		s.cron.Remove(entryID)
		delete(s.entries, id)
	}

	// Re-add if the job exists and is enabled.
	job := s.store.Get(id)
	if job == nil || !job.Enabled {
		return
	}

	s.addJobLocked(job)
}

// SetTimezone loads the named timezone and restarts the cron runner with it.
// Re-registers all enabled jobs after the restart.
func (s *Scheduler) SetTimezone(tz string) error {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone %q: %w", tz, err)
	}
	s.RestartWithLocation(loc)
	return nil
}

// RestartWithLocation stops the cron runner, recreates it with the
// given location, and re-registers all enabled jobs.
func (s *Scheduler) RestartWithLocation(loc *time.Location) {
	if loc == nil {
		loc = time.Local
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.cron.Stop()
	<-ctx.Done()

	s.cron = cron.New(cron.WithLocation(loc))
	s.entries = make(map[string]cron.EntryID)

	for _, job := range s.store.List() {
		if !job.Enabled {
			continue
		}
		s.addJobLocked(&job)
	}

	s.cron.Start()
	log.Printf("scheduler: restarted with timezone %s, %d active job(s)", loc.String(), len(s.entries))
}

// heartbeatLoop ticks every 60 seconds. If it detects a gap > 3 minutes
// between ticks, the system was likely asleep — check for missed jobs.
func (s *Scheduler) heartbeatLoop() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	lastCheck := time.Now()
	for {
		select {
		case <-s.stopCh:
			return
		case now := <-ticker.C:
			gap := now.Sub(lastCheck)
			lastCheck = now // update BEFORE any work

			if gap < 3*time.Minute {
				continue
			}

			log.Printf("scheduler: detected system sleep/wake, gap was %v, checking for missed jobs", gap)
			s.catchUpMissedRuns(now.Add(-gap), now)
		}
	}
}

// catchUpMissedRuns finds the latest missed firing for each enabled job
// within [sleepStart, now] and executes it if it hasn't already run.
// Strategy: fire only the LATEST missed slot (not every missed slot).
func (s *Scheduler) catchUpMissedRuns(sleepStart, now time.Time) {
	// Read the cron timezone under lock to avoid racing with RestartWithLocation.
	s.mu.Lock()
	loc := s.cron.Location()
	s.mu.Unlock()

	// Create a parser for standard 5-field cron expressions.
	// The timezone is carried by the scheduler's cron instance, so
	// times are compared/located via s.cron.Location().
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	// Pad slightly before sleepStart so we don't miss a firing that
	// occurred right at the boundary.
	padStart := sleepStart.Add(-1 * time.Minute)

	for _, job := range s.store.List() {
		if !job.Enabled {
			continue
		}

		schedule, err := parser.Parse(job.Cron)
		if err != nil {
			log.Printf("scheduler: heartbeat: invalid cron for job %q: %v", job.ID, err)
			continue
		}

		// Find the LATEST firing time within the gap window.
		var latestMissed *time.Time
		t := padStart.In(loc)
		until := now.In(loc)

		for {
			next := schedule.Next(t)
			if next.IsZero() || next.After(until) {
				break
			}
			latestMissed = &next
			t = next
		}

		if latestMissed == nil {
			// Nothing missed for this job.
			continue
		}

		// Avoid double-firing: check if a run already exists at or after
		// the missed scheduled time.
		alreadyFired, err := s.wasAlreadyFired(job.ID, *latestMissed)
		if err != nil {
			log.Printf("scheduler: heartbeat: check runs for job %q: %v", job.ID, err)
			continue
		}
		if alreadyFired {
			log.Printf("scheduler: heartbeat: job %q already fired for %v, skipping", job.ID, latestMissed.Format(time.RFC3339))
			continue
		}

		// Re-fetch the job under lock to make sure it's still enabled.
		s.mu.Lock()
		curJob := s.store.Get(job.ID)
		if curJob == nil || !curJob.Enabled {
			s.mu.Unlock()
			continue
		}
		jobCopy := *curJob
		s.mu.Unlock()

		log.Printf("scheduler: catch-up firing job %q (%s) missed at %v", jobCopy.ID, jobCopy.Name, latestMissed.Format(time.RFC3339))
		s.runner.Run(&jobCopy, newRunID())
	}
}

// wasAlreadyFired returns true if the given job has a run whose
// started_at is on or after scheduleTime.
func (s *Scheduler) wasAlreadyFired(jobID string, scheduleTime time.Time) (bool, error) {
	recentRuns, err := s.store.ListRuns(jobID, 1)
	if err != nil {
		if err.Error() == "scheduler: database not available" {
			return false, nil
		}
		return false, err
	}
	if len(recentRuns) == 0 {
		return false, nil
	}

	lastRunTime, err := time.Parse(time.RFC3339, recentRuns[0].StartedAt)
	if err != nil {
		return false, fmt.Errorf("parse started_at %q: %w", recentRuns[0].StartedAt, err)
	}

	// Already fired if the most recent run's started_at is on or after
	// the scheduled time.
	return !lastRunTime.Before(scheduleTime), nil
}

// addJobLocked registers a job with the cron scheduler.
// Must be called with s.mu held.
func (s *Scheduler) addJobLocked(job *Job) {
	jobCopy := *job // capture by value
	entryID, err := s.cron.AddFunc(job.Cron, func() {
		log.Printf("scheduler: firing job %q (%s)", jobCopy.ID, jobCopy.Name)
		s.runner.Run(&jobCopy, newRunID())
	})
	if err != nil {
		log.Printf("scheduler: invalid cron expression for job %q: %v", job.ID, err)
		return
	}
	s.entries[job.ID] = entryID
}
