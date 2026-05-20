package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Scheduler manages cron-driven job execution.
type Scheduler struct {
	mu      sync.Mutex
	store   *Store
	runner  *Runner
	cron    *cron.Cron
	entries map[string]cron.EntryID // job ID -> cron entry ID
}

// New creates a Scheduler. Call Start() to begin execution.
func New(store *Store, runner *Runner) *Scheduler {
	return &Scheduler{
		store:   store,
		runner:  runner,
		cron:    cron.New(cron.WithSeconds()),
		entries: make(map[string]cron.EntryID),
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
	log.Printf("scheduler: started with %d active job(s)", len(s.entries))
}

// Stop gracefully shuts down the cron runner, waiting for any
// running jobs to finish.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx := s.cron.Stop() // returns context that's done when all jobs finish
	<-ctx.Done()
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
