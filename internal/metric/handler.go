package metric

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Handler struct {
	repository RepositoryInterface
}

func NewHandler(repository RepositoryInterface) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h Handler) SaveMetricJSON() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var metric Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		switch metric.MType {
		case GaugeType:
			h.repository.UpdateMetric(metric.MType, Gauge(*metric.Value))

		case CounterType:
			h.repository.UpdateCount(metric.MType, Counter(*metric.Delta))
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
		case GaugeType:
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateMetric(key, Gauge(f))

		case CounterType:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateCount(key, Counter(i))

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
		var metric Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

			return
		}

		switch metric.MType {
		case GaugeType:
			var value float64

			gauge, err := h.repository.GetGaugeByKey(metric.MType)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

				return
			}

			value = float64(gauge)
			metric.Value = &value

		case CounterType:
			var value int64

			counter, err := h.repository.GetCounterByKey(metric.MType)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

				return
			}

			value = int64(counter)
			metric.Delta = &value
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
