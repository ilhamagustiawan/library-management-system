package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type fakeUsers struct{ user *entity.User }

func (f *fakeUsers) FindByEmail(context.Context, string) (*entity.User, error) {
	if f.user == nil {
		return nil, errs.ErrNotFound
	}
	copy := *f.user
	return &copy, nil
}
func (f *fakeUsers) Create(_ context.Context, user *entity.User) error {
	copy := *user
	f.user = &copy
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) { return "hashed:" + password, nil }

func TestCreateCreatesAdminWithoutPersistingPlaintext(t *testing.T) {
	users := &fakeUsers{}
	created, err := New(users, fakeHasher{}).Create(context.Background(), Input{
		Name: "Ada Admin", Email: " ADA@example.com ", Password: "correct horse battery staple",
	})
	if err != nil || created.Role != entity.RoleAdmin || users.user.PasswordHash != "hashed:correct horse battery staple" {
		t.Fatalf("Create() = (%#v, %v), stored = %#v", created, err, users.user)
	}
}

func TestCreateReturnsExistingAdminWithoutMutation(t *testing.T) {
	existing := &entity.User{ID: "admin-1", Email: "ada@example.com", Role: entity.RoleAdmin, PasswordHash: "original"}
	users := &fakeUsers{user: existing}
	created, err := New(users, fakeHasher{}).Create(context.Background(), Input{
		Name: "Changed", Email: "ada@example.com", Password: "different secure password",
	})
	if err != nil || created.PasswordHash != "original" || users.user.PasswordHash != "original" {
		t.Fatalf("Create() = (%#v, %v), stored = %#v", created, err, users.user)
	}
}

func TestCreateRejectsExistingMemberWithoutPromotion(t *testing.T) {
	users := &fakeUsers{user: &entity.User{ID: "member-1", Email: "ada@example.com", Role: entity.RoleMember}}
	_, err := New(users, fakeHasher{}).Create(context.Background(), Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	if !errors.Is(err, ErrExistingMember) || users.user.Role != entity.RoleMember {
		t.Fatalf("Create() error = %v, user = %#v", err, users.user)
	}
}

func TestCreateRejectsInvalidInputBeforePersistence(t *testing.T) {
	for _, input := range []Input{
		{Name: "", Email: "ada@example.com", Password: "correct horse battery staple"},
		{Name: "Ada", Email: "invalid", Password: "correct horse battery staple"},
		{Name: "Ada", Email: "ada@example.com", Password: "short"},
	} {
		users := &fakeUsers{}
		if _, err := New(users, fakeHasher{}).Create(context.Background(), input); err == nil || users.user != nil {
			t.Fatalf("Create(%#v) persisted %#v, want validation failure", input, users.user)
		}
	}
}
