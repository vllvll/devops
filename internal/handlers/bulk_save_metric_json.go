package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
)

// BulkSaveMetricJSON Сохранение всех описанных метрик (Gauge, Counter) в запросе за один раз
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
