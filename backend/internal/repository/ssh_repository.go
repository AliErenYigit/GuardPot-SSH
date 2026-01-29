package repository

import (
	"context"

	"backend/internal/domain"
)

type SSHConnectionRepository interface {
	Create(ctx context.Context, c domain.SSHConnection) (domain.SSHConnection, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.SSHConnection, error)
	GetByID(ctx context.Context, userID, id int64) (domain.SSHConnection, error)
	Delete(ctx context.Context, userID, id int64) error
}
