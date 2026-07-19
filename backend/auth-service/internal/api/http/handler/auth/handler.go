package auth

import (
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/request"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	authusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/auth"
)

type Config struct {
	Issuer              string
	SessionCookieName   string
	SessionCookieDomain string
	SessionCookieSecure bool
}

type Handler interface {
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
}

type handler struct {
	usecase       authusecase.Usecase
	validate      *validator.Validate
	issuer        string
	sessionCookie helper.SessionCookieConfig
}

func NewHandler(usecase authusecase.Usecase, validate *validator.Validate, config Config) Handler {
	return &handler{
		usecase:  usecase,
		validate: validate,
		issuer:   config.Issuer,
		sessionCookie: helper.SessionCookieConfig{
			Name: config.SessionCookieName, Domain: config.SessionCookieDomain, Secure: config.SessionCookieSecure,
		},
	}
}

// Register creates a user account.
//
// @Summary Register user
// @Description Creates a user account. Requests from browsers must use the configured trusted Origin.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.Register true "Registration details"
// @Success 201 {object} response.UserSuccess "User created"
// @Failure 403 {object} response.ErrorResponse "Untrusted request origin"
// @Failure 409 {object} response.ErrorResponse "Email already registered"
// @Failure 422 {object} response.ErrorResponse "Invalid registration data"
// @Failure 429 {object} response.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "Internal error"
// @Router /api/v1/auth/register [post]
func (h *handler) Register(c *fiber.Ctx) error {
	var input request.Register
	if err := request.DecodeStrictJSON(c, &input); err != nil || h.validate.Struct(input) != nil {
		return response.ValidationError(c, "invalid registration data")
	}
	user, err := h.usecase.Register(c.UserContext(), authusecase.RegisterInput{
		Name: input.Name, Email: input.Email, Password: input.Password,
	})
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusCreated, response.NewUser(user))
}

// Login creates an auth-service session.
//
// @Summary Log in
// @Description Authenticates a user and sets an HttpOnly session cookie. HTML form requests are also accepted. A valid returnTo redirects with 303.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body request.Login true "Login details"
// @Success 200 {object} response.UserSuccess "Authenticated user"
// @Success 303 {string} string "Redirect to authorization flow"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 403 {object} response.ErrorResponse "Untrusted request origin"
// @Failure 422 {object} response.ErrorResponse "Invalid login data or return URL"
// @Failure 429 {object} response.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "Internal error"
// @Router /api/v1/auth/login [post]
func (h *handler) Login(c *fiber.Ctx) error {
	var input request.Login
	var err error
	if c.Is("json") {
		err = request.DecodeStrictJSON(c, &input)
	} else {
		err = c.BodyParser(&input)
	}
	if err != nil || h.validate.Struct(input) != nil {
		return response.ValidationError(c, "invalid login data")
	}
	if input.ReturnTo != "" && !helper.IsSafeAuthorizeURL(input.ReturnTo, h.issuer) {
		return response.ValidationError(c, "invalid return_to URL")
	}

	result, err := h.usecase.Login(c.UserContext(), authusecase.LoginInput{Email: input.Email, Password: input.Password})
	if err != nil {
		return response.Error(c, err)
	}
	helper.SetSessionCookie(c, h.sessionCookie, result.SessionToken, result.ExpiresAt)
	if input.ReturnTo != "" {
		return c.Redirect(input.ReturnTo, http.StatusSeeOther)
	}
	return response.Success(c, http.StatusOK, response.NewUser(result.User))
}

// Logout revokes the current session.
//
// @Summary Log out
// @Description Revokes the session identified by the lms_session cookie and clears that cookie.
// @Tags Authentication
// @Param Cookie header string true "Session cookie: lms_session=token"
// @Success 204 "Session revoked"
// @Failure 401 {object} response.ErrorResponse "Invalid session"
// @Failure 403 {object} response.ErrorResponse "Untrusted request origin"
// @Failure 500 {object} response.ErrorResponse "Internal error"
// @Router /api/v1/auth/logout [post]
func (h *handler) Logout(c *fiber.Ctx) error {
	token := c.Cookies(h.sessionCookie.Name)
	if err := h.usecase.Logout(c.UserContext(), token); err != nil {
		return response.Error(c, err)
	}
	helper.SetSessionCookie(c, h.sessionCookie, "", time.Unix(1, 0))
	return c.SendStatus(http.StatusNoContent)
}

// Me returns the current session user.
//
// @Summary Get current user
// @Description Returns the user identified by the lms_session cookie.
// @Tags Authentication
// @Produce json
// @Param Cookie header string true "Session cookie: lms_session=token"
// @Success 200 {object} response.UserSuccess "Current user"
// @Failure 401 {object} response.ErrorResponse "Invalid session"
// @Failure 500 {object} response.ErrorResponse "Internal error"
// @Router /api/v1/auth/me [get]
func (h *handler) Me(c *fiber.Ctx) error {
	user, err := h.usecase.AuthenticateSession(c.UserContext(), c.Cookies(h.sessionCookie.Name))
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, response.NewUser(user))
}
