package handlers

import (
	"database/sql"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"net/http"
)

type Handler struct {
	repository repositories.StatsRepository
	signer     services.Signer
	db         *sql.DB
}

func NewHandler(repository repositories.StatsRepository, signer services.Signer, db *sql.DB) *Handler {
	return &Handler{
		repository: repository,
		signer:     signer,
		db:         db,
	}
}

type MetricHandlers interface {
	SaveMetricJSON() http.HandlerFunc
	SaveMetric() http.HandlerFunc
	GetAll() http.HandlerFunc
	GetMetricJSON() http.HandlerFunc
	GetGauge() http.HandlerFunc
	GetCounter() http.HandlerFunc
	Ping() http.HandlerFunc
	BulkSaveMetricJSON() http.HandlerFunc
}
