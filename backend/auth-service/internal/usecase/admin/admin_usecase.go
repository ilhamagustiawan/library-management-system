package admin

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

var ErrExistingMember = errors.New("email belongs to an existing member")

type Input struct {
	Name     string
	Email    string
	Password string
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type Usecase struct {
	users  UserRepository
	hasher PasswordHasher
	now    func() time.Time
	newID  func() string
}

func New(users UserRepository, hasher PasswordHasher) *Usecase {
	return &Usecase{users: users, hasher: hasher, now: func() time.Time { return time.Now().UTC() }, newID: uuid.NewString}
}

func (u *Usecase) Create(ctx context.Context, input Input) (*entity.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if len(name) < 2 || len(name) > 100 || len(email) > 254 || !strings.Contains(email, "@") {
		return nil, fmt.Errorf("valid --name and --email are required")
	}
	existing, err := u.users.FindByEmail(ctx, email)
	if err == nil {
		if existing.Role == entity.RoleAdmin {
			return existing, nil
		}
		return nil, ErrExistingMember
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return nil, fmt.Errorf("find admin identity: %w", err)
	}
	if len(input.Password) < 12 || len([]byte(input.Password)) > 72 {
		return nil, fmt.Errorf("password must contain 12 to 72 bytes")
	}
	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash admin password: %w", err)
	}
	now := u.now()
	user := &entity.User{
		ID: u.newID(), Name: name, Email: email, PasswordHash: passwordHash, Role: entity.RoleAdmin,
		CreatedAt: now, UpdatedAt: now,
	}
	if err := u.users.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create admin identity: %w", err)
	}
	return user, nil
}
