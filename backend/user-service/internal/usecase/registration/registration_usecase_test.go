package registration

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
)

type fakeRepository struct {
	prepared  *entity.Registration
	completed *entity.User
	event     *entity.OutboxEvent
	conflicts int
}

func (r *fakeRepository) Prepare(_ context.Context, candidate *entity.Registration) (*entity.Registration, error) {
	if r.prepared != nil {
		copy := *r.prepared
		return &copy, nil
	}
	copy := *candidate
	r.prepared = &copy
	return &copy, nil
}

func (r *fakeRepository) Complete(_ context.Context, _ string, user *entity.User, event *entity.OutboxEvent) error {
	userCopy := *user
	eventCopy := *event
	r.completed = &userCopy
	r.event = &eventCopy
	return nil
}

func (r *fakeRepository) MarkConflict(context.Context, string) error {
	r.conflicts++
	return nil
}

type fakeIdentityClient struct {
	identity *Identity
	err      error
	calls    int
	key      string
}

func (c *fakeIdentityClient) Create(_ context.Context, key string, _ IdentityInput) (*Identity, error) {
	c.calls++
	c.key = key
	return c.identity, c.err
}

func TestRegisterCompletesProfileAndSecretFreeEvent(t *testing.T) {
	repository := &fakeRepository{}
	identities := &fakeIdentityClient{identity: &Identity{
		ID: "user-123", Name: "Ada Lovelace", Email: "ada@example.com", Role: entity.RoleMember,
	}}
	usecase := New(repository, identities)
	usecase.now = func() time.Time { return time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC) }
	usecase.newID = func() string { return "registration-123" }

	user, err := usecase.Register(context.Background(), Input{
		Name: "  Ada Lovelace  ", Email: " ADA@example.com ", Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.ID != "user-123" || user.Role != entity.RoleMember || identities.key != "registration-123" {
		t.Fatalf("user = %#v, key = %q", user, identities.key)
	}
	var event entity.UserRegistered
	if err := json.Unmarshal(repository.event.Payload, &event); err != nil {
		t.Fatalf("event payload: %v", err)
	}
	if event.EventID != "registration-123" || event.Data.UserID != "user-123" || event.Data.Role != entity.RoleMember {
		t.Fatalf("event = %#v", event)
	}
	if string(repository.event.Payload) == "" || containsSecret(string(repository.event.Payload)) {
		t.Fatalf("unsafe event payload = %s", repository.event.Payload)
	}
}

func TestRegisterRejectsCompletedEmailWithoutIdentityCall(t *testing.T) {
	repository := &fakeRepository{prepared: &entity.Registration{
		ID: "registration-123", Name: "Ada", Email: "ada@example.com", Status: entity.RegistrationCompleted,
	}}
	identities := &fakeIdentityClient{}
	_, err := New(repository, identities).Register(context.Background(), Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeEmailExists || identities.calls != 0 {
		t.Fatalf("error = %v, identity calls = %d", err, identities.calls)
	}
}

func TestRegisterLeavesPendingOperationWhenIdentityUnavailable(t *testing.T) {
	repository := &fakeRepository{}
	identities := &fakeIdentityClient{err: errs.ErrIdentityUnavailable}
	_, err := New(repository, identities).Register(context.Background(), Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeDependency {
		t.Fatalf("error = %v", err)
	}
	if repository.completed != nil || repository.conflicts != 0 {
		t.Fatalf("pending operation changed: %#v", repository)
	}
}

func TestRegisterMarksIdentityConflict(t *testing.T) {
	repository := &fakeRepository{}
	identities := &fakeIdentityClient{err: errs.ErrIdentityConflict}
	_, err := New(repository, identities).Register(context.Background(), Input{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeEmailExists || repository.conflicts != 1 {
		t.Fatalf("error = %v, conflicts = %d", err, repository.conflicts)
	}
}

func TestRegisterPreservesPendingOperationAfterChangedIdentityReplay(t *testing.T) {
	repository := &fakeRepository{}
	identities := &fakeIdentityClient{err: errs.ErrIdentityIdempotency}
	_, err := New(repository, identities).Register(context.Background(), Input{
		Name: "Ada", Email: "ada@example.com", Password: "different secure password",
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeIdempotencyConflict || repository.conflicts != 0 {
		t.Fatalf("error = %v, conflicts = %d", err, repository.conflicts)
	}
}

func containsSecret(payload string) bool {
	for _, forbidden := range []string{"password", "correct horse", "email", "name"} {
		if len(payload) >= len(forbidden) {
			for start := 0; start+len(forbidden) <= len(payload); start++ {
				if payload[start:start+len(forbidden)] == forbidden {
					return true
				}
			}
		}
	}
	return false
}
