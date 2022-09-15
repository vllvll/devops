package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vllvll/devops/internal/dictionaries"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/types"
)

func Example_getAll() {
	client := resty.New()
	_, _ = client.R().Get("/")
}

func TestHandler_GetAll(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name      string
		signerKey string
		metric    types.Metrics
		want      want
	}{
		{
			name:      "gauge success",
			signerKey: "",
			metric: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: getGauge(0.1),
				Hash:  "",
			},
			want: want{
				code:        200,
				response:    "Gauges:\nAlloc - 0.100\nCounters:",
				contentType: "text/html",
			},
		},
		{
			name:      "counter success",
			signerKey: "",
			metric: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: getCounter(100),
				Hash:  "",
			},
			want: want{
				code:        200,
				response:    "Gauges:\nCounters:\nPollCount - 100",
				contentType: "text/html",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := repositories.NewStatsMemoryRepository()
			switch tt.metric.MType {
			case dictionaries.CounterType:
				repository.UpdateCount(tt.metric.ID, types.Counter(*tt.metric.Delta))
			case dictionaries.GaugeType:
				repository.UpdateGauge(tt.metric.ID, types.Gauge(*tt.metric.Value))
			}

			signer := services.NewMetricSigner(tt.signerKey)
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Get("/", handler.GetAll())

			client := resty.New()

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().Get(ts.URL)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}
