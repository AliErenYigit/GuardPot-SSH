package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SSHAuditSQLiteRepo struct {
	db *sql.DB
}

func NewSSHAuditSQLiteRepo(db *sql.DB) *SSHAuditSQLiteRepo {
	return &SSHAuditSQLiteRepo{db: db}
}

func (r *SSHAuditSQLiteRepo) Log(ctx context.Context, userID, connectionID int64, remoteIP, event, detail string) error {
	const q = `
	INSERT INTO ssh_audit_logs (user_id, connection_id, remote_ip, event, detail, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, q, userID, connectionID, remoteIP, event, detail, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("audit log failed: %w", err)
	}
	return nil
}
