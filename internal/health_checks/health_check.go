package health_checks

import (
	"context"
	"database/sql"

	"awesomeProjectCr/internal/config"
)

type HealthStatus struct {
	Status      string      `json:"status"`
	Version     string      `json:"version"`
	DBStatus    sql.DBStats `json:"db_status"`
	CacheStatus string      `json:"redis_connection_status"`
}

type Service struct {
	DB *sql.DB
}

type HealthCheck interface {
	GetStatus(ctx context.Context) (HealthStatus, error)
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (h *Service) GetStatus(ctx context.Context) (HealthStatus, error) {
	var status HealthStatus
	var err error

	status.Status = "OK"
	status.CacheStatus = "connected"

	if err = h.DB.Ping(); err != nil {
		status.Status = "Errror"
	}

	if err != nil {
		status.Status = "Error"
		status.CacheStatus = "disconnected"
	}

	status.Version = config.App.Version
	status.DBStatus = h.DB.Stats()
	return status, err
}
