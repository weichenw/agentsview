package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	maxOutputBytes = 10 * 1024 * 1024 // 10MB cap
	runTimeout     = 30 * time.Minute
)

// Runner spawns Pi processes for scheduled jobs.
type Runner struct {
	store       *Store
	postRunHook func()
}

// NewRunner creates a Runner that records run history to store.
// postRunHook is called after each subprocess run completes (success or failure).
func NewRunner(store *Store, postRunHook func()) *Runner {
	return &Runner{store: store, postRunHook: postRunHook}
}

// RunResult holds the outcome of a single job execution.
type RunResult struct {
	RunID     string
	SessionID string
	Status    string
	ExitCode  int
	Error     string
}

// Run executes a job immediately as a subprocess.
// It creates a run entry, spawns the process, waits for completion,
// and updates the run entry with the result.
func (r *Runner) Run(job *Job, runID string) *RunResult {
	startedAt := time.Now().UTC().Format(time.RFC3339)

	// Create initial run entry.
	run := &SchedulerRun{
		ID:        runID,
		JobID:     job.ID,
		StartedAt: startedAt,
		Status:    "running",
	}
	if err := r.store.CreateRun(run); err != nil {
		log.Printf("scheduler: create run entry: %v", err)
	}

	// Check for cmux spawn mode.
	if job.SpawnMode == "cmux" {
		r.runCmux(job, run)
		log.Printf("scheduler: job %q run %s launched via cmux", job.ID, runID)
		return &RunResult{RunID: runID, Status: "running"}
	}

	// Default: subprocess mode (synchronous, blocks until complete).
	r.runSubprocess(job, run)
	log.Printf("scheduler: job %q run %s finished", job.ID, runID)
	return &RunResult{RunID: runID, Status: run.Status}
}

// runCmux spawns a job via cmux, creating a visible terminal window.
// It writes a run entry to the database but cannot track completion
// since the pi process runs in a separate window.
func (r *Runner) runCmux(job *Job, run *SchedulerRun) {
	// Check if cmux is available.
	cmuxPath, err := exec.LookPath("cmux")
	if err != nil {
		log.Printf("scheduler: cmux not found, falling back to subprocess for job %q", job.ID)
		r.runSubprocess(job, run)
		return
	}

	// Build the command string for the new window.
	// cmux new-window -n "{id}" "cd {dir} && pi ..."
	piArgs := []string{}
	if !job.InheritProjectContext {
		piArgs = append(piArgs, "--no-context-files")
	}
	sessionDir := filepath.Join(os.Getenv("HOME"), ".pi", "agent", "sessions", "scheduler")
	piArgs = append(piArgs, "--session-dir", sessionDir)
	if job.Model != "" {
		piArgs = append(piArgs, "--model", job.Model)
	}
	piArgs = append(piArgs, job.Prompt)

	piCmd := "pi"
	for _, a := range piArgs {
		piCmd += " " + a
	}

	windowCmd := piCmd
	if job.WorkingDir != "" {
		windowCmd = "cd " + job.WorkingDir + " && " + piCmd
	}

	cmd := exec.Command(cmuxPath, "new-window", "-n", job.ID, windowCmd)
	if err := cmd.Start(); err != nil {
		log.Printf("scheduler: cmux spawn failed for job %q: %v, falling back to subprocess", job.ID, err)
		// Update run entry to record the failure before fallback.
		run.Status = "failed"
		if uErr := r.store.UpdateRun(run.ID, "failed", -1, err.Error()); uErr != nil {
			log.Printf("scheduler: update run %s: %v", run.ID, uErr)
		}
		r.runSubprocess(job, run)
		return
	}

	// cmux spawn succeeded — fire-and-forget.
	// The pi session will be picked up by the existing sync pipeline.
	log.Printf("scheduler: cmux window launched for job %q", job.ID)
}

