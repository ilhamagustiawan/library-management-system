package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/repository"
)

type Config struct {
	SessionTTL        time.Duration
	DummyPasswordHash string
}

type authUsecase struct {
	users    repository.UserRepository
	sessions repository.SessionRepository
	hasher   PasswordHasher
	tokens   TokenGenerator
	config   Config
	now      func() time.Time
	newID    func() string
}

func NewAuthUsecase(
	users repository.UserRepository,
	sessions repository.SessionRepository,
	hasher PasswordHasher,
	tokens TokenGenerator,
	config Config,
) Usecase {
	return &authUsecase{
		users: users, sessions: sessions, hasher: hasher, tokens: tokens, config: config,
		now: utcNow, newID: uuid.NewString,
	}
}

func utcNow() time.Time {
	return time.Now().UTC()
}

func (u *authUsecase) Register(ctx context.Context, input RegisterInput) (*entity.User, error) {
	name := strings.TrimSpace(input.Name)
	email := normalizeEmail(input.Email)
	if name == "" || email == "" || len(input.Password) < 12 || len([]byte(input.Password)) > 72 {
		return nil, errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, "invalid registration data", nil, nil)
	}

	passwordHash, err := u.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	now := u.now()
	user := &entity.User{
		ID: u.newID(), Name: name, Email: email, PasswordHash: passwordHash, Role: entity.RoleMember,
		CreatedAt: now, UpdatedAt: now,
	}
	if err := u.users.Create(ctx, user); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return nil, errs.New(http.StatusConflict, errs.CodeEmailExists, "email is already registered", nil, err)
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}

func (u *authUsecase) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	user, err := u.users.FindByEmail(ctx, normalizeEmail(input.Email))
	if errors.Is(err, errs.ErrNotFound) {
		_ = u.hasher.Compare(u.config.DummyPasswordHash, input.Password)
		return nil, invalidCredentials(err)
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if err := u.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, invalidCredentials(err)
	}

	rawToken, err := u.tokens.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}
	now := u.now()
	session := &entity.Session{
		ID: u.newID(), UserID: user.ID, TokenHash: sessionTokenHash(rawToken),
		CreatedAt: now, ExpiresAt: now.Add(u.config.SessionTTL),
	}
	if err := u.sessions.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &LoginResult{User: user, SessionToken: rawToken, ExpiresAt: session.ExpiresAt}, nil
}

func (u *authUsecase) AuthenticateSession(ctx context.Context, token string) (*entity.User, error) {
	if token == "" {
		return nil, invalidSession(nil)
	}
	hash := sessionTokenHash(token)
	session, err := u.sessions.FindByTokenHash(ctx, hash)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, invalidSession(err)
	}
	if err != nil {
		return nil, fmt.Errorf("find session: %w", err)
	}
	if !session.ExpiresAt.After(u.now()) {
		_ = u.sessions.DeleteByTokenHash(ctx, hash)
		return nil, invalidSession(nil)
	}
	user, err := u.users.FindByID(ctx, session.UserID)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, invalidSession(err)
	}
	if err != nil {
		return nil, fmt.Errorf("find session user: %w", err)
	}
	return user, nil
}

func (u *authUsecase) FindUser(ctx context.Context, id string) (*entity.User, error) {
	return u.users.FindByID(ctx, id)
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	if token == "" {
		return nil
	}
	if err := u.sessions.DeleteByTokenHash(ctx, sessionTokenHash(token)); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func sessionTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func invalidCredentials(cause error) error {
	return errs.New(http.StatusUnauthorized, errs.CodeInvalidCredentials, "invalid email or password", nil, cause)
}

func invalidSession(cause error) error {
	return errs.New(http.StatusUnauthorized, errs.CodeInvalidToken, "authentication required", nil, cause)
}
