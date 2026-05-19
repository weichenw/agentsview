package db

import "database/sql"

// Writer returns the single writer connection for direct SQL access.
// Used by the scheduler to write run history entries.
func (db *DB) Writer() *sql.DB { return db.getWriter() }
