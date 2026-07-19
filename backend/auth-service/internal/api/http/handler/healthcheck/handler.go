package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	healthusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/healthcheck"
)

type Handler struct {
	usecase healthusecase.Usecase
}

func NewHandler(usecase healthusecase.Usecase) *Handler {
	return &Handler{usecase: usecase}
}

// Liveness reports whether the process can serve requests.
//
// @Summary Check liveness
// @Tags Health
// @Produce json
// @Success 200 {object} response.Health "Process is live"
// @Router /health/liveness [get]
func (h *Handler) Liveness(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(response.Health{Status: "ok"})
}

// Readiness reports whether required dependencies are available.
//
// @Summary Check readiness
// @Tags Health
// @Produce json
// @Success 200 {object} response.Health "Service is ready"
// @Failure 503 {object} response.Health "A required dependency is unavailable"
// @Router /health/readiness [get]
func (h *Handler) Readiness(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()
	if err := h.usecase.Readiness(ctx); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(response.Health{Status: "unavailable"})
	}
	return c.Status(http.StatusOK).JSON(response.Health{Status: "ready"})
}
