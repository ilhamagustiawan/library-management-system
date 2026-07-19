package oauth

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	oauthusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/oauth"
)

type Handler struct {
	usecase oauthusecase.Usecase
}

func NewHandler(usecase oauthusecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// Metadata returns OAuth authorization server metadata.
//
// @Summary Get authorization server metadata
// @Description Returns RFC 8414 discovery metadata for this OAuth server.
// @Tags OAuth
// @Produce json
// @Success 200 {object} response.OAuthMetadata "Authorization server metadata"
// @Router /.well-known/oauth-authorization-server [get]
func (h *Handler) Metadata(c *fiber.Ctx) error {
	return c.JSON(response.NewOAuthMetadata(h.usecase.Metadata()))
}

// UserInfo returns the resource owner represented by an access token.
//
// @Summary Get token user
// @Tags OAuth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserSuccess "Token user"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired access token"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Failure 500 {object} response.ErrorResponse "Internal error"
// @Router /api/v1/oauth/userinfo [get]
func (h *Handler) UserInfo(c *fiber.Ctx) error {
	user, err := h.usecase.UserInfo(c.UserContext(), helper.BearerToken(c.Get(fiber.HeaderAuthorization)))
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, response.NewUser(user))
}
