package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	"net/http"
	"strconv"
)

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
