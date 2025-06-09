package counter

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"module6/internal/services"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

type CounterResponse struct {
	Count int `json:"count"`
}

// @Summary IncrementCounter
// @Description Increments counter by 1
// @Tags counter
// @Success 200 {object} CounterResponse
// @Router /api/counter/increment [get]
func (h *Handler) IncrementCounter(ctx *fiber.Ctx) error {
	count, err := h.svc.Counter.IncrementOne(ctx.Context())
	if err != nil {
		return fmt.Errorf("increment counter: %w", err)
	}
	return ctx.JSON(&CounterResponse{Count: count})
}
