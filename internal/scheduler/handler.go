package scheduler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"
)

// Handler provides HTTP handlers for the scheduler REST API.
type Handler struct {
	store  *Store
	runner *Runner
	sched  *Scheduler
}

// NewHandler creates a new Handler backed by the given Store.
func NewHandler(store *Store, runner *Runner, sched *Scheduler) *Handler {
	return &Handler{store: store, runner: runner, sched: sched}
}

var validIDPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9_-]*[a-z0-9])?$`)

// ListJobs handles GET /api/v1/scheduler/jobs
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	jobs := h.store.List()
	if jobs == nil {
		jobs = []Job{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"jobs": jobs})
}

// GetJob handles GET /api/v1/scheduler/jobs/{id}
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job := h.store.Get(id)
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

// CreateJob handles POST /api/v1/scheduler/jobs
func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID                   string `json:"id"`
		Name                 string `json:"name"`
		Cron                 string `json:"cron"`
		Enabled              *bool  `json:"enabled"`
		Agent                string `json:"agent"`
		Prompt               string `json:"prompt"`
		Model                string `json:"model,omitempty"`
		InheritProjectContext bool   `json:"inherit_project_context"`
		WorkingDir           string `json:"working_dir"`
		SpawnMode            string `json:"spawn_mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if input.ID == "" {
		// Generate a short unique ID if not provided.
		b := make([]byte, 4)
		rand.Read(b)
		input.ID = hex.EncodeToString(b)
	}
	if err := validateJobInput(input.ID, input.Name, input.Cron, input.Agent, input.Prompt, input.SpawnMode); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	if input.SpawnMode == "" {
		input.SpawnMode = "cmux"
	}

	job := &Job{
		ID:                   input.ID,
		Name:                 input.Name,
		Cron:                 input.Cron,
		Enabled:              enabled,
		Agent:                input.Agent,
		Prompt:               input.Prompt,
		Model:                input.Model,
		InheritProjectContext: input.InheritProjectContext,
		WorkingDir:           input.WorkingDir,
		SpawnMode:            input.SpawnMode,
	}

	if err := h.store.Create(job); err != nil {
		if errors.Is(err, ErrJobExists) {
			writeError(w, http.StatusConflict, "job already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, h.store.Get(job.ID))

	// Reload cron entry for the new job.
	if h.sched != nil {
		h.sched.Reload(job.ID)
	}
}

// UpdateJob handles PUT /api/v1/scheduler/jobs/{id}
func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := h.store.Get(id)
	if existing == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}

	var input struct {
		Name                 string `json:"name"`
		Cron                 string `json:"cron"`
		Enabled              *bool  `json:"enabled"`
		Agent                string `json:"agent"`
		Prompt               string `json:"prompt"`
		Model                string `json:"model,omitempty"`
		InheritProjectContext *bool  `json:"inherit_project_context"`
		WorkingDir           string `json:"working_dir"`
		SpawnMode            string `json:"spawn_mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	// Merge: use existing values for omitted fields.
	job := *existing
	if input.Name != "" {
		job.Name = input.Name
	}
	if input.Cron != "" {
		job.Cron = input.Cron
	}
	if input.Enabled != nil {
		job.Enabled = *input.Enabled
	}
	if input.Agent != "" {
		job.Agent = input.Agent
	}
	if input.Prompt != "" {
		job.Prompt = input.Prompt
	}
	if input.Model != "" {
		job.Model = input.Model
	}
	if input.InheritProjectContext != nil {
		job.InheritProjectContext = *input.InheritProjectContext
	}
	if input.WorkingDir != "" {
		job.WorkingDir = input.WorkingDir
	}
	if input.SpawnMode != "" {
		job.SpawnMode = input.SpawnMode
	}

	if err := validateJobInput(job.ID, job.Name, job.Cron, job.Agent, job.Prompt, job.SpawnMode); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.Update(id, &job); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, h.store.Get(id))

	// Reload cron entry with updated job.
	if h.sched != nil {
		h.sched.Reload(id)
	}
}

// DeleteJob handles DELETE /api/v1/scheduler/jobs/{id}
func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := h.store.Get(id)
	if existing == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if err := h.store.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// Remove cron entry for deleted job.
	if h.sched != nil {
		h.sched.Reload(id)
	}
}

// EnableJob handles POST /api/v1/scheduler/jobs/{id}/enable
func (h *Handler) EnableJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.setEnabled(w, id, true)
}

// DisableJob handles POST /api/v1/scheduler/jobs/{id}/disable
func (h *Handler) DisableJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.setEnabled(w, id, false)
}

// RunJob handles POST /api/v1/scheduler/jobs/{id}/run
// Triggers immediate execution and returns 202 Accepted with run ID.
func (h *Handler) RunJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job := h.store.Get(id)
	if job == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if h.runner == nil {
		writeError(w, http.StatusInternalServerError, "runner not available")
		return
	}

	// Generate run ID and create the run entry synchronously
	// so the response always includes a valid run_id.
	runID := newRunID()
	startedAt := time.Now().UTC().Format(time.RFC3339)
	run := &SchedulerRun{
		ID:        runID,
		JobID:     job.ID,
		StartedAt: startedAt,
		Status:    "running",
	}
	if err := h.store.CreateRun(run); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create run entry")
		return
	}

	// Return the run_id immediately before the (possibly slow) subprocess.
	writeJSON(w, http.StatusAccepted, map[string]any{
		"run_id": runID,
		"status": "launched",
		"job_id": id,
	})

	// Fire the actual run in a background goroutine.
	go h.runner.Run(job, runID)
}

// KillRun handles POST /api/v1/scheduler/runs/{id}/kill
// Kills a running subprocess for the given run ID.
func (h *Handler) KillRun(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("id")
	if h.runner == nil {
		writeError(w, http.StatusInternalServerError, "runner not available")
		return
	}
	if err := h.runner.KillRun(runID); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	// Update the run entry in the database.
	_ = h.store.UpdateRun(runID, "killed", -1, "killed by user")
	w.WriteHeader(http.StatusNoContent)
}

// ListRuns handles GET /api/v1/scheduler/runs
func (h *Handler) ListRuns(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	jobID := q.Get("job_id")
	limit := 20
	if v := q.Get("limit"); v != "" {
		if n, err := parseInt(v); err == nil && n > 0 {
			limit = n
		}
	}

	runs, err := h.store.ListRuns(jobID, limit)
	if err != nil {
		log.Printf("scheduler: ListRuns jobID=%q limit=%d err=%v", jobID, limit, err)
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("scheduler: ListRuns jobID=%q limit=%d runs=%d", jobID, limit, len(runs))
	writeJSON(w, http.StatusOK, map[string]any{"runs": runs})
}

func parseInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number")
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func (h *Handler) setEnabled(w http.ResponseWriter, id string, enabled bool) {
	existing := h.store.Get(id)
	if existing == nil {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	existing.Enabled = enabled
	if err := h.store.Update(id, existing); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, h.store.Get(id))

	// Reload cron entry with updated enabled state.
	if h.sched != nil {
		h.sched.Reload(id)
	}
}

// GetSettings handles GET /api/v1/scheduler/settings
func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.store.LoadSettings()
	if err != nil {
		log.Printf("scheduler: LoadSettings: %v", err)
		writeJSON(w, http.StatusOK, &SchedulerSettings{Timezone: defaultTimezone})
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

// UpdateSettings handles POST /api/v1/scheduler/settings
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings SchedulerSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if settings.Timezone == "" {
		writeError(w, http.StatusBadRequest, "timezone is required")
		return
	}
	// Validate timezone by attempting to load it.
	if settings.Timezone != "Local" {
		if _, err := time.LoadLocation(settings.Timezone); err != nil {
			writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid timezone: %v", err))
			return
		}
	}
	if err := h.store.SaveSettings(&settings); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, settings)

	// Restart the scheduler with the new timezone.
	if h.sched != nil {
		if err := h.sched.SetTimezone(settings.Timezone); err != nil {
			log.Printf("scheduler: SetTimezone %q: %v", settings.Timezone, err)
		}
	}
}

// Health handles GET /api/v1/scheduler/health
// Returns lightweight scheduler status for the status-bar heartbeat indicator.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	settings, _ := h.store.LoadSettings()
	if settings == nil {
		settings = &SchedulerSettings{Timezone: defaultTimezone}
	}

	activeJobs := 0
	for _, job := range h.store.List() {
		if job.Enabled {
			activeJobs++
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"healthy":      h.sched != nil,
		"timezone":     settings.Timezone,
		"active_jobs":  activeJobs,
	})
}

func validateJobInput(id, name, cron, agent, prompt, spawnMode string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if !validIDPattern.MatchString(id) {
		return errors.New("id must match [a-z0-9]([a-z0-9_-]*[a-z0-9])?")
	}
	if name == "" {
		return errors.New("name is required")
	}
	if cron == "" {
		return errors.New("cron is required")
	}
	if agent == "" {
		return errors.New("agent is required")
	}
	if prompt == "" {
		return errors.New("prompt is required")
	}
	if spawnMode != "" && spawnMode != "cmux" && spawnMode != "subprocess" {
		return errors.New("spawn_mode must be \"cmux\" or \"subprocess\"")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
