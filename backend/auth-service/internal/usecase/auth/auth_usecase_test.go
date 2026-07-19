package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type fakeUserRepository struct {
	usersByID    map[string]*entity.User
	usersByEmail map[string]*entity.User
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{usersByID: make(map[string]*entity.User), usersByEmail: make(map[string]*entity.User)}
}

func (r *fakeUserRepository) Create(_ context.Context, user *entity.User) error {
	if _, exists := r.usersByEmail[user.Email]; exists {
		return errs.ErrConflict
	}
	copy := *user
	r.usersByID[user.ID] = &copy
	r.usersByEmail[user.Email] = &copy
	return nil
}

func (r *fakeUserRepository) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	user, exists := r.usersByEmail[email]
	if !exists {
		return nil, errs.ErrNotFound
	}
	copy := *user
	return &copy, nil
}

func (r *fakeUserRepository) FindByID(_ context.Context, id string) (*entity.User, error) {
	user, exists := r.usersByID[id]
	if !exists {
		return nil, errs.ErrNotFound
	}
	copy := *user
	return &copy, nil
}

type fakeSessionRepository struct {
	byHash map[string]*entity.Session
}

func newFakeSessionRepository() *fakeSessionRepository {
	return &fakeSessionRepository{byHash: make(map[string]*entity.Session)}
}

func (r *fakeSessionRepository) Create(_ context.Context, session *entity.Session) error {
	copy := *session
	r.byHash[session.TokenHash] = &copy
	return nil
}

func (r *fakeSessionRepository) FindByTokenHash(_ context.Context, hash string) (*entity.Session, error) {
	session, exists := r.byHash[hash]
	if !exists {
		return nil, errs.ErrNotFound
	}
	copy := *session
	return &copy, nil
}

func (r *fakeSessionRepository) DeleteByTokenHash(_ context.Context, hash string) error {
	delete(r.byHash, hash)
	return nil
}

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(password string) (string, error) { return "hashed:" + password, nil }
func (fakePasswordHasher) Compare(hash, password string) error {
	if hash != "hashed:"+password {
		return errors.New("password mismatch")
	}
	return nil
}

type fakeTokenGenerator struct{ token string }

func (g fakeTokenGenerator) Generate() (string, error) { return g.token, nil }

func TestRegisterNormalizesIdentityAndHashesPassword(t *testing.T) {
	users := newFakeUserRepository()
	uc := NewAuthUsecase(users, newFakeSessionRepository(), fakePasswordHasher{}, fakeTokenGenerator{"session-token"}, testConfig())

	user, err := uc.Register(context.Background(), RegisterInput{
		Name: "  Ada Lovelace  ", Email: "  ADA@Example.COM ", Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.Name != "Ada Lovelace" || user.Email != "ada@example.com" {
		t.Fatalf("Register() user = %#v, want normalized identity", user)
	}
	if user.PasswordHash != "hashed:correct horse battery staple" {
		t.Fatalf("password hash = %q, want generated hash", user.PasswordHash)
	}
	if user.Role != entity.RoleMember {
		t.Fatalf("role = %q, want member", user.Role)
	}
}

func TestRegisterMapsDuplicateEmailToConflict(t *testing.T) {
	uc := NewAuthUsecase(newFakeUserRepository(), newFakeSessionRepository(), fakePasswordHasher{}, fakeTokenGenerator{"token"}, testConfig())
	input := RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "a secure password"}
	_, _ = uc.Register(context.Background(), input)

	_, err := uc.Register(context.Background(), input)
	assertDomainCode(t, err, errs.CodeEmailExists)
}

func TestLoginCreatesHashedOpaqueSession(t *testing.T) {
	users := newFakeUserRepository()
	sessions := newFakeSessionRepository()
	uc := NewAuthUsecase(users, sessions, fakePasswordHasher{}, fakeTokenGenerator{"raw-session-token"}, testConfig())
	_, _ = uc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "right password"})

	result, err := uc.Login(context.Background(), LoginInput{Email: " ADA@example.com ", Password: "right password"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.SessionToken != "raw-session-token" {
		t.Fatalf("session token = %q, want raw token returned once", result.SessionToken)
	}
	if _, storedRaw := sessions.byHash["raw-session-token"]; storedRaw {
		t.Fatal("session repository stored the raw session token")
	}
	if len(sessions.byHash) != 1 {
		t.Fatalf("stored sessions = %d, want 1", len(sessions.byHash))
	}
}

func TestLoginUsesSamePublicErrorForUnknownEmailAndWrongPassword(t *testing.T) {
	users := newFakeUserRepository()
	uc := NewAuthUsecase(users, newFakeSessionRepository(), fakePasswordHasher{}, fakeTokenGenerator{"token"}, testConfig())
	_, _ = uc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "right password"})

	for name, input := range map[string]LoginInput{
		"unknown email": {Email: "missing@example.com", Password: "right password"},
		"wrong secret":  {Email: "ada@example.com", Password: "wrong password"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := uc.Login(context.Background(), input)
			assertDomainCode(t, err, errs.CodeInvalidCredentials)
		})
	}
}

func TestAuthenticateAndLogoutSession(t *testing.T) {
	users := newFakeUserRepository()
	sessions := newFakeSessionRepository()
	uc := NewAuthUsecase(users, sessions, fakePasswordHasher{}, fakeTokenGenerator{"raw-session-token"}, testConfig())
	created, _ := uc.Register(context.Background(), RegisterInput{Name: "Ada", Email: "ada@example.com", Password: "right password"})
	_, _ = uc.Login(context.Background(), LoginInput{Email: "ada@example.com", Password: "right password"})

	user, err := uc.AuthenticateSession(context.Background(), "raw-session-token")
	if err != nil || user.ID != created.ID {
		t.Fatalf("AuthenticateSession() = (%#v, %v), want created user", user, err)
	}
	if err := uc.Logout(context.Background(), "raw-session-token"); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	_, err = uc.AuthenticateSession(context.Background(), "raw-session-token")
	assertDomainCode(t, err, errs.CodeInvalidToken)
}

func testConfig() Config {
	return Config{SessionTTL: 24 * 60 * 60 * 1e9, DummyPasswordHash: "hashed:not-the-user-password"}
}

func assertDomainCode(t *testing.T, err error, code string) {
	t.Helper()
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != code {
		t.Fatalf("error = %v, want domain code %s", err, code)
	}
}
