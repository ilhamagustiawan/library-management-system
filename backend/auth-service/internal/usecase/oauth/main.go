package oauth

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type Usecase interface {
	Metadata() entity.OAuthMetadata
	UserInfo(ctx context.Context, token string) (*entity.User, error)
}
