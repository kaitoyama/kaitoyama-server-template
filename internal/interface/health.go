package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) HealthCheck(c echo.Context) error {
	healthStatus, err := h.healthUsecase.CheckHealth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, healthStatus)
	}
	return c.JSON(http.StatusOK, healthStatus)
}
