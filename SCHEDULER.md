# Scheduler Feature Spec

## Overview

Add a Scheduler tab to agentsview that lets you define, manage, and monitor
scheduled Pi agent jobs. The scheduler fires `pi` commands via cmux or directly
as subprocesses. Session history for each job appears automatically in the
existing Sessions view via the existing SSE/SQLite pipeline.

---

## What We Are Building

1. **Backend** — Go scheduler service that reads `schedules.json`, runs jobs on
   cron triggers, spawns Pi processes, and exposes CRUD + run-history REST endpoints.
2. **Frontend** — New "Scheduler" tab in the React UI with a job list, create/edit
   form, run history per job, and live status feed.
3. **Storage** — `~/.agentsview/schedules.json` for job definitions.
   Run history reuses the existing SQLite sessions database — no new tables needed.

---

## Non-Goals

- No inter-agent coordination (no shared task lists, no subagent spawning from here)
- No modification to Pi internals or any Pi extension
- No Docker-specific changes — macOS launchd is the OS scheduler;
  agentsview scheduler is the *management UI*, not a replacement for the OS

---

## Data Model

### Schedule definition (`~/.agentsview/schedules.json`)

```json
[
  {
    "id": "morning-briefer",
    "name": "Morning Portfolio Briefing",
    "cron": "0 7 * * *",
    "enabled": true,
    "agent": "morning-briefer",
    "prompt": "Run the morning portfolio briefing.",
    "model": "anthropic/claude-sonnet-4",
    "inherit_project_context": false,
    "working_dir": "/Users/Wei/projects/finserv",
    "spawn_mode": "cmux",
    "created_at": "2026-05-19T00:00:00Z",
    "updated_at": "2026-05-19T00:00:00Z"
  }
]
```

**Fields:**

