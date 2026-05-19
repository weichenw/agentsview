package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

const (
	maxOutputBytes = 10 * 1024 * 1024 // 10MB cap
	runTimeout     = 30 * time.Minute
)

// Runner spawns Pi processes for scheduled jobs.
type Runner struct {
	store *Store
}

// NewRunner creates a Runner that records run history to store.
func NewRunner(store *Store) *Runner {
	return &Runner{store: store}
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
func (r *Runner) Run(job *Job) *RunResult {
	runID := newRunID()
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

	result := &RunResult{RunID: runID}

	// Build command.
	args := []string{}
	if !job.InheritProjectContext {
		args = append(args, "--no-project-context")
	}
	if job.Model != "" {
		args = append(args, "--model", job.Model)
	}
	args = append(args, job.Prompt)

	cmd := exec.Command("pi", args...)
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

		status := "completed"
		if exitCode != 0 {
			status = "failed"
		}

		errMsg := ""
		if exitCode != 0 {
			errMsg = strings.TrimSpace(stderr.String())
			if errMsg == "" {
				errMsg = fmt.Sprintf("exit code %d", exitCode)
			}
		}

		// Try to extract session ID from stdout.
		sessionID := extractSessionID(stdout.String())

		result.Status = status
		result.ExitCode = exitCode
		result.Error = errMsg
		result.SessionID = sessionID

		if err := r.store.UpdateRun(runID, status, exitCode, errMsg); err != nil {
			log.Printf("scheduler: update run %s: %v", runID, err)
		}

	case <-ctx.Done():
		// Timeout: kill the process.
		if cmd.Process != nil {
			cmd.Process.Kill()
			// Wait for it to finish so we don't leak the goroutine.
			<-done
		}
		result.Status = "killed"
		result.ExitCode = -1
		result.Error = "process timed out after 30m"

		if err := r.store.UpdateRun(runID, "killed", -1, "process timed out after 30m"); err != nil {
			log.Printf("scheduler: update run %s: %v", runID, err)
		}
	}

	log.Printf("scheduler: job %q run %s: %s (exit=%d)", job.ID, runID, result.Status, result.ExitCode)
	return result
}

// newRunID generates a short unique run ID.
func newRunID() string {
	ts := time.Now().UnixNano()
	return fmt.Sprintf("run_%x", ts)
}

// extractSessionID attempts to find a session ID from pi's stdout.
// Looks for a line containing "session_id" or "Session ID" emitted by pi.
func extractSessionID(output string) string {
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		// pi outputs session info on stderr typically, but check
		// stdout too as a fallback.
		if strings.Contains(line, "session_id") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				id := strings.TrimSpace(parts[len(parts)-1])
				if id != "" {
					return id
				}
			}
		}
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
