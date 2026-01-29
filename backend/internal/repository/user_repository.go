package repository

import (
	"context"

	"backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
