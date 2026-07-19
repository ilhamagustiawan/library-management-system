package identity

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/request"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
	oauthinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/oauth"
	identityusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/identity"
)

type Creator interface {
	Create(ctx context.Context, idempotencyKey string, input identityusecase.Input) (*entity.User, error)
}

type TokenAuthenticator interface {
	AuthenticateServiceToken(ctx context.Context, token, expectedClientID, expectedAudience, requiredScope string) error
}

type Handler struct {
	creator       Creator
	authenticator TokenAuthenticator
	validate      *validator.Validate
}

func NewHandler(creator Creator, authenticator TokenAuthenticator, validate *validator.Validate) *Handler {
	return &Handler{creator: creator, authenticator: authenticator, validate: validate}
}

// Create creates a member identity for User Service.
//
// @Summary Create member identity
// @Description Internal idempotent endpoint. Requires a User Service bearer token with identities:create.
// @Tags Internal
// @Accept json
// @Produce json
// @Param Idempotency-Key header string true "Stable registration command key"
// @Param request body request.Register true "Member identity"
// @Success 201 {object} response.UserSuccess "Identity created or replayed"
// @Failure 401 {object} response.ErrorResponse "Invalid service token"
// @Failure 403 {object} response.ErrorResponse "Service scope denied"
// @Failure 409 {object} response.ErrorResponse "Email or idempotency conflict"
// @Failure 422 {object} response.ErrorResponse "Invalid identity data"
// @Router /internal/identities [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	token, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
	if !ok {
		return response.Error(c, errs.New(http.StatusUnauthorized, errs.CodeInvalidToken, "valid User Service token required; no identity was changed", nil, nil))
	}
	if err := h.authenticator.AuthenticateServiceToken(c.UserContext(), token, "user-service", "auth-service", "identities:create"); err != nil {
		if errors.Is(err, oauthinfra.ErrInsufficientServiceGrant) {
			return response.Error(c, errs.New(http.StatusForbidden, errs.CodeForbidden, "User Service token lacks identities:create; no identity was changed", nil, err))
		}
		return response.Error(c, errs.New(http.StatusUnauthorized, errs.CodeInvalidToken, "User Service token is invalid; no identity was changed", nil, err))
	}

	var input request.Register
	if err := request.DecodeStrictJSON(c, &input); err != nil || h.validate.Struct(input) != nil {
		return response.ValidationError(c, "invalid identity data; no identity was changed")
	}
	user, err := h.creator.Create(c.UserContext(), c.Get("Idempotency-Key"), identityusecase.Input{
		Name: input.Name, Email: input.Email, Password: input.Password,
	})
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusCreated, response.NewUser(user))
}

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)
	returnValue := ""
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		returnValue = parts[1]
	}
	return returnValue, returnValue != ""
}
