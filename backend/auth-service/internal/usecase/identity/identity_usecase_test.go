package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type fakeRepository struct {
	created *entity.User
	key     string
	err     error
}

func (r *fakeRepository) Create(_ context.Context, key string, user *entity.User) (*entity.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.created != nil && r.key == key {
		copy := *r.created
		return &copy, nil
	}
	r.key = key
	copy := *user
	r.created = &copy
	return &copy, nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) { return "hashed:" + password, nil }
func (fakeHasher) Compare(hash, password string) error {
	if hash != "hashed:"+password {
		return errors.New("password mismatch")
	}
	return nil
}

func TestCreateAssignsMemberAndStoresNoPlaintextPassword(t *testing.T) {
	repository := &fakeRepository{}
	usecase := New(repository, fakeHasher{})
	user, err := usecase.Create(context.Background(), "registration-123", Input{
		Name: "  Ada Lovelace  ", Email: " ADA@example.com ", Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.Role != entity.RoleMember || user.Email != "ada@example.com" || user.Name != "Ada Lovelace" {
		t.Fatalf("user = %#v", user)
	}
	if repository.key != "registration-123" || repository.created.PasswordHash != "hashed:correct horse battery staple" {
		t.Fatalf("stored identity = %#v, key = %q", repository.created, repository.key)
	}
}

func TestCreateRejectsInvalidIdentityBeforePersistence(t *testing.T) {
	for _, input := range []Input{
		{Name: "", Email: "ada@example.com", Password: "correct horse battery staple"},
		{Name: "Ada", Email: "", Password: "correct horse battery staple"},
		{Name: "Ada", Email: "ada@example.com", Password: "short"},
	} {
		repository := &fakeRepository{}
		_, err := New(repository, fakeHasher{}).Create(context.Background(), "registration-123", input)
		if err == nil || repository.created != nil {
			t.Fatalf("Create(%#v) = %v, persisted = %#v", input, err, repository.created)
		}
	}
}

func TestCreatePropagatesRepositoryConflict(t *testing.T) {
	want := errors.New("conflict")
	_, err := New(&fakeRepository{err: want}, fakeHasher{}).Create(context.Background(), "registration-123", Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	if !errors.Is(err, want) {
		t.Fatalf("Create() error = %v, want conflict", err)
	}
}

func TestCreateRejectsIdempotencyReplayWithChangedIdentityData(t *testing.T) {
	repository := &fakeRepository{}
	usecase := New(repository, fakeHasher{})
	_, err := usecase.Create(context.Background(), "registration-123", Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("first Create() error = %v", err)
	}

	for name, input := range map[string]Input{
		"name":     {Name: "Grace", Email: "ada@example.com", Password: "correct horse battery staple"},
		"password": {Name: "Ada", Email: "ada@example.com", Password: "different secure password"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := usecase.Create(context.Background(), "registration-123", input)
			if err == nil {
				t.Fatal("Create() error = nil, want idempotency conflict")
			}
		})
	}
}
