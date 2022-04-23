package metric

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type Handler struct {
	repository *Repository
	constants  *Сonstants
}

func NewHandler(repository *Repository, constants *Сonstants) *Handler {
	return &Handler{
		repository: repository,
		constants:  constants,
	}
}

func (h Handler) SaveMetric() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		format := chi.URLParam(r, "format")
		key := chi.URLParam(r, "key")
		value := chi.URLParam(r, "value")

		switch format {
		case "gauge":
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

				return
			}

			h.repository.UpdateMetric(key, Gauge(f))

		case "counter":
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
