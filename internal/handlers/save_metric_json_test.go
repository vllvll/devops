package handlers

import (
	"encoding/json"
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

func TestHandler_SaveMetricJSON(t *testing.T) {
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
				response:    "",
				contentType: "",
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
				response:    "",
				contentType: "",
			},
		},
		{
			name:      "error counter signer",
			signerKey: "6d9d04f1f54f1b11944a9bb143b4ad786d502f29f801ee75da2e612e459f98f4",
			metric: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: getCounter(100),
				Hash:  "errorhash",
			},
			want: want{
				code:        400,
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := repositories.NewStatsMemoryRepository()
			signer := services.NewMetricSigner(tt.signerKey)
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Post("/update/", handler.SaveMetricJSON())

			client := resty.New().
				SetHeader("Content-Type", "application/json")
			client.JSONMarshal = json.Marshal
			client.JSONUnmarshal = json.Unmarshal

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().SetBody(tt.metric).Post(ts.URL + "/update/")
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}

func getGauge(value float64) *float64 {
	return &value
}

func getCounter(value int64) *int64 {
	return &value
}
