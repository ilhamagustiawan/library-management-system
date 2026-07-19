package oauth

import (
	"context"
	"errors"
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type fakeAccessTokenRepository struct {
	userID string
	err    error
}

func (r fakeAccessTokenRepository) FindUserIDByAccessToken(context.Context, string) (string, error) {
	return r.userID, r.err
}

func (r fakeAccessTokenRepository) Metadata() entity.OAuthMetadata {
	return entity.NewOAuthMetadata("http://auth.example/", []string{"books:read"})
}

type fakeUserRepository struct {
	user *entity.User
	err  error
}

func (r fakeUserRepository) Create(context.Context, *entity.User) error { return nil }
func (r fakeUserRepository) FindByEmail(context.Context, string) (*entity.User, error) {
	return nil, errs.ErrNotFound
}
func (r fakeUserRepository) FindByID(context.Context, string) (*entity.User, error) {
	return r.user, r.err
}

func TestUserInfoResolvesAccessTokenUser(t *testing.T) {
	user := &entity.User{ID: "user-1", Email: "ada@example.com"}
	uc := NewUsecase(
		fakeAccessTokenRepository{userID: user.ID},
		fakeUserRepository{user: user},
	)

	got, err := uc.UserInfo(context.Background(), "access-token")
	if err != nil || got.ID != user.ID {
		t.Fatalf("UserInfo() = (%#v, %v), want user", got, err)
	}
}

func TestUserInfoMapsInvalidAccessTokenToDomainError(t *testing.T) {
	uc := NewUsecase(
		fakeAccessTokenRepository{err: errors.New("invalid token")},
		fakeUserRepository{},
	)

	_, err := uc.UserInfo(context.Background(), "invalid-token")
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeInvalidToken {
		t.Fatalf("UserInfo() error = %v, want %s", err, errs.CodeInvalidToken)
	}
}

func TestMetadataUsesNormalizedIssuer(t *testing.T) {
	uc := NewUsecase(fakeAccessTokenRepository{}, fakeUserRepository{})

	metadata := uc.Metadata()
	if metadata.Issuer != "http://auth.example" || metadata.TokenEndpoint != "http://auth.example/oauth/token" {
		t.Fatalf("Metadata() = %#v, want normalized endpoints", metadata)
	}
}
