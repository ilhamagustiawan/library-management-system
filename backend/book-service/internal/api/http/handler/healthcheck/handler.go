package healthcheck

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/healthcheck"
)

type Handler struct{ usecase *healthcheck.Usecase }

func NewHandler(usecase *healthcheck.Usecase) *Handler { return &Handler{usecase: usecase} }

// Liveness reports process health.
// @Summary Liveness
// @Tags Health
// @Produce json
// @Success 200 {object} response.Health
// @Router /health/liveness [get]
func (h *Handler) Liveness(ctx *fiber.Ctx) error {
	if err := h.usecase.Liveness(ctx.UserContext()); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Status(http.StatusOK).JSON(response.Health{Status: "ok"})
}

// Readiness reports database readiness.
// @Summary Readiness
// @Tags Health
// @Produce json
// @Success 200 {object} response.Health
// @Failure 500 {object} response.ErrorResponse
// @Router /health/readiness [get]
func (h *Handler) Readiness(ctx *fiber.Ctx) error {
	if err := h.usecase.Readiness(ctx.UserContext()); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.Status(http.StatusOK).JSON(response.Health{Status: "ok"})
}
