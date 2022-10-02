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

func Example_bulkSaveMetricJSON() {
	metrics := make([]types.Metrics, 0)
	metrics = append(metrics, types.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Value: getGauge(0.1),
		Hash:  "",
	})

	body, _ := json.Marshal(metrics)

	client := resty.New().
		SetHeader("Content-Type", "application/json")
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	_, _ = client.R().SetBody(body).Post("/updates/")
}

func TestHandler_BulkSaveMetricJSON(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name           string
		signerKey      string
		privateKeyPath string
		metric         types.Metrics
		want           want
	}{
		{
			name: "gauge success",
			metric: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
				Value: getGauge(0.1),
				Hash:  "",
			},
			want: want{
				code:        200,
				response:    "OK",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "counter success",
			metric: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
				Delta: getCounter(100),
				Hash:  "",
			},
			want: want{
				code:        200,
				response:    "OK",
				contentType: "text/plain; charset=utf-8",
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
			decrypt, _ := services.NewMetricDecrypt(tt.privateKeyPath)
			handler := NewHandler(repository, signer, nil, decrypt)

			r := chi.NewRouter()
			r.Post("/updates/", handler.BulkSaveMetricJSON())

			client := resty.New().
				SetHeader("Content-Type", "application/json")
			client.JSONMarshal = json.Marshal
			client.JSONUnmarshal = json.Unmarshal

			ts := httptest.NewServer(r)
			defer ts.Close()

			metrics := make([]types.Metrics, 0)
			metrics = append(metrics, tt.metric)

			response, err := client.R().SetBody(metrics).Post(ts.URL + "/updates/")
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}
