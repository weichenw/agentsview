package scheduler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"go.kenn.io/agentsview/internal/db"
)

// ErrJobExists is returned by Create when a job ID already exists.
var ErrJobExists = errors.New("job already exists")

// Job defines a scheduled Pi agent job.
type Job struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Cron                  string `json:"cron"`
	Enabled               bool   `json:"enabled"`
	Agent                 string `json:"agent"`
	Prompt                string `json:"prompt"`
	Model                 string `json:"model,omitempty"`
	InheritProjectContext bool   `json:"inherit_project_context"`
	WorkingDir            string `json:"working_dir"`
	SpawnMode             string `json:"spawn_mode"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`
}

// SchedulerRun represents a single execution of a scheduled job.
type SchedulerRun struct {
	ID         string `json:"id"`
	JobID      string `json:"job_id"`
	SessionID  string `json:"session_id,omitempty"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at,omitempty"`
	Status     string `json:"status"`
	ExitCode   *int   `json:"exit_code,omitempty"`
	Error      string `json:"error,omitempty"`
}

// Store reads and writes schedules.json and run history.
type Store struct {
	mu       sync.RWMutex
	filePath string
	jobs     map[string]*Job
	db       *db.DB
}

// NewStore creates a Store rooted at dataDir.
// If dataDir is empty it defaults to $HOME/.agentsview.
// Pass a non-nil *db.DB to enable run history persistence.
func NewStore(dataDir string, database *db.DB) (*Store, error) {
	if dataDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("scheduler: cannot determine home dir: %w", err)
		}
		dataDir = filepath.Join(home, ".agentsview")
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		return nil, fmt.Errorf("scheduler: create data dir: %w", err)
	}
	s := &Store{
		filePath: filepath.Join(dataDir, "schedules.json"),
		jobs:     make(map[string]*Job),
		db:       database,
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

// List returns all jobs sorted by created_at descending.
func (s *Store) List() []Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		result = append(result, *j)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})
	return result
}

// Get returns a job by ID, or nil if not found.
func (s *Store) Get(id string) *Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil
	}
	copy := *j
	return &copy
}

// Create adds a new job. The job's ID must be unique.
func (s *Store) Create(job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.ID]; exists {
		return fmt.Errorf("scheduler: job %q already exists: %w", job.ID, ErrJobExists)
	}
	now := time.Now().UTC().Format(time.RFC3339)
	job.CreatedAt = now
	job.UpdatedAt = now
	copy := *job
	s.jobs[job.ID] = &copy
	return s.save()
}

// Update replaces an existing job identified by id.
func (s *Store) Update(id string, job *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return fmt.Errorf("scheduler: job %q not found", id)
	}
	job.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	job.CreatedAt = s.jobs[id].CreatedAt
	job.ID = id
	copy := *job
	s.jobs[id] = &copy
	return s.save()
}

// Delete removes a job by ID.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return fmt.Errorf("scheduler: job %q not found", id)
	}
	delete(s.jobs, id)
	return s.save()
}

// load reads schedules.json from disk.
func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // empty store
		}
		return fmt.Errorf("scheduler: read schedules: %w", err)
	}
	var jobs []Job
	if err := json.Unmarshal(data, &jobs); err != nil {
		return fmt.Errorf("scheduler: parse schedules: %w", err)
	}
	for i := range jobs {
		j := jobs[i]
		s.jobs[j.ID] = &j
	}
	return nil
}

// save atomically writes schedules.json to disk.
func (s *Store) save() error {
	jobs := make([]Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		jobs = append(jobs, *j)
	}
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return fmt.Errorf("scheduler: marshal schedules: %w", err)
	}
	// Atomic write: write to temp file, then rename.
	tmpPath := s.filePath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o600); err != nil {
		return fmt.Errorf("scheduler: write schedules temp: %w", err)
	}
	if err := os.Rename(tmpPath, s.filePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("scheduler: rename schedules: %w", err)
	}
	return nil
}

// CreateRun inserts a new scheduler_runs row.
func (s *Store) CreateRun(run *SchedulerRun) error {
	if s.db == nil {
		return fmt.Errorf("scheduler: database not available")
	}
	_, err := s.db.Writer().Exec(
		`INSERT INTO scheduler_runs (id, job_id, session_id, started_at, status)
		VALUES (?, ?, ?, ?, ?)`,
		run.ID, run.JobID, run.SessionID, run.StartedAt, run.Status,
	)
	if err != nil {
		return fmt.Errorf("scheduler: insert run: %w", err)
	}
	return nil
}

// UpdateRun updates a scheduler_runs row with completion info.
func (s *Store) UpdateRun(id, status string, exitCode int, errMsg string) error {
	if s.db == nil {
		return fmt.Errorf("scheduler: database not available")
	}
	finishedAt := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Writer().Exec(
		`UPDATE scheduler_runs
		SET status = ?, finished_at = ?, exit_code = ?, error = ?
		WHERE id = ?`,
		status, finishedAt, exitCode, errMsg, id,
	)
	if err != nil {
		return fmt.Errorf("scheduler: update run: %w", err)
	}
	return nil
}

// ListRuns returns recent runs for a job, newest first.
func (s *Store) ListRuns(jobID string, limit int) ([]SchedulerRun, error) {
	if s.db == nil {
		return nil, fmt.Errorf("scheduler: database not available")
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Reader().Query(
		`SELECT id, job_id, COALESCE(session_id, ''), started_at,
		COALESCE(finished_at, ''), status, exit_code, COALESCE(error, '')
		FROM scheduler_runs
		WHERE job_id = ?
		ORDER BY started_at DESC
		LIMIT ?`,
		jobID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("scheduler: list runs: %w", err)
	}
	defer rows.Close()

	var runs []SchedulerRun
	for rows.Next() {
		var r SchedulerRun
		if err := rows.Scan(&r.ID, &r.JobID, &r.SessionID, &r.StartedAt,
			&r.FinishedAt, &r.Status, &r.ExitCode, &r.Error); err != nil {
			return nil, fmt.Errorf("scheduler: scan run: %w", err)
		}
		runs = append(runs, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scheduler: rows iteration: %w", err)
	}
	if runs == nil {
		runs = []SchedulerRun{}
	}
	return runs, nil
}

// LoadSettings reads scheduler-settings.json from the data dir.
// Returns a default settings struct if the file does not exist.
func (s *Store) LoadSettings() (*SchedulerSettings, error) {
	return LoadSettings(filepath.Dir(s.filePath))
}

// SaveSettings atomically writes scheduler-settings.json to the data dir.
func (s *Store) SaveSettings(settings *SchedulerSettings) error {
	return SaveSettings(filepath.Dir(s.filePath), settings)
}
