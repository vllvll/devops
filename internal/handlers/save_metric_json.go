package handlers

import (
	"encoding/json"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	"net/http"
)

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
