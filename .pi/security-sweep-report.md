# Security Sweep Report

**Project:** agents-view
**Root:** `/Users/Wei/Workspace/Agents/agents-view`
**Languages:** Go
**Frameworks:** None detected
**Files Scanned:** 604
**Date:** 2026-05-25T13:12:04.323Z

## Security Score: 🔴 0/100 (Critical)

## AI Services Detected

- **Google AI / Gemini** (`@google/generative-ai`)
  - `.roborev.toml`
  - `cmd/agentsview/cli.go`
  - `cmd/agentsview/session_get.go`
  - `cmd/testfixture/main.go`
  - `frontend/src/lib/api/types/insights.ts`
  - ... and 41 more
- **Anthropic** (`@anthropic-ai/sdk`)
  - `internal/db/analytics.go`
  - `internal/db/analytics_test.go`
  - `internal/insight/generate.go`
  - `internal/insight/generate_test.go`

## Findings Summary

| Severity | Count |
|----------|-------|
| 🔴 Critical | 57 |
| 🟠 High | 137 |
| 🟡 Medium | 2 |
| 🔵 Low | 0 |

## Findings

### SEC-001: 🔴 AI output passed to eval/Function constructor

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `desktop/src-tauri/src/lib.rs:92`
- **Description:** AI-generated output is being evaluated as code via eval() or the Function constructor. This is extremely dangerous — a prompt injection could lead to arbitrary code execution.
- **Evidence:**
  ```
                      if let Some(window) = app_handle.get_webview_window("main") {
                        let _ = window.eval("window.dispatchEvent(new CustomEvent('show-about'));");
                    }
                }
  ```
- **Recommendation:** Never eval() AI outputs. Use structured output parsing (JSON.parse with validation) instead.

### SEC-002: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `desktop/src-tauri/src/lib.rs:97`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
                      let handle = app_handle.clone();
                    tauri::async_runtime::spawn(async move {
                        check_for_updates(&handle, false).await;
                    });
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-003: 🔴 AI output passed to eval/Function constructor

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `frontend/public/d3.v7.min.js:2`
- **Description:** AI-generated output is being evaluated as code via eval() or the Function constructor. This is extremely dangerous — a prompt injection could lead to arbitrary code execution.
- **Evidence:**
  ```
  // https://d3js.org v7.9.0 Copyright 2010-2023 Mike Bostock
!function(t,n){"object"==typeof exports&&"undefined"!=typeof module?n(exports):"function"==typeof define&&define.amd?define(["exports"],n):n((t="undefined"!=typeof globalThis?globalThis:t||self).d3=t.d3||{})}(this,(function(t){"use strict";...
  ```
- **Recommendation:** Never eval() AI outputs. Use structured output parsing (JSON.parse with validation) instead.