| Field | Type | Notes |
|---|---|---|
| id | string | slug, unique |
| name | string | display name |
| cron | string | standard 5-field cron expression |
| enabled | bool | pause without deleting |
| agent | string | Pi agent name (maps to .pi/agents/*.md) |
| prompt | string | prompt injected at runtime |
| model | string | optional model override |
| inherit_project_context | bool | default false |
| working_dir | string | cwd for the pi process |
| spawn_mode | string | "cmux" or "subprocess" |
| created_at | string | RFC3339 |
| updated_at | string | RFC3339 |

---

## spawn_mode

### cmux (default, recommended)

```bash
cmux new-window -n "{id}" \
  "cd {working_dir} && pi --no-project-context --model {model} '{prompt}'"
```

- Visible pane — you can watch, scroll, kill it
- Survives agentsview restart
- Works with HazAT subagents if Pi session is running

### subprocess

```go
cmd := exec.Command("pi", "--no-project-context", "--model", model, prompt)
cmd.Dir = workingDir
```

- Headless, no terminal pane
- Result captured to log file
- Use when cmux is not available (Docker)

---

## Backend

### New files

```
internal/
  scheduler/
    scheduler.go       ← cron runner, job lifecycle
    store.go           ← read/write schedules.json
    runner.go          ← spawn pi via cmux or subprocess
    handler.go         ← HTTP handlers for REST API
```

### REST API

```
GET    /api/v1/scheduler/jobs              list all jobs
POST   /api/v1/scheduler/jobs              create job
GET    /api/v1/scheduler/jobs/{id}         get job
PUT    /api/v1/scheduler/jobs/{id}         update job
DELETE /api/v1/scheduler/jobs/{id}         delete job
POST   /api/v1/scheduler/jobs/{id}/run     trigger immediately
POST   /api/v1/scheduler/jobs/{id}/enable  enable
POST   /api/v1/scheduler/jobs/{id}/disable disable

GET    /api/v1/scheduler/runs              recent runs (all jobs)
GET    /api/v1/scheduler/runs/{job_id}     runs for a specific job
```

### Run history

Each run entry stored in SQLite alongside existing sessions:

```sql
CREATE TABLE scheduler_runs (
  id          TEXT PRIMARY KEY,
  job_id      TEXT NOT NULL,
  session_id  TEXT,           -- links to existing sessions table
  started_at  TEXT NOT NULL,
  finished_at TEXT,
  status      TEXT NOT NULL,  -- running | completed | failed | killed
  exit_code   INTEGER,
  error       TEXT
);
```

`session_id` links to the existing `sessions` table — so the full session
replay, token cost, and health score are all available for free via the
existing session API.

### Cron runner (scheduler.go)

```go
// pseudocode
func (s *Scheduler) Start() {
    c := cron.New()
    for _, job := range s.store.List() {
        if !job.Enabled { continue }
        c.AddFunc(job.Cron, func() {
            s.runner.Run(job)
        })
    }
    c.Start()
}
```

Use `github.com/robfig/cron/v3` — already a common Go cron library,
check go.mod before adding.

---

## Frontend

### New tab

Add "Scheduler" between "Sessions" and "Usage" in the nav.

### Job list view

```
┌─ Scheduler ──────────────────────────────────────────────────────┐
│  + New Job                                        [Search jobs]  │
│                                                                   │
│  ● Morning Briefing        0 7 * * *   next: tomorrow 07:00      │
│    last run: completed · 42s · 18k tok · 2 hours ago             │
│                                                                   │
│  ○ Code Review (disabled)  0 14 * * 1  —                         │
│    last run: failed · 3 days ago                                  │
└──────────────────────────────────────────────────────────────────┘
```

### Job detail / edit form

Fields:
- Name
- Cron expression (with human-readable preview: "Every day at 7:00 AM")
- Agent name
- Prompt (multiline)
- Model override
- Working directory
- Spawn mode (cmux / subprocess)
- Inherit project context toggle
- Enabled toggle

### Run history panel

Per-job run list linking directly to the session viewer for full replay:

```
Run history — Morning Briefing
──────────────────────────────────────────────────────
2026-05-19 07:00  completed  44s  22k tok  A  → view session
2026-05-18 07:00  completed  38s  19k tok  A  → view session
2026-05-17 07:00  failed     12s  4k tok   F  → view session
```

Clicking "view session" opens the existing session detail view.

---

## Implementation Order

Do these in order. Each is a shippable increment.

### Phase 1 — Storage + API (no UI)

1. `internal/scheduler/store.go` — read/write `schedules.json`
2. `internal/scheduler/handler.go` — REST CRUD endpoints
3. Wire handlers into existing server setup in `cmd/serve.go`
4. Test: `curl /api/v1/scheduler/jobs` returns empty list

### Phase 2 — Runner (no cron yet)

1. `internal/scheduler/runner.go` — subprocess spawn only first
2. `POST /api/v1/scheduler/jobs/{id}/run` triggers a job immediately
3. Write run entry to SQLite `scheduler_runs` table
4. Test: trigger morning-briefer manually, session appears in Sessions tab

### Phase 3 — Cron

1. `internal/scheduler/scheduler.go` — cron runner using robfig/cron
2. Start scheduler in `cmd/serve.go` alongside existing server start
3. Test: set a job to fire every minute, confirm run appears in history

### Phase 4 — cmux spawn mode

1. Add cmux spawn path to `runner.go`
2. Detect if cmux is available, fall back to subprocess
3. Test: job opens a new cmux window

### Phase 5 — Frontend

1. Scheduler tab scaffold (empty page, nav item)
2. Job list component
3. Create/edit form
4. Run history panel with session links
5. Live status via existing SSE endpoint

---

## Constraints

- Never modify existing session parsing, SQLite schema (additive only),
  or SSE infrastructure — the scheduler plugs into them, does not replace them
- Keep `schedules.json` human-editable — no binary formats
- cmux spawn must not block the agentsview server process
- subprocess spawn must cap output at 10MB to avoid disk fill
- All new API endpoints under `/api/v1/scheduler/` prefix
- Follow existing handler patterns in `internal/` — look at how
  session handlers are structured before writing new ones

---

## Files to Read Before Starting

Read these files first to understand existing patterns:

```
AGENTS.md                          ← project conventions
internal/server/server.go          ← how handlers are registered
internal/sessions/handler.go       ← example handler pattern
internal/sessions/store.go         ← example store pattern  
internal/db/schema.go              ← existing SQLite schema
frontend/src/components/nav/       ← how tabs are added
frontend/src/pages/                ← existing page structure
go.mod                             ← check available dependencies
```

---

## Definition of Done

- [ ] `GET /api/v1/scheduler/jobs` returns job list
- [ ] `POST /api/v1/scheduler/jobs/{id}/run` fires Pi, session appears in Sessions tab
- [ ] Cron fires morning-briefer at 7am, run recorded in `scheduler_runs`
- [ ] Scheduler tab shows job list with last run status
- [ ] Create/edit form saves to `schedules.json`
- [ ] Run history links to session detail view
- [ ] cmux spawn mode opens a new window
- [ ] No existing tests broken
