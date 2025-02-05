package db

import (
	"context"
	"database/sql"

	"github.com/kaitoyama/kaitoyama-server-template/internal/domain"
)

type dbHealthChecker struct {
	db *sql.DB
}

func NewDBHealthChecker(db *sql.DB) domain.DatabaseHealthChecker {
	return &dbHealthChecker{
		db: db,
	}
}

func (d *dbHealthChecker) Check() (*domain.HealthStatus, error) {
	if err := d.PingDB(); err != nil {
		return &domain.HealthStatus{
			Status: "ok",
		}, nil
	}
	return &domain.HealthStatus{
		Status: "ok",
	}, nil
}

func (d *dbHealthChecker) PingDB() error {
	return d.db.PingContext(context.Background())
}
