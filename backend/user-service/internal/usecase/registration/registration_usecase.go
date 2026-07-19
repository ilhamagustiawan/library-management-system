package registration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/repository"
)

type registrationUsecase struct {
	repository repository.RegistrationRepository
	identities IdentityCreator
	now        func() time.Time
	newID      func() string
}

func New(repository repository.RegistrationRepository, identities IdentityCreator) *registrationUsecase {
	return &registrationUsecase{
		repository: repository,
		identities: identities,
		now:        func() time.Time { return time.Now().UTC() },
		newID:      uuid.NewString,
	}
}

func (u *registrationUsecase) Register(ctx context.Context, input Input) (*entity.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if name == "" || email == "" || len(input.Password) < 12 || len([]byte(input.Password)) > 72 {
		return nil, validationError()
	}

	now := u.now()
	registration, err := u.repository.Prepare(ctx, &entity.Registration{
		ID: u.newID(), Name: name, Email: email, Status: entity.RegistrationPending,
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("prepare registration: %w", err)
	}
	if registration.Name != name || registration.Email != email {
		return nil, errs.New(
			http.StatusConflict, errs.CodeIdempotencyConflict,
			"registration details changed while a prior attempt remains pending; retry with the original name and email",
			nil, nil,
		)
	}
	if registration.Status != entity.RegistrationPending {
		return nil, emailExists(nil)
	}

	identity, err := u.identities.Create(ctx, registration.ID, IdentityInput{
		Name: name, Email: email, Password: input.Password,
	})
	if err != nil {
		return nil, u.handleIdentityError(ctx, registration.ID, err)
	}
	if identity == nil || identity.ID == "" || identity.Name != name || identity.Email != email || identity.Role != entity.RoleMember {
		return nil, dependencyError(errs.ErrInvalidIdentity)
	}

	user := &entity.User{
		ID: identity.ID, Name: identity.Name, Email: identity.Email, Role: identity.Role,
		CreatedAt: now, UpdatedAt: now,
	}
	eventPayload, err := json.Marshal(entity.UserRegistered{
		EventID: registration.ID, Type: entity.UserRegisteredV1, OccurredAt: now,
		Data: entity.UserRegisteredData{UserID: user.ID, Role: user.Role},
	})
	if err != nil {
		return nil, fmt.Errorf("encode user registration event: %w", err)
	}
	if err := u.repository.Complete(ctx, registration.ID, user, &entity.OutboxEvent{
		ID: registration.ID, Type: entity.UserRegisteredV1, AggregateID: user.ID,
		Payload: eventPayload, OccurredAt: now,
	}); err != nil {
		if errors.Is(err, errs.ErrConflict) {
			return nil, emailExists(err)
		}
		return nil, fmt.Errorf("complete registration: %w", err)
	}
	return user, nil
}

func (u *registrationUsecase) handleIdentityError(ctx context.Context, registrationID string, err error) error {
	if errors.Is(err, errs.ErrIdentityIdempotency) {
		return errs.New(
			http.StatusConflict, errs.CodeIdempotencyConflict,
			"registration retry changed identity data; pending registration remains preserved, retry with the original details",
			nil, err,
		)
	}
	if errors.Is(err, errs.ErrIdentityConflict) {
		if markErr := u.repository.MarkConflict(ctx, registrationID); markErr != nil {
			return fmt.Errorf("mark registration conflict: %w", markErr)
		}
		return emailExists(err)
	}
	if errors.Is(err, errs.ErrIdentityUnavailable) || errors.Is(err, errs.ErrInvalidIdentity) {
		return dependencyError(err)
	}
	return fmt.Errorf("create identity: %w", err)
}

func validationError() error {
	return errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, "invalid registration data", nil, nil)
}

func emailExists(cause error) error {
	return errs.New(http.StatusConflict, errs.CodeEmailExists, "email is already registered; existing account remains unchanged", nil, cause)
}

func dependencyError(cause error) error {
	return errs.New(
		http.StatusServiceUnavailable, errs.CodeDependency,
		"registration dependency unavailable; pending registration was preserved, retry with the same details",
		nil, cause,
	)
}