// runSubprocess spawns a job as a subprocess with output capture,
// timeout handling, and run history tracking.
func (r *Runner) runSubprocess(job *Job, run *SchedulerRun) {
	// Build command.
	args := []string{}
	if !job.InheritProjectContext {
		args = append(args, "--no-context-files")
	}
	if job.Model != "" {
		args = append(args, "--model", job.Model)
	}
	sessionDir := filepath.Join(os.Getenv("HOME"), ".pi", "agent", "sessions", "scheduler")
	args = append(args, "--session-dir", sessionDir)
	args = append(args, job.Prompt)

	piPath, err := exec.LookPath("pi")
	if err != nil {
		// fallback to known homebrew location
		piPath = "/opt/homebrew/bin/pi"
	}
	cmd := exec.Command(piPath, args...)
	if job.WorkingDir != "" {
		cmd.Dir = job.WorkingDir
	}

	// Cap output at 10MB.
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &maxSizeWriter{w: &stdout, max: maxOutputBytes}
	cmd.Stderr = &maxSizeWriter{w: &stderr, max: maxOutputBytes}

	// Run with 30-minute timeout.
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				exitCode = -1
			}
		}

		run.Status = "completed"
		if exitCode != 0 {
			run.Status = "failed"
		}

		errMsg := ""
		if exitCode != 0 {
			errMsg = strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = fmt.Sprintf("exit code %d", exitCode)
			}
		}

		if err := r.store.UpdateRun(run.ID, run.Status, exitCode, errMsg); err != nil {
			log.Printf("scheduler: update run %s: %v", run.ID, err)
		}

	case <-ctx.Done():
		// Timeout: kill the process.
		if cmd.Process != nil {
			cmd.Process.Kill()
			<-done
		}
		run.Status = "killed"
		if err := r.store.UpdateRun(run.ID, "killed", -1, "process timed out after 30m"); err != nil {
			log.Printf("scheduler: update run %s: %v", run.ID, err)
		}
	}

	if r.postRunHook != nil {
		r.postRunHook()
	}

	// Find the session file created by this run and link it.
	startedAt, parseErr := time.Parse(time.RFC3339, run.StartedAt)
	if parseErr == nil {
		if sessionID := findLatestSessionID(startedAt); sessionID != "" {
			run.SessionID = sessionID
			// Update the session_id in the database.
			if _, dbErr := r.store.db.Writer().Exec(
				`UPDATE scheduler_runs SET session_id = ? WHERE id = ?`,
				sessionID, run.ID,
			); dbErr != nil {
				log.Printf("scheduler: update session_id for run %s: %v", run.ID, dbErr)
			}
		}
	}

	log.Printf("scheduler: subprocess job %q run %s finished", job.ID, run.ID)
}

// newRunID generates a short unique run ID.
func newRunID() string {
	ts := time.Now().UnixNano()
	return fmt.Sprintf("run_%x", ts)
}

// findLatestSessionID scans the default pi sessions directory for the
// most recently modified .jsonl file created after startedAt, and
// returns the session UUID extracted from its filename.
func findLatestSessionID(startedAt time.Time) string {
	sessionsDir := filepath.Join(os.Getenv("HOME"), ".pi", "agent", "sessions")
	var latestPath string
	var latestMod time.Time

	filepath.WalkDir(sessionsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".jsonl") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().After(startedAt) && info.ModTime().After(latestMod) {
			latestMod = info.ModTime()
			latestPath = path
		}
		return nil
	})

	if latestPath == "" {
		return ""
	}

	// Extract UUID from filename (format: session_<uuid>.jsonl)
	base := filepath.Base(latestPath)
	parts := strings.Split(base, "_")
	if len(parts) >= 2 {
		return strings.TrimSuffix(parts[len(parts)-1], ".jsonl")
	}
	return ""
}

// maxSizeWriter wraps a bytes.Buffer but stops writing after max bytes.
type maxSizeWriter struct {
	w    *bytes.Buffer
	max  int
	done bool
}

func (mw *maxSizeWriter) Write(p []byte) (int, error) {
	if mw.done {
		return len(p), nil // silently drop
	}
	remaining := mw.max - mw.w.Len()
	if remaining <= 0 {
		mw.done = true
		return len(p), nil
	}
	if len(p) > remaining {
		p = p[:remaining]
	}
	return mw.w.Write(p)
}
