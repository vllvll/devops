// Package handlers содержит обработчики проекта
package handlers

import (
	"database/sql"
	"net/http"

	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
)

type Handler struct {
	repository repositories.StatsRepository // Сервис для чтения и записи данных метрик
	signer     services.Signer              // Сервис для создания подписи
	db         *sql.DB                      // База данных
}

// NewHandler Получение хендлера
func NewHandler(repository repositories.StatsRepository, signer services.Signer, db *sql.DB) *Handler {
	return &Handler{
		repository: repository,
		signer:     signer,
		db:         db,
	}
}

// MetricHandlers Список методов для хендлеров (сервер)
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
