package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"
)

type UserSQLiteRepo struct {
	db *sql.DB
}

func NewUserSQLiteRepo(db *sql.DB) *UserSQLiteRepo {
	return &UserSQLiteRepo{db: db}
}

func (r *UserSQLiteRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const q = `SELECT 1 FROM users WHERE email = ? LIMIT 1`
	var one int
	err := r.db.QueryRowContext(ctx, q, email).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("exists by email failed: %w", err)
	}
	return true, nil
}

func (r *UserSQLiteRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	const q = `
	INSERT INTO users (email, password_hash, created_at)
	VALUES (?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, q, u.Email, u.PasswordHash, u.CreatedAt.Format(time.RFC3339Nano))
	if err != nil {
		return domain.User{}, fmt.Errorf("create user failed: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, fmt.Errorf("create user last insert id failed: %w", err)
	}
	u.ID = id
	return u, nil
}

func (r *UserSQLiteRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	const q = `
	SELECT id, email, password_hash, created_at
	FROM users
	WHERE email = ?
	LIMIT 1
	`

	var u domain.User
	var createdAt string

	err := r.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &createdAt)
	if err == sql.ErrNoRows {
		return domain.User{}, repository.ErrNotFound
	}
	if err != nil {
		return domain.User{}, fmt.Errorf("find by email failed: %w", err)
	}

	t, err := time.Parse(time.RFC3339Nano, createdAt)
	if err != nil {
		// DB'de format bozulursa patlamasın; ama hata olarak dönelim
		return domain.User{}, fmt.Errorf("parse created_at failed: %w", err)
	}
	u.CreatedAt = t

	return u, nil
}
