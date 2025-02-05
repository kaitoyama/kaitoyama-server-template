package usecase

import (
	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

type HealthUsecase struct {
	dbHealthChecker domain.DatabaseHealthChecker
}

func NewHealthUsecase(dbChecker domain.DatabaseHealthChecker) *HealthUsecase {
	return &HealthUsecase{
		dbHealthChecker: dbChecker,
	}
}

func (h *HealthUsecase) CheckHealth() (*domain.HealthStatus, error) {
	if err := h.dbHealthChecker.PingDB(); err != nil {
		return &domain.HealthStatus{
			Status: "error",
			Error:  "Database connection failed",
		}, err
	}

	return &domain.HealthStatus{
		Status: "ok",
		Error:  "",
	}, nil
}
