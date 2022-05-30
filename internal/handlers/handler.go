package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/types"
	"net/http"
	"strconv"
)

type Handler struct {
	repository repositories.StatsRepository
	signer     services.Signer
	db         *sql.DB
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

func NewHandler(repository repositories.StatsRepository, signer services.Signer, db *sql.DB) *Handler {
	return &Handler{
		repository: repository,
		signer:     signer,
		db:         db,
	}
}

func (h Handler) SaveMetricJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var metric types.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		switch metric.MType {
		case dictionaries.GaugeType:
			if !h.signer.IsEqualHashGauge(metric.ID, *metric.Value, metric.Hash) {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateGauge(metric.ID, types.Gauge(*metric.Value))

		case dictionaries.CounterType:
			if !h.signer.IsEqualHashCounter(metric.ID, *metric.Delta, metric.Hash) {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateCount(metric.ID, types.Counter(*metric.Delta))
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func (h Handler) SaveMetric() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		format := chi.URLParam(r, "format")
		key := chi.URLParam(r, "key")
		value := chi.URLParam(r, "value")

		switch format {
		case dictionaries.GaugeType:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateGauge(key, types.Gauge(f))

		case dictionaries.CounterType:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateCount(key, types.Counter(i))

		default:
			http.Error(rw, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)

			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}

func (h Handler) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		gauges, counters := h.repository.GetAll()

		answer := "Gauges:\n"
		for key, value := range gauges {
			answer += fmt.Sprintf("%s - %s\n", key, strconv.FormatFloat(float64(value), 'f', 3, 64))
		}

		answer += "Counters:\n"
		for key, value := range counters {
			answer += fmt.Sprintf("%s - %s\n", key, strconv.FormatInt(int64(value), 10))
		}

		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(answer))
	}
}

func (h Handler) GetMetricJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var metric types.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		switch metric.MType {
		case dictionaries.GaugeType:
			var value float64

			gauge, err := h.repository.GetGaugeByKey(metric.ID)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

				return
			}

			value = float64(gauge)
			metric.Value = &value
			metric.Hash = h.signer.GetHashGauge(metric.ID, value)

		case dictionaries.CounterType:
			var value int64

			counter, err := h.repository.GetCounterByKey(metric.ID)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

				return
			}

			value = int64(counter)
			metric.Delta = &value
			metric.Hash = h.signer.GetHashCounter(metric.ID, value)
		}

		response, err := json.Marshal(metric)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(response)
	}
}

func (h Handler) GetGauge() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		value, err := h.repository.GetGaugeByKey(key)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(strconv.FormatFloat(float64(value), 'f', 3, 64)))
	}
}

func (h Handler) GetCounter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		value, err := h.repository.GetCounterByKey(key)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(strconv.FormatInt(int64(value), 10)))
	}
}

func (h Handler) Ping() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.db.Ping()
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}

func (h Handler) BulkSaveMetricJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var metrics []types.Metrics
		var counters = types.Counters{}
		var gauges = types.Gauges{}

		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		for _, metric := range metrics {
			switch metric.MType {
			case dictionaries.GaugeType:
				if !h.signer.IsEqualHashGauge(metric.ID, *metric.Value, metric.Hash) {
					http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

					return
				}

				gauges[metric.ID] = types.Gauge(*metric.Value)
			case dictionaries.CounterType:
				if !h.signer.IsEqualHashCounter(metric.ID, *metric.Delta, metric.Hash) {
					http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

					return
				}

				counters[metric.ID] += types.Counter(*metric.Delta)
			}
		}

		err := h.repository.UpdateAll(gauges, counters)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}
