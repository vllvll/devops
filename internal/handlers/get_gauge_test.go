package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/types"
)

func Example_getGauge() {
	client := resty.New()
	_, _ = client.R().Get("/value/gauge/Alloc")
}

func TestHandler_GetGauge(t *testing.T) {
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
				response:    "0.100",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "without gauge key",
			signerKey: "",
			metric: types.Metrics{
				Value: getGauge(0),
			},
			want: want{
				code:        404,
				response:    "404 page not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := repositories.NewStatsMemoryRepository()
			repository.UpdateGauge(tt.metric.ID, types.Gauge(*tt.metric.Value))

			signer := services.NewMetricSigner(tt.signerKey)
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Get("/value/gauge/{key:[A-Za-z0-9]+}", handler.GetGauge())

			client := resty.New()

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().Get(ts.URL + "/value/gauge/" + tt.metric.ID)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}
