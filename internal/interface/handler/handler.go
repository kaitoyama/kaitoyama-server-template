package handler

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
	"github.com/kaitoyama/kaitoyama-server-template/internal/usecase"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	healthUsecase usecase.HealthUsecase
}

func (h *Handler) GetHealth(c echo.Context) error {
	return h.HealthCheck(c)
}

func NewHandler(dbChecker domain.DatabaseHealthChecker) *Handler {
	return &Handler{
		healthUsecase: *usecase.NewHealthUsecase(dbChecker),
	}
}
