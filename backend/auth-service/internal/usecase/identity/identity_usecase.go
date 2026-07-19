package identity

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type Input struct {
	Name     string
	Email    string
	Password string
}

type Repository interface {
	Create(ctx context.Context, idempotencyKey string, user *entity.User) (*entity.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type Usecase struct {
	repository Repository
	hasher     PasswordHasher
	now        func() time.Time
	newID      func() string
}

func New(repository Repository, hasher PasswordHasher) *Usecase {
	return &Usecase{repository: repository, hasher: hasher, now: func() time.Time { return time.Now().UTC() }, newID: uuid.NewString}
}

func (u *Usecase) Create(ctx context.Context, idempotencyKey string, input Input) (*entity.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	key := strings.TrimSpace(idempotencyKey)
	if len(key) < 8 || len(key) > 255 || len(name) < 2 || len(name) > 100 || len(email) > 254 || !strings.Contains(email, "@") || len(input.Password) < 12 || len([]byte(input.Password)) > 72 {
		return nil, errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, "invalid identity data or idempotency key", nil, nil)
	}
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return nil, err
	}
	now := u.now()
	user, err := u.repository.Create(ctx, key, &entity.User{
		ID: u.newID(), Name: name, Email: email, PasswordHash: passwordHash, Role: entity.RoleMember,
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	if user.Name != name || user.Email != email || u.hasher.Compare(user.PasswordHash, input.Password) != nil {
		return nil, errs.New(
			http.StatusConflict,
			errs.CodeIdempotencyConflict,
			"idempotency key was already used with different identity data; existing identity remains unchanged",
			nil,
			nil,
		)
	}
	return user, nil
}
