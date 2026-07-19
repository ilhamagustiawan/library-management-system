package oauth

import (
	"context"
	"net/http"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/repository"
)

type oauthUsecase struct {
	oauth repository.OAuthRepository
	users repository.UserRepository
}

func NewUsecase(
	oauth repository.OAuthRepository,
	users repository.UserRepository,
) Usecase {
	return &oauthUsecase{oauth: oauth, users: users}
}

func (u *oauthUsecase) Metadata() entity.OAuthMetadata {
	return u.oauth.Metadata()
}

func (u *oauthUsecase) UserInfo(ctx context.Context, token string) (*entity.User, error) {
	if token == "" {
		return nil, invalidAccessToken(nil)
	}
	userID, err := u.oauth.FindUserIDByAccessToken(ctx, token)
	if err != nil || userID == "" {
		return nil, invalidAccessToken(err)
	}
	return u.users.FindByID(ctx, userID)
}

func invalidAccessToken(cause error) error {
	return errs.New(http.StatusUnauthorized, errs.CodeInvalidToken, "invalid access token", nil, cause)
}
