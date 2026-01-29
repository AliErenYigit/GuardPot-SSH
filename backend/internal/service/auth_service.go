package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"backend/internal/domain"
	"backend/internal/repository"
)

type AuthService struct {
	users    repository.UserRepository
	password PasswordManager
	tokens   TokenManager
}

func NewAuthService(users repository.UserRepository, password PasswordManager, tokens TokenManager) AuthService {
	return AuthService{
		users:    users,
		password: password,
		tokens:   tokens,
	}
}

type PublicUser struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginResult struct {
	AccessToken string     `json:"accessToken"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	User        PublicUser `json:"user"`
}

func (a AuthService) Register(ctx context.Context, email string, password string) (PublicUser, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if !isValidEmail(email) {
		return PublicUser{}, ErrInvalidEmail
	}
	if len(password) < 8 {
		return PublicUser{}, ErrInvalidPassword
	}

	exists, err := a.users.ExistsByEmail(ctx, email)
	if err != nil {
		return PublicUser{}, fmt.Errorf("exists check failed: %w", err)
	}
	if exists {
		return PublicUser{}, ErrEmailAlreadyInUse
	}

	hash, err := a.password.Hash(password)
	if err != nil {
		return PublicUser{}, err
	}

	u := domain.User{
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
	}

	created, err := a.users.Create(ctx, u)
	if err != nil {
		// race condition: unique constraint (iki istek aynı anda)
		if isUniqueConstraintErr(err) {
			return PublicUser{}, ErrEmailAlreadyInUse
		}
		return PublicUser{}, err
	}

	return PublicUser{
		ID:        created.ID,
		Email:     created.Email,
		CreatedAt: created.CreatedAt,
	}, nil
}

func (a AuthService) Login(ctx context.Context, email string, password string) (LoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if !isValidEmail(email) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if password == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	u, err := a.users.FindByEmail(ctx, email)
	if err != nil {
		if err == repository.ErrNotFound {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, fmt.Errorf("find user failed: %w", err)
	}

	if !a.password.Compare(u.PasswordHash, password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	tok, err := a.tokens.GenerateAccessToken(u.ID, u.Email)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken: tok.Token,
		ExpiresAt:   tok.ExpiresAt,
		User: PublicUser{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		},
	}, nil
}

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(s string) bool {
	return emailRe.MatchString(s)
}

// SQLite unique hatası basit tespit (driver bağımlı olmadan)
func isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique") || strings.Contains(msg, "constraint")
}
