package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/request"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/response"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

type Handler struct {
	usecase  registrationusecase.Usecase
	validate *validator.Validate
}

func NewHandler(usecase registrationusecase.Usecase, validate *validator.Validate) *Handler {
	return &Handler{usecase: usecase, validate: validate}
}

// Register creates a member profile and Auth identity.
//
// @Summary Register member
// @Description Creates a member profile. Role assignment is not accepted.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body request.Register true "Registration details"
// @Success 201 {object} response.SuccessResponse "Member created"
// @Failure 409 {object} response.ErrorResponse "Email already registered"
// @Failure 422 {object} response.ErrorResponse "Invalid registration data"
// @Failure 429 {object} response.ErrorResponse "Rate limit exceeded"
// @Failure 503 {object} response.ErrorResponse "Registration dependency unavailable"
// @Router /api/v1/users [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var input request.Register
	if err := request.DecodeStrictJSON(c, &input); err != nil || h.validate.Struct(input) != nil {
		return response.ValidationError(c)
	}
	user, err := h.usecase.Register(c.UserContext(), registrationusecase.Input{
		Name: input.Name, Email: input.Email, Password: input.Password,
	})
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusCreated, response.NewUser(user))
}
