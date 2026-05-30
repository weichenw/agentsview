package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Writer returns the single writer connection for direct SQL access.
// Used by the scheduler to write run history entries.
func (db *DB) Writer() *sql.DB { return db.getWriter() }

// CopySchedulerRunsFrom copies scheduler run history from the source
// database into this one. It is safe to call when the source lacks
// the table (older databases) or when rows already exist (INSERT OR
// IGNORE).
func (db *DB) CopySchedulerRunsFrom(sourcePath string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	ctx := context.Background()
	conn, err := db.getWriter().Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx, "ATTACH DATABASE ? AS old_db", sourcePath,
	); err != nil {
		return fmt.Errorf("attaching source db: %w", err)
	}
	defer func() {
		_, _ = conn.ExecContext(
			ctx, "DETACH DATABASE old_db",
		)
	}()

	var tableExists int
	err = conn.QueryRowContext(ctx,
		"SELECT 1 FROM old_db.sqlite_master WHERE type='table' AND name='scheduler_runs'",
	).Scan(&tableExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("probing scheduler_runs table: %w", err)
	}
	if tableExists != 1 {
		return nil
	}

	_, err = conn.ExecContext(ctx, `
		INSERT OR IGNORE INTO scheduler_runs
			(id, job_id, session_id, started_at, finished_at, status, exit_code, error)
		SELECT
			id, job_id, session_id, started_at, finished_at, status, exit_code, error
		FROM old_db.scheduler_runs`)
	if err != nil {
		return fmt.Errorf("copying scheduler runs: %w", err)
	}
	return nil
}
