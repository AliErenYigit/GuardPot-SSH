package sqlite

import (
	"database/sql"
	"fmt"
)

func InitSchema(db *sql.DB) error {
	const q = `
	CREATE TABLE IF NOT EXISTS ssh_connections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,

		name TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		username TEXT NOT NULL,

		auth_type TEXT NOT NULL,              
		secret_enc TEXT NOT NULL,           
		created_at TEXT NOT NULL,

		FOREIGN KEY(user_id) REFERENCES users(id)
	);

	CREATE INDEX IF NOT EXISTS idx_ssh_user_id ON ssh_connections(user_id);

		CREATE TABLE IF NOT EXISTS ssh_audit_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		connection_id INTEGER NOT NULL,
		remote_ip TEXT NOT NULL,
		event TEXT NOT NULL,       -- "connect_ok" | "connect_fail" | "disconnect"
		detail TEXT NOT NULL,      -- error text vs.
		created_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_ssh_audit_user ON ssh_audit_logs(user_id);
	CREATE INDEX IF NOT EXISTS idx_ssh_audit_conn ON ssh_audit_logs(connection_id);


	`
	if _, err := db.Exec(q); err != nil {
		return fmt.Errorf("init schema failed: %w", err)
	}
	return nil
}
