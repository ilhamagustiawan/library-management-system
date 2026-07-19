package auth

import (
	"context"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	User         *entity.User
	SessionToken string
	ExpiresAt    time.Time
}

type Usecase interface {
	Register(ctx context.Context, input RegisterInput) (*entity.User, error)
	Login(ctx context.Context, input LoginInput) (*LoginResult, error)
	AuthenticateSession(ctx context.Context, token string) (*entity.User, error)
	FindUser(ctx context.Context, id string) (*entity.User, error)
	Logout(ctx context.Context, token string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenGenerator interface {
	Generate() (string, error)
}
