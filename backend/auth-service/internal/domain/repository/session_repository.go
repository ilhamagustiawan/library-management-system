package repository

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
}
