package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"
)

type SSHSQLiteRepo struct {
	db *sql.DB
}

func NewSSHSQLiteRepo(db *sql.DB) *SSHSQLiteRepo {
	return &SSHSQLiteRepo{db: db}
}

func (r *SSHSQLiteRepo) Create(ctx context.Context, c domain.SSHConnection) (domain.SSHConnection, error) {
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}

	const q = `
	INSERT INTO ssh_connections
	(user_id, name, host, port, username, auth_type, secret_enc, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, q,
		c.UserID, c.Name, c.Host, c.Port, c.Username, string(c.AuthType), c.SecretEnc, c.CreatedAt.Format(time.RFC3339Nano),
	)
	if err != nil {
		return domain.SSHConnection{}, fmt.Errorf("create ssh connection failed: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.SSHConnection{}, fmt.Errorf("create ssh connection last insert id failed: %w", err)
	}
	c.ID = id
	return c, nil
}

func (r *SSHSQLiteRepo) ListByUser(ctx context.Context, userID int64) ([]domain.SSHConnection, error) {
	const q = `
	SELECT id, user_id, name, host, port, username, auth_type, secret_enc, created_at
	FROM ssh_connections
	WHERE user_id = ?
	ORDER BY id DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("list ssh connections failed: %w", err)
	}
	defer rows.Close()

	var out []domain.SSHConnection
	for rows.Next() {
		var c domain.SSHConnection
		var authType string
		var createdAt string
		if err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Host, &c.Port, &c.Username, &authType, &c.SecretEnc, &createdAt); err != nil {
			return nil, fmt.Errorf("list scan failed: %w", err)
		}
		c.AuthType = domain.SSHAuthType(authType)
		t, err := time.Parse(time.RFC3339Nano, createdAt)
		if err != nil {
			return nil, fmt.Errorf("list parse created_at failed: %w", err)
		}
		c.CreatedAt = t
		out = append(out, c)
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("list rows err: %w", err)
		}
	}
	return out, nil
}

func (r *SSHSQLiteRepo) GetByID(ctx context.Context, userID, id int64) (domain.SSHConnection, error) {
	const q = `
	SELECT id, user_id, name, host, port, username, auth_type, secret_enc, created_at
	FROM ssh_connections
	WHERE user_id = ? AND id = ?
	LIMIT 1
	`
	var c domain.SSHConnection
	var authType string
	var createdAt string
	err := r.db.QueryRowContext(ctx, q, userID, id).
		Scan(&c.ID, &c.UserID, &c.Name, &c.Host, &c.Port, &c.Username, &authType, &c.SecretEnc, &createdAt)
	if err == sql.ErrNoRows {
		return domain.SSHConnection{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.SSHConnection{}, fmt.Errorf("get ssh connection failed: %w", err)
	}
	c.AuthType = domain.SSHAuthType(authType)
	t, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		return domain.SSHConnection{}, fmt.Errorf("get parse created_at failed: %w", err)
	}
	c.CreatedAt = t
	return c, nil
}

func (r *SSHSQLiteRepo) Delete(ctx context.Context, userID, id int64) error {
	const q = `DELETE FROM ssh_connections WHERE user_id = ? AND id = ?`
	res, err := r.db.ExecContext(ctx, q, userID, id)
	if err != nil {
		return fmt.Errorf("delete ssh connection failed: %w", err)
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		return repository.ErrNotFound
	}
	return nil
}
