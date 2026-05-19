package scheduler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

// Handler provides HTTP handlers for the scheduler REST API.
type Handler struct {
	store *Store
}

// NewHandler creates a new Handler backed by the given Store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
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