### SEC-004: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `frontend/public/d3.v7.min.js:2`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  // https://d3js.org v7.9.0 Copyright 2010-2023 Mike Bostock
!function(t,n){"object"==typeof exports&&"undefined"!=typeof module?n(exports):"function"==typeof define&&define.amd?define(["exports"],n):n((t="undefined"!=typeof globalThis?globalThis:t||self).d3=t.d3||{})}(this,(function(t){"use strict";...
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-005: 🔴 AI output passed to eval/Function constructor

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `frontend/src/lib/components/layout/ThreeColumnLayout.test.ts:114`
- **Description:** AI-generated output is being evaluated as code via eval() or the Function constructor. This is extremely dangerous — a prompt injection could lead to arbitrary code execution.
- **Evidence:**
  ```
        configurable: true,
      value: function () {
        if (
          this instanceof HTMLElement &&
  ```
- **Recommendation:** Never eval() AI outputs. Use structured output parsing (JSON.parse with validation) instead.

### SEC-006: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `frontend/src/lib/utils/markdown.test.ts:83`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
    let match: RegExpExecArray | null;
  while ((match = hrefPattern.exec(html)) !== null) {
    const value = match[1] ?? match[2] ?? match[3] ?? "";
    const norm = normalizeHref(value);
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-007: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `frontend/src/lib/utils/markdown.ts:22`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
      start(src) {
      const m = startRe.exec(src);
      return m?.index;
    },
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-008: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/automated_backfill_test.go:22`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// Force is_automated to 0 to simulate pre-migration state.
	_, err := d.getWriter().Exec(
		"UPDATE sessions SET is_automated = 0 WHERE id = 'missed'",
	)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-009: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/db.go:166`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  type Reader interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-010: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/db_test.go:415`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	requireNoError(t, err, "raw open")
	_, err = conn.Exec("PRAGMA user_version = 0")
	requireNoError(t, err, "reset version")
	conn.Close()
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-011: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/insights.go:72`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
	res, err := db.getWriter().Exec(`
		INSERT INTO insights (
			type, date_from, date_to, project,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-012: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/messages.go:237`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  		)
		if _, err := tx.Exec(query, args...); err != nil {
			first := batch[0].Ordinal
			last := batch[len(batch)-1].Ordinal
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-013: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/messages_test.go:43`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	requireNoError(t, d.Update(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"INSERT INTO excluded_sessions (id) VALUES (?)",
			"excluded",
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-014: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/pins.go:64`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// the subsequent SELECT to detect a missing pin.
	if _, err := db.getWriter().Exec(
		`INSERT INTO pinned_messages (session_id, message_id, ordinal, note)
		 SELECT ?, m.id, m.ordinal, ?
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-015: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/pricing.go:50`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	for _, p := range prices {
		if _, err := stmt.Exec(
			p.ModelPattern,
			p.InputPerMTok,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-016: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/search_content_test.go:222`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	} {
		if _, err := d.getWriter().Exec(sql); err != nil {
			t.Fatalf("force empty tool_use_id: %v", err)
		}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-017: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/search_test.go:437`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// JOIN to return duplicate rows before the MATCH clause was added.
	if _, err := d.getWriter().Exec(
		"INSERT INTO messages_fts(messages_fts) VALUES('optimize')",
	); err != nil {
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-018: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/secret_findings.go:57`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  ) error {
	if _, err := tx.Exec(
		"DELETE FROM secret_findings WHERE session_id = ?",
		sessionID,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-019: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/secret_findings_test.go:171`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// it rather than regenerating from the destination default.
	_, err = srcDB.getWriter().Exec(
		"UPDATE secret_findings SET created_at = ? WHERE session_id = 's1'",
		"2020-01-01T00:00:00.000Z")
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-020: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/session_batch.go:60`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  		savepoint := fmt.Sprintf("session_batch_%d", i)
		if _, err := tx.Exec("SAVEPOINT " + savepoint); err != nil {
			return result, fmt.Errorf(
				"creating savepoint %s: %w", savepoint, err,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-021: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/session_stats_test.go:189`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	}
	if _, err := d.getWriter().Exec(
		"UPDATE sessions SET is_automated = ? WHERE id = ?",
		want, f.id,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-022: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/sessions.go:875`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	defer db.mu.Unlock()
	_, err := db.getWriter().Exec(
		"DELETE FROM sessions WHERE id IN (SELECT id FROM excluded_sessions)",
	)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-023: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/signals.go:68`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  ) error {
	_, err := tx.Exec(`
		UPDATE sessions SET
			tool_failure_signal_count = ?,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-024: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/usage_events.go:32`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  func (db *DB) ensureUsageEventsSchemaLocked(w *sql.DB) error {
	if _, err := w.Exec(`
		CREATE TABLE IF NOT EXISTS usage_events (
			id INTEGER PRIMARY KEY,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-025: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/db/usage_test.go:708`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  		t.Helper()
		if _, err := d.getWriter().Exec(q, args...); err != nil {
			t.Fatalf("exec %q: %v", q, err)
		}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-026: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/antigravity_test.go:636`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	defer db.Close()
	mustExec(t, db, `CREATE TABLE trajectory_meta (
		trajectory_id text, cascade_id text,
		trajectory_type integer, source integer,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-027: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/forge_test.go:36`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	s.t.Helper()
	_, err := s.db.Exec(
		`INSERT INTO conversations
		 (conversation_id, title, workspace_id, context, created_at, updated_at, metrics)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-028: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/hermes_test.go:51`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	t.Cleanup(func() { _ = db.Close() })
	_, err = db.Exec(`
		CREATE TABLE sessions (
			id TEXT PRIMARY KEY,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-029: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/kiro_sqlite_test.go:31`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	t.Cleanup(func() { db.Close() })
	if _, err := db.Exec(kiroSQLiteSchema); err != nil {
		t.Fatalf("create kiro sqlite schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-030: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/opencode_test.go:67`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	s.t.Helper()
	_, err := s.db.Exec(`INSERT INTO project (id, worktree) VALUES (?, ?)`, id, worktree)
	if err != nil {
		s.t.Fatalf("add project: %v", err)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-031: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/piebald_test.go:88`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("exec schema: %v", err)
		}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-032: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/parser/warp_test.go:47`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	s.t.Helper()
	_, err := s.db.Exec(
		`INSERT INTO agent_conversations
		 (conversation_id, conversation_data, last_modified_at)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-033: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/curation.go:15`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  func (s *Store) StarSession(sessionID string) (bool, error) {
	res, err := s.pg.Exec(`
		INSERT INTO starred_sessions (session_id)
		SELECT $1 WHERE EXISTS (
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-034: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/messages_pgtest_test.go:80`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// Insert a newer session that also matches "hello".
	_, err = pg.Exec(`
		INSERT INTO sessions
			(id, machine, project, agent, first_message,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-035: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/push_pgtest_test.go:36`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	ctx := context.Background()
	if _, err := pg.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE`); err != nil {
		t.Fatalf("drop schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-036: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/search_content_pgtest_test.go:29`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
	if _, err := pg.Exec(`DROP SCHEMA IF EXISTS ` + contentSearchSchema + ` CASCADE`); err != nil {
		t.Fatalf("drop schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-037: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/secret_findings_pgtest_test.go:27`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
	if _, err := pg.Exec(
		`DELETE FROM secret_findings WHERE session_id = $1`, sid,
	); err != nil {
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-038: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/session_timing_pgtest_test.go:15`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	t.Helper()
	if _, err := pg.Exec(
		`DELETE FROM sessions WHERE id = $1`, sessionID,
	); err != nil {
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-039: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/sessions_pgtest_test.go:27`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	// Seed a session with leaks and one without.
	_, err = pg.Exec(`
		INSERT INTO sessions
			(id, machine, project, agent, first_message,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-040: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/store_test.go:23`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
	_, err = pg.Exec(`
		DROP SCHEMA IF EXISTS ` + testSchema + ` CASCADE;
	`)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-041: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/sync_test.go:33`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	defer pg.Close()
	_, _ = pg.Exec(
		"DROP SCHEMA IF EXISTS agentsview CASCADE",
	)
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-042: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/token_usage_pgtest_test.go:23`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
	_, err = pg.Exec(`
		INSERT INTO sessions (
			id, machine, project, agent,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-043: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/postgres/usage_pgtest_test.go:26`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	ctx := context.Background()
	if _, err := pg.Exec(`DROP SCHEMA IF EXISTS ` + schema + ` CASCADE`); err != nil {
		t.Fatalf("drop schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-044: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/scheduler/runner.go:289`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  			run.SessionID = sessionID
			if _, dbErr := r.store.db.Writer().Exec(
				`UPDATE scheduler_runs SET session_id = ? WHERE id = ?`,
				sessionID, run.ID,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-045: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/scheduler/store.go:197`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	}
	_, err := s.db.Writer().Exec(
		`INSERT INTO scheduler_runs (id, job_id, session_id, started_at, status)
		VALUES (?, ?, ?, ?, ?)`,
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-046: 🔴 Hardcoded API key or secret

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/secrets/bench_test.go:18`
- **Description:** A credential or API key appears to be hardcoded in the source code.
- **Evidence:**
  ```
  	hash := sha256.Sum256(data) // returns a 32-byte digest
	token = "ghp_xxxxNOTAREALTOKENxxxxxxxxxxxxxxxxxxxx"
}
Here is a normal paragraph of assistant prose explaining the change in
  ```
- **Recommendation:** Move all secrets to environment variables. Use a secrets manager for production. Never commit secrets to version control.

### SEC-047: 🔴 GitHub token in source code

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/secrets/bench_test.go:18`
- **Description:** A GitHub personal access token or OAuth token is hardcoded in the source code.
- **Evidence:**
  ```
  	hash := sha256.Sum256(data) // returns a 32-byte digest
	token = "ghp_xxxxNOTAREALTOKENxxxxxxxxxxxxxxxxxxxx"
}
Here is a normal paragraph of assistant prose explaining the change in
  ```
- **Recommendation:** Remove the token, rotate it on GitHub, and use environment variables.

### SEC-048: 🔴 GitHub token in source code

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/secrets/placeholder_test.go:308`
- **Description:** A GitHub personal access token or OAuth token is hardcoded in the source code.
- **Evidence:**
  ```
  		{"slack", "TOKEN=xoxb-549271836401-fHk7Bm3Pz9Wt5Vx2Yq8Nc done"},
		{"github_pat", "tok ghp_8Hk3Wn7Dz4Rp2Vx9Mb6Tj0Qc5Lm1Yp8Bv4Hg"},
		{"stripe", "key=sk_live_7Qh3Wn8Dk4Rp9Vx2Mb6Tj0Qc5Lm done"},
		{"google_api", "key=AIza7Qh3Wn8Dk4Rp9Vx2Mb6Tj0Qc5Lm1Yp8Bv4H end"},
  ```
- **Recommendation:** Remove the token, rotate it on GitHub, and use environment variables.

### SEC-049: 🔴 GitHub token in source code

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/secrets/rules_test.go:16`
- **Description:** A GitHub personal access token or OAuth token is hardcoded in the source code.
- **Evidence:**
  ```
  		{"github classic", "github-pat",
			"tok ghp_8Hk3Wn7Dz4Rp2Vx9Mb6Tj0Qc5Lm1Yp8Bv4Hg x", true},
		{"github fine-grained", "github-pat",
			"github_pat_8Hk3Wn7Dz4Rp2Vx9Mb6Tj0Qc5Lm1Yp8Bv4HgN_X2cWp9", true},
  ```
- **Recommendation:** Remove the token, rotate it on GitHub, and use environment variables.

### SEC-050: 🔴 GitHub token in source code

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/secrets/secrets_test.go:159`
- **Description:** A GitHub personal access token or OAuth token is hardcoded in the source code.
- **Evidence:**
  ```
  		"AKIA7QHWN2DKR4FYPLJM",
		"ghp_8Hk3Wn7Dz4Rp2Vx9Mb6Tj0Qc5Lm1Yp8Bv4Hg",
		"xoxb-549271836401-fHk7Bm3Pz9Wt5Vx2Yq8Nc",
		"xoxs-302846159270-xPk9Bm3Wv8Qt5Lz2Yh7Fc",
  ```
- **Recommendation:** Remove the token, rotate it on GitHub, and use environment variables.

### SEC-051: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/server/server_test.go:1344`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	err := te.db.Update(func(tx *sql.Tx) error {
		_, err := tx.Exec("DROP TABLE IF EXISTS messages_fts")
		return err
	})
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-052: 🔴 Hardcoded API key or secret

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/service/http_test.go:386`
- **Description:** A credential or API key appears to be hardcoded in the source code.
- **Evidence:**
  ```
  	t.Parallel()
	const goodToken = "correct-horse-battery-staple"
	baseURL, d := newHTTPTestServerWithCfg(t, config.Config{
		RequireAuth: true,
  ```
- **Recommendation:** Move all secrets to environment variables. Use a secrets manager for production. Never commit secrets to version control.

### SEC-053: 🔴 Hardcoded API key or secret

- **Severity:** CRITICAL
- **Category:** credential exposure
- **File:** `internal/service/secrets_test.go:153`
- **Description:** A credential or API key appears to be hardcoded in the source code.
- **Evidence:**
  ```
  	d := dbtest.OpenTestDB(t)
	const secret = "AKIA7QHWN2DKR4FYPLJM"
	content := "my key is " + secret + " ok"
	start := strings.Index(content, secret)
  ```
- **Recommendation:** Move all secrets to environment variables. Use a secrets manager for production. Never commit secrets to version control.

### SEC-054: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/sync/engine_integration_test.go:1414`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  
func TestSyncAllImportsCodexExec(
	t *testing.T,
) {
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-055: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/sync/forge_integration_test.go:39`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	`
	if _, err := d.Exec(schema); err != nil {
		t.Fatalf("creating forge schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-056: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/sync/piebald_integration_test.go:93`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	`
	if _, err := d.Exec(schema); err != nil {
		t.Fatalf("creating piebald schema: %v", err)
	}
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-057: 🔴 AI output passed to shell execution

- **Severity:** CRITICAL
- **Category:** unsafe eval
- **File:** `internal/sync/test_helpers_test.go:87`
- **Description:** AI-generated output is being passed to a shell command execution function. Prompt injection could lead to arbitrary command execution.
- **Evidence:**
  ```
  	err := e.db.Update(func(tx *sql.Tx) error {
		_, err := tx.Exec(
			"UPDATE sessions SET file_mtime = NULL"+
				" WHERE id = ?",
  ```
- **Recommendation:** Never execute AI outputs as shell commands. If code execution is needed, use a sandboxed environment.

### SEC-058: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/agentsview/health_test.go:225`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:           42,
		UserMessageCount:       12,
		HealthGrade:            &a,
		HealthScore:            &score,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-059: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/agentsview/session_list.go:20`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		minMessages, maxMessages                int
		minUserMessages                         int
		includeOneShot                          bool
		includeAutomated, includeChildren       bool
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-060: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/agentsview/session_test.go:65`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	t.Cleanup(func() { d.Close() })
	// UserMessageCount >= 2 so seeded sessions pass the default
	// ExcludeOneShot filter in `session list` (one-shot means
	// user_message_count <= 1). See internal/db/analytics.go.
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-061: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/agentsview/stats.go:254`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		fmtInt(s.Totals.MessagesTotal),
		fmtInt(s.Totals.UserMessagesTotal))
	fmt.Fprintln(w)
}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-062: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/agentsview/stats_test.go:42`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			MessagesTotal:      109324,
			UserMessagesTotal:  3012,
		},
		Archetypes: db.StatsArchetypes{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-063: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `cmd/testfixture/main.go:135`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     spec.msgCount,
		UserMessageCount: spec.userMsgCount,
		RelationshipType: spec.relationshipType,
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-064: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `desktop/src-tauri/src/lib.rs:1420`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
          format!("GET /api/v1/version HTTP/1.1\r\nHost: {HOST}:{port}\r\nConnection: close\r\n\r\n");
    let response = match read_http_response(port, request.as_str()) {
        Some(resp) => resp,
        None => return false,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-065: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/e2e/message-loading.spec.ts:248`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
                thinking_text: "",
              has_tool_use: false,
              content_length: content.length,
              model: "",
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-066: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/public/d3.v7.min.js:2`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
  // https://d3js.org v7.9.0 Copyright 2010-2023 Mike Bostock
!function(t,n){"object"==typeof exports&&"undefined"!=typeof module?n(exports):"function"==typeof define&&define.amd?define(["exports"],n):n((t="undefined"!=typeof globalThis?globalThis:t||self).d3=t.d3||{})}(this,(function(t){"use strict";...
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-067: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/analytics/TopSessions.test.ts:53`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
      // before their unmount call.
    document.body.innerHTML = "";
  });

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-068: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/content/CompactBoundaryDivider.test.ts:8`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
  afterEach(() => {
  document.body.innerHTML = "";
});

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-069: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/content/MessageContent.test.ts:82`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
  afterEach(() => {
  document.body.innerHTML = "";
  vi.clearAllMocks();
});
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-070: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/components/content/MessageList.test.ts:51`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
      thinking_text: "",
    has_tool_use: false,
    content_length: 6,
    model: "",
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-071: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/content/SubagentInline.test.ts:66`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
    getSession.mockReset();
  document.body.innerHTML = "";
});

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-072: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/content/ToolBlock.test.ts:20`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
      if (component) unmount(component);
    document.body.innerHTML = "";
  });

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-073: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/layout/AppHeader.test.ts:49`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
      }
    document.body.innerHTML = "";
  });

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-074: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/layout/SessionBreadcrumb.test.ts:56`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
  afterEach(() => {
  document.body.innerHTML = "";
});

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-075: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/layout/StatusBar.test.ts:29`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
    afterEach(() => {
    document.body.innerHTML = "";
    vi.useRealTimers();
    sync.lastSync = null;
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-076: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/layout/ThreeColumnLayout.test.ts:206`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
    document.body.className = "";
  document.body.innerHTML = "";
  restoreMeasuredLayoutWidth?.();
  ui.sidebarOpen = true;
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-077: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/sidebar/SessionList.test.ts:90`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
      }
    document.body.innerHTML = "";
    Object.defineProperty(globalThis, "ResizeObserver", {
      configurable: true,
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-078: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/system/SystemBoundaryCard.test.ts:8`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
  afterEach(() => {
  document.body.innerHTML = "";
});

  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-079: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/components/trends/TrendsPage.test.ts:74`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
      }
    document.body.innerHTML = "";
    window.history.replaceState(null, "", "/");
    vi.unstubAllGlobals();
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-080: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `frontend/src/lib/stores/analytics.svelte.ts:74`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
    termination: string = $state("");
  minUserMessages: number = $state(0);
  includeOneShot: boolean = $state(true);
  includeAutomated: boolean = $state(false);
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-081: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/stores/analytics.svelte.ts:88`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
    velocity = $state<VelocityResponse | null>(null);
  tools = $state<ToolsAnalyticsResponse | null>(null);
  topSessions = $state<TopSessionsResponse | null>(null);
  signals = $state<SignalsAnalyticsResponse | null>(null);
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-082: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/stores/analytics.test.ts:31`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
    getAnalyticsVelocity: vi.fn(),
  getAnalyticsTools: vi.fn(),
  getAnalyticsTopSessions: vi.fn(),
  getAnalyticsSignals: vi.fn(),
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-083: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/stores/messages.test.ts:58`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
      thinking_text: "",
    has_tool_use: false,
    content_length: 6,
    model: "",
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-084: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `frontend/src/lib/stores/sessions.svelte.ts:63`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
    maxMessages: number;
  minUserMessages: number;
  includeOneShot: boolean;
  includeAutomated: boolean;
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-085: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `frontend/src/lib/stores/sessions.test.ts:254`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
        sessions.filters.maxMessages = 20;
      sessions.filters.minUserMessages = 1;
      sessions.filters.includeOneShot = false;
      sessions.filters.includeAutomated = true;
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-086: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `frontend/src/lib/stores/usage.svelte.ts:211`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
        min_user_messages:
        sessionFilters.minUserMessages > 0
          ? sessionFilters.minUserMessages
          : undefined,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-087: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `frontend/src/lib/stores/usage.test.ts:181`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
      sessions.filters.agent = "claude,codex";
    sessions.filters.minUserMessages = 5;
    sessions.filters.includeOneShot = false;
    sessions.filters.includeAutomated = true;
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-088: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/utils/csv-export.ts:15`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
    projects: ProjectsAnalyticsResponse | null;
  tools: ToolsAnalyticsResponse | null;
  velocity: VelocityResponse | null;
}
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-089: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/utils/highlight.test.ts:7`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
    const div = document.createElement("div");
  div.innerHTML = html;
  return div;
}
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-090: 🟠 AI output rendered as raw HTML

- **Severity:** HIGH
- **Category:** missing output filtering
- **File:** `frontend/src/lib/utils/markdown.test.ts:11`
- **Description:** AI-generated output is being rendered as raw HTML without sanitization. An injection attack could produce malicious HTML/JS.
- **Evidence:**
  ```
    const div = document.createElement("div");
  div.innerHTML = html;
  return div;
}
  ```
- **Recommendation:** Always sanitize AI outputs before rendering as HTML. Use DOMPurify or a similar library. Prefer text rendering over HTML.

### SEC-091: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `frontend/src/lib/utils/messages.test.ts:19`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
      thinking_text: "",
    has_tool_use: false,
    content_length: 5,
    model: "",
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-092: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/analytics.go:53`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	Hour             *int   // nil = all, 0-23
	MinUserMessages  int    // user_message_count >= N
	ExcludeOneShot   bool   // exclude sessions with user_message_count <= 1
	ExcludeAutomated bool   // exclude automated (roborev) sessions
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-093: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `internal/db/analytics.go:1436`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
  		SUM(CASE WHEN role='assistant'
			AND has_tool_use=1 THEN 1 ELSE 0 END)
		FROM messages
		WHERE session_id IN ` + ph + `
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-094: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/analytics_test.go:14`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	TotalMessages          int
	TotalUserMessages      int
	TotalAssistantMessages int
	ActiveProjects         int
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-095: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/automated_backfill_test.go:19`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.MessageCount = 3
		s.UserMessageCount = 1
	})
	// Force is_automated to 0 to simulate pre-migration state.
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-096: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/automated_test.go:16`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		{"EmptyMessage", "", false},
		{"NormalUserPrompt", "fix the login bug", false},

		// Code review
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-097: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/db_test.go:5564`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:         5,
		UserMessageCount:     2,
		TotalOutputTokens:    500,
		PeakContextTokens:    1500,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-098: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/filter_test.go:170`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
func TestSessionFilterMinUserMessages(t *testing.T) {
	d := testDB(t)

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-099: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/messages_test.go:68`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  				MessageCount:     2,
				UserMessageCount: 1,
			},
			Messages: []Message{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-100: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/search_content_test.go:15`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	t.Helper()
	// UserMessageCount > 1 so the session is not treated as one-shot and
	// excluded by the default session-list-parity filter.
	insertSession(t, d, id, project, func(s *Session) {
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-101: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/session_stats.go:355`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	messageCount         int
	userMessageCount     int
	totalOutputTokens    int64
	hasTotalOutputTokens bool
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-102: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/session_stats_buckets.go:16`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  // Represented as two separate edge lists for clarity.
var userMessagesEdgesAll = []float64{0, 2, 6, 16, 31, 51, math.Inf(1)}
var userMessagesEdgesHuman = []float64{2, 6, 16, 31, 51, math.Inf(1)}

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-103: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/session_stats_buckets_test.go:25`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
func TestAssignBucketUserMessagesAll(t *testing.T) {
	cases := []struct {
		v    float64
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-104: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/session_stats_test.go:156`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.Agent = agent
		s.UserMessageCount = f.userMsgs
		s.MessageCount = mc
		if f.startedAt != "" {
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-105: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/session_stats_types.go:52`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	MessagesTotal      int `json:"messages_total"`
	UserMessagesTotal  int `json:"user_messages_total"`
}

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-106: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/sessions.go:121`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		&s.FirstMessage, &s.DisplayName, &s.StartedAt, &s.EndedAt,
		&s.MessageCount, &s.UserMessageCount,
		&s.ParentSessionID, &s.RelationshipType,
		&s.TotalOutputTokens, &s.PeakContextTokens,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-107: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/sessions_test.go:68`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			s.MessageCount = 5
			s.UserMessageCount = 3
		})
		err := d.UpdateSessionSignals(tc.id, SessionSignalUpdate{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-108: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `internal/db/timing.go:111`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
  		  CASE
		    WHEN m2.has_tool_use = 0 THEN NULL
		    WHEN m2.delta_ms < 0    THEN NULL
		    ELSE m2.delta_ms
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-109: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/token_usage_test.go:482`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:         5,
		UserMessageCount:     2,
		TotalOutputTokens:    1000,
		PeakContextTokens:    8000,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-110: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/trends_test.go:164`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.MessageCount = 3
		s.UserMessageCount = 2
	})
	insertMessages(t, d,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-111: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/usage.go:27`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	Timezone         string // IANA timezone, "" for UTC
	MinUserMessages  int    // user_message_count >= N
	ExcludeOneShot   bool   // user_message_count > 1
	ExcludeAutomated bool   // is_automated = false
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-112: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/db/usage_test.go:44`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.StartedAt = new("2026-05-14T10:00:00Z")
		s.UserMessageCount = 2
	})

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-113: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/importer/importer.go:194`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     s.MessageCount,
		UserMessageCount: s.UserMessageCount,
	}

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-114: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/insight/prompt_test.go:177`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  					s.MessageCount = 5
					s.UserMessageCount = 2
					s.StartedAt = new("2025-01-15T10:00:00Z")
					s.EndedAt = new("2025-01-15T11:00:00Z")
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-115: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/amp.go:156`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-116: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/amp_test.go:58`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertMessageCount(t, sess.MessageCount, 2)
	assert.Equal(t, 1, sess.UserMessageCount)

	// Start time from created (epoch ms).
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-117: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/antigravity.go:161`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-118: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/antigravity_cli.go:199`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-119: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/antigravity_test.go:255`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if sess.MessageCount != 2 || sess.UserMessageCount != 1 {
		t.Fatalf(
			"counts: msg=%d user=%d",
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-120: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/chatgpt.go:199`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(msgs),
		UserMessageCount: userCount,
	}

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-121: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/chatgpt_test.go:110`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assert.Equal(t, 3, s.MessageCount)
	assert.Equal(t, 2, s.UserMessageCount)

	msgs := results[0].Messages
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-122: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/claude.go:671`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File:             fileInfo,
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-123: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/claude_ai.go:131`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		userCount        int
		firstUserMessage string
	)

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-124: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/claude_ai_test.go:96`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assert.Equal(t, 2, s.MessageCount)
	assert.Equal(t, 1, s.UserMessageCount)
	assert.Equal(t,
		"2026-01-15T10:00:00.000000Z",
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-125: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/claude_parser_test.go:139`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		assert.Equal(t, 2, sess.MessageCount)
		// Only real user messages count toward UserMessageCount.
		assert.Equal(t, 1, sess.UserMessageCount)
		// FirstMessage is from the first real user message.
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-126: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/codex.go:1194`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:      len(b.messages),
		UserMessageCount:  userCount,
		TerminationStatus: classifyCodexTermination(b.lastTaskEvent),
		File: FileInfo{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-127: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/codex_parser_test.go:1123`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		assert.Equal(t, "OK", msgs[1].Content)
		assert.Equal(t, 1, sess.UserMessageCount,
			"skill injection must not count as a user turn")
	})
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-128: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/copilot.go:17`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	copilotEventSessionStart    = "session.start"
	copilotEventUserMessage     = "user.message"
	copilotEventAssistantMsg    = "assistant.message"
	copilotEventToolComplete    = "tool.execution_complete"
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-129: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/copilot_test.go:439`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
func TestCopilotUserMessageCount(t *testing.T) {
	// Tool-result user messages (Content == "") should not count
	// as user prompts. This was the exact bug: Copilot emits
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-130: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/cortex.go:55`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	InternalOnly *bool             `json:"internalOnly"`
	IsUserPrompt *bool             `json:"is_user_prompt"`
	MessageID    string            `json:"message_id"`
	ToolUse      *cortexToolUse    `json:"tool_use"`
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-131: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/cortex_test.go:53`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertMessageCount(t, sess.MessageCount, 2)
	assert.Equal(t, 1, sess.UserMessageCount)

	require.Len(t, msgs, 2)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-132: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/cursor.go:188`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	if role == RoleUser {
		content := extractUserQuery(lines)
		return content, false, nil
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-133: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/forge.go:353`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File:             fileInfo,
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-134: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/forge_test.go:183`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertEq(t, "DisplayName", s.Session.DisplayName, "Add Forge Support")
	assertEq(t, "UserMessageCount", s.Session.UserMessageCount, 1)
	assertEq(t, "FirstMessage", s.Session.FirstMessage, "Please add Forge support.")
	assertEq(t, "Cwd", s.Session.Cwd, "/home/mj/dev/projects/agentsview")
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-135: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/gemini.go:303`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-136: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/hermes.go:322`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-137: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/hermes_test.go:309`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assert.Equal(t, "compact_boundary", msgs[0].SourceSubtype)
	assert.Equal(t, 1, sess.UserMessageCount)
	assert.Equal(t, "real prompt", sess.FirstMessage)
}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-138: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/iflow.go:193`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File:             fileInfo,
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-139: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kimi.go:226`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
			userText := strings.Join(userParts, "\n")
			if userText == "" {
				continue
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-140: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kimi_test.go:53`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertMessageCount(t, sess.MessageCount, 2)
	assert.Equal(t, 1, sess.UserMessageCount)

	wantStart := time.Unix(1704067200, 0)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-141: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kiro.go:260`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-142: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kiro_ide.go:392`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-143: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kiro_sqlite.go:332`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  KiroSQLiteVirtualPath(s.dbPath, row.conversationID),
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-144: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/kiro_sqlite_test.go:87`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertEq(t, "MessageCount", sess.MessageCount, 4)
	assertEq(t, "UserMessageCount", sess.UserMessageCount, 2)
	assertEq(t, "File.Path", sess.File.Path, dbPath+"#sqlite-session")
	assertEq(t, "File.Mtime", sess.File.Mtime, int64(1779012030000)*1_000_000)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-145: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/openclaw.go:223`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-146: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/openclaw_test.go:76`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if sess.UserMessageCount != 1 {
		t.Errorf("expected 1 user message, got %d", sess.UserMessageCount)
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-147: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/opencode.go:608`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(parsed),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  filePath,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-148: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/openhands.go:335`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File:             snapshot,
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-149: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/openhands_test.go:132`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assert.Equal(t, 4, sess.MessageCount)
	assert.Equal(t, 1, sess.UserMessageCount)
	assert.Equal(t, sessionDir, sess.File.Path)
	assert.NotEmpty(t, sess.File.Hash)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-150: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/parser_test.go:1066`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
func TestCodexUserMessageCount(t *testing.T) {
	content := testjsonl.JoinJSONL(
		testjsonl.CodexSessionMetaJSON(
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-151: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/pi.go:116`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			case "user":
				msg := parsePiUserMessage(line, ordinal)
				if msg == nil {
					continue
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-152: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/pi_test.go:55`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
// TestParsePiSession_UserMessages verifies user message content and ordinals
// (PRSR-02, PRSR-01).
func TestParsePiSession_UserMessages(t *testing.T) {
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-153: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/piebald.go:310`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		sess.MessageCount = len(messages)
		sess.UserMessageCount = realUserCount
		accumulateMessageTokenUsage(&sess, messages)

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-154: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/qclaw.go:217`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-155: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/qclaw_test.go:76`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if sess.UserMessageCount != 1 {
		t.Errorf("expected 1 user message, got %d", sess.UserMessageCount)
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-156: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/qwen.go:256`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-157: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/qwen_test.go:33`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assert.Equal(t, 2, sess.MessageCount)
	assert.Equal(t, 1, sess.UserMessageCount)
	assert.True(t, sess.HasTotalOutputTokens)
	assert.Equal(t, 47, sess.TotalOutputTokens)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-158: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/types.go:535`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	MessageCount     int
	UserMessageCount int
	File             FileInfo

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-159: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/vscode_copilot.go:269`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: userCount,
	}

  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-160: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/warp.go:355`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(parsed),
		UserMessageCount: userCount,
		File: FileInfo{
			Path:  dbPath + "#" + c.id,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-161: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/warp_test.go:163`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertEq(t, "Project", s.Session.Project, "myproject")
	assertEq(t, "UserMessageCount", s.Session.UserMessageCount, 2)
	assertEq(t, "FirstMessage",
		s.Session.FirstMessage,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-162: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/workbuddy.go:265`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     len(messages),
		UserMessageCount: realUserCount,
		File: FileInfo{
			Path:  path,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-163: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/workbuddy_test.go:65`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if sess.Cwd != "/tmp/proj" || sess.FirstMessage != "hello" || sess.UserMessageCount != 1 {
		t.Fatalf("session = %+v", sess)
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-164: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/zencoder.go:90`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	case "user":
		b.handleUserMessage(line, ts)
	case "assistant":
		b.handleAssistantMessage(line, ts)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-165: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/parser/zencoder_test.go:45`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertMessageCount(t, sess.MessageCount, 4)
	assert.Equal(t, 1, sess.UserMessageCount)

	wantStart := mustParseTime(t, "2024-01-01T00:00:00Z")
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-166: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/analytics.go:156`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if f.MinUserMessages > 0 {
		preds = append(preds,
			"user_message_count >= "+
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-167: 🟠 AI tool/function calling with network access

- **Severity:** HIGH
- **Category:** data exfiltration
- **File:** `internal/postgres/analytics.go:1358`
- **Description:** AI function calling or tool use configuration allows network requests. A prompt injection could instruct the AI to exfiltrate data via tool calls.
- **Evidence:**
  ```
  		SUM(CASE WHEN role='assistant'
			AND has_tool_use=true THEN 1 ELSE 0 END)
		FROM messages
		WHERE session_id IN ` + ph + `
  ```
- **Recommendation:** Restrict AI tool capabilities to a minimum. Validate all tool call parameters. Block network-accessing tools or whitelist allowed URLs.

### SEC-168: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/push.go:653`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		fmt.Sprintf("%d", sess.MessageCount),
		fmt.Sprintf("%d", sess.UserMessageCount),
		fmt.Sprintf("%t", sess.IsAutomated),
		fmt.Sprintf("%d", sess.TotalOutputTokens),
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-169: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/push_pgtest_test.go:176`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:     1,
		UserMessageCount: 1,
		// CreatedAt must be parseable by ParseSQLiteTimestamp;
		// PG's NOT NULL on created_at would otherwise reject NULL.
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-170: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/sessions.go:141`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		&createdAt, &startedAt, &endedAt,
		&s.MessageCount, &s.UserMessageCount,
		&s.ParentSessionID, &s.RelationshipType,
		&s.TotalOutputTokens, &s.PeakContextTokens,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-171: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/store_test.go:259`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		messageCount:     3,
		userMessageCount: 2,
	}
	for _, opt := range opts {
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-172: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/token_usage_pgtest_test.go:128`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:         1,
		UserMessageCount:     0,
		TotalOutputTokens:    500,
		PeakContextTokens:    900,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-173: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/postgres/usage.go:74`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
	if f.MinUserMessages > 0 {
		query += " AND u.user_message_count >= " + pb.add(f.MinUserMessages)
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-174: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/analytics.go:107`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		Hour:             hour,
		MinUserMessages:  minUserMsgs,
		ExcludeOneShot:   !includeOneShot,
		ExcludeAutomated: !includeAutomated,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-175: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/analytics_test.go:507`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  					for _, e := range resp.Series {
						totalUser += e.UserMessages
						totalAsst += e.AssistantMessages
					}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-176: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/export.go:107`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	rawURL := fmt.Sprintf(
		"https://gist.githubusercontent.com/%s/%s/raw/%s",
		gist.Owner.Login, gist.ID, encoded,
	)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-177: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/server_test.go:386`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.MessageCount = msgCount
		s.UserMessageCount = max(msgCount, 2)
		s.StartedAt = new(tsSeed)
		s.EndedAt = new(tsSeedEnd)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-178: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/session_timing_test.go:59`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			s.MessageCount = 3
			s.UserMessageCount = 2
			s.StartedAt = new(string(startedAt))
			s.EndedAt = new(string(endedAt))
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-179: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/sessions.go:62`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MaxMessages:      maxMsgs,
		MinUserMessages:  minUserMsgs,
		IncludeOneShot:   q.Get("include_one_shot") == "true",
		IncludeAutomated: q.Get("include_automated") == "true",
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-180: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/upload.go:239`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MessageCount:         sess.MessageCount,
		UserMessageCount:     sess.UserMessageCount,
		ParentSessionID:      strPtr(sess.ParentSessionID),
		RelationshipType:     string(sess.RelationshipType),
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-181: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/usage.go:133`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		Timezone:         tz,
		MinUserMessages:  minUserMsgs,
		ExcludeOneShot:   !includeOneShot,
		ExcludeAutomated: !includeAutomated,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-182: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/server/usage_test.go:72`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			sess.EndedAt = &ts
			sess.UserMessageCount = 1
		},
	)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-183: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/direct.go:145`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		MaxMessages:          f.MaxMessages,
		MinUserMessages:      f.MinUserMessages,
		ExcludeOneShot:       !f.IncludeOneShot,
		ExcludeAutomated:     !f.IncludeAutomated,
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-184: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/direct_test.go:71`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  			s.MessageCount = 2
			s.UserMessageCount = 2
		})
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-185: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/http.go:103`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	}
	if f.MinUserMessages > 0 {
		q.Set("min_user_messages", strconv.Itoa(f.MinUserMessages))
	}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-186: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/http_test.go:344`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	baseURL, d := newHTTPTestServer(t)
	// Seed a session with UserMessageCount=2 so content search includes it.
	dbtest.SeedSession(t, d, "cs-1", "search-proj", func(s *db.Session) {
		s.MessageCount = 3
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-187: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/search_content_test.go:17`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  // seedServiceSearchSession creates a session with a single user message
// whose content contains the given text. The session has UserMessageCount=2
// so it is not excluded by the default one-shot filter.
func seedServiceSearchSession(
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-188: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/secrets_test.go:61`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		s.MessageCount = 1
		s.UserMessageCount = 1
	})
	require.NoError(t, d.InsertMessages([]db.Message{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-189: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/service/service.go:125`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	MaxMessages      int    `json:"max_messages,omitempty"`
	MinUserMessages  int    `json:"min_user_messages,omitempty"`
	IncludeOneShot   bool   `json:"include_one_shot,omitempty"`
	IncludeAutomated bool   `json:"include_automated,omitempty"`
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-190: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/sync/engine.go:4339`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  func isAutomatedFromSession(s db.Session) bool {
	return s.UserMessageCount <= 1 &&
		s.FirstMessage != nil &&
		db.IsAutomatedSession(*s.FirstMessage)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-191: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/sync/engine_integration_test.go:4613`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	assertSessionState(t, env.db, "filter-count", func(sess *db.Session) {
		if sess.UserMessageCount != 1 {
			t.Errorf("user_message_count = %d, want 1", sess.UserMessageCount)
		}
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-192: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/sync/forge_integration_test.go:65`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  
func forgeTestContext(userPrompt, finalAnswer string) string {
	messages := []map[string]any{
		{
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-193: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/sync/secret_scan_test.go:149`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  		ID: id, Project: "proj", Machine: "m", Agent: "claude",
		MessageCount: 1, UserMessageCount: 1,
	}); err != nil {
		t.Fatalf("UpsertSession: %v", err)
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-194: 🟠 Unsanitized user input in AI prompt

- **Severity:** HIGH
- **Category:** prompt injection
- **File:** `internal/sync/test_helpers_test.go:501`
- **Description:** User-controlled input appears to flow directly into an AI prompt without sanitization or validation.
- **Evidence:**
  ```
  	sessionID string,
	userContent, assistantContent string,
	timeCreated int64,
) {
  ```
- **Recommendation:** Sanitize and validate all user inputs before including them in AI prompts. Use a dedicated input filter to strip injection patterns.

### SEC-195: 🟡 System prompt defined in client-accessible code

- **Severity:** MEDIUM
- **Category:** system prompt leakage
- **File:** `frontend/src/lib/utils/messages.test.ts:61`
- **Description:** A system prompt is defined in code that may be accessible to clients (frontend, API response, etc.).
- **Evidence:**
  ```
      ["stop hook", "Stop hook feedback: blocked"],
  ])("detects prefix-based system message: %s", (_label, content) => {
    expect(isSystemMessage(msg({ content }))).toBe(true);
  });
  ```
- **Recommendation:** Store system prompts server-side only. Never expose them in client bundles, API responses, or source maps.

### SEC-196: 🟡 System prompt defined in client-accessible code

- **Severity:** MEDIUM
- **Category:** system prompt leakage
- **File:** `internal/postgres/messages_pgtest_test.go:368`
- **Description:** A system prompt is defined in code that may be accessible to clients (frontend, API response, etc.).
- **Evidence:**
  ```
  	if err != nil {
		t.Fatalf("inserting system message: %v", err)
	}

  ```
- **Recommendation:** Store system prompts server-side only. Never expose them in client bundles, API responses, or source maps.
