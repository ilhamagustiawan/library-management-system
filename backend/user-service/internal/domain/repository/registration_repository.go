package repository

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type RegistrationRepository interface {
	Prepare(ctx context.Context, registration *entity.Registration) (*entity.Registration, error)
	Complete(ctx context.Context, registrationID string, user *entity.User, event *entity.OutboxEvent) error
	MarkConflict(ctx context.Context, registrationID string) error
}
