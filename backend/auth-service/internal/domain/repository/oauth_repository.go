package repository

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type OAuthRepository interface {
	FindUserIDByAccessToken(ctx context.Context, token string) (string, error)
	Metadata() entity.OAuthMetadata
}
