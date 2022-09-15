package handlers

import (
	"encoding/json"
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

func Example_getMetricJSON() {
	client := resty.New().
		SetHeader("Content-Type", "application/json")
	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	_, _ = client.R().SetBody(types.Metrics{
		ID:    "Alloc",
		MType: "gauge",
	}).Post("/value/")
}

func TestHandler_GetMetricJSON(t *testing.T) {
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
				response:    "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":0.1,\"hash\":\"d030b4a233bcff5054e6381220698c3a0ff6ceaa067d50bf886658db57cdf983\"}",
				contentType: "application/json",
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
				response:    "{\"id\":\"PollCount\",\"type\":\"counter\",\"delta\":100,\"hash\":\"e6bb86b15f8b233556b65d3ae04e0e79581ae6fb6a8370262a8fe8f0f5dba78f\"}",
				contentType: "application/json",
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
			r.Post("/value/", handler.GetMetricJSON())

			client := resty.New().
				SetHeader("Content-Type", "application/json")
			client.JSONMarshal = json.Marshal
			client.JSONUnmarshal = json.Unmarshal

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().SetBody(types.Metrics{
				ID:    tt.metric.ID,
				MType: tt.metric.MType,
			}).Post(ts.URL + "/value/")
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}

	failTests := []struct {
		name      string
		signerKey string
		metric    types.Metrics
		want      want
	}{
		{
			name:      "gauge not found",
			signerKey: "",
			metric: types.Metrics{
				ID:    "Alloc",
				MType: "gauge",
			},
			want: want{
				code:        404,
				response:    "Not Found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "counter not found",
			signerKey: "",
			metric: types.Metrics{
				ID:    "PollCount",
				MType: "counter",
			},
			want: want{
				code:        404,
				response:    "Not Found",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range failTests {
		t.Run(tt.name, func(t *testing.T) {
			repository := repositories.NewStatsMemoryRepository()
			signer := services.NewMetricSigner("")
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Post("/value/", handler.GetMetricJSON())

			client := resty.New().
				SetHeader("Content-Type", "application/json")
			client.JSONMarshal = json.Marshal
			client.JSONUnmarshal = json.Unmarshal

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().SetBody(types.Metrics{
				ID:    tt.metric.ID,
				MType: tt.metric.MType,
			}).Post(ts.URL + "/value/")
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}
