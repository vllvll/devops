package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
)

// GetMetricJSON Получение метрики в формате JSON с хешем
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
