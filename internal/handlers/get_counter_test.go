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

func Example_getCounter() {
	client := resty.New()
	_, _ = client.R().Get("/value/counter/PollCount")
}

func TestHandler_GetCounter(t *testing.T) {
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
				response:    "100",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:      "without counter key",
			signerKey: "",
			metric: types.Metrics{
				Delta: getCounter(100),
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
			repository.UpdateCount(tt.metric.ID, types.Counter(*tt.metric.Delta))

			signer := services.NewMetricSigner(tt.signerKey)
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Get("/value/counter/{key:[A-Za-z0-9]+}", handler.GetCounter())

			client := resty.New()

			ts := httptest.NewServer(r)
			defer ts.Close()

			response, err := client.R().Get(ts.URL + "/value/counter/" + tt.metric.ID)
			require.NoError(t, err)

			assert.Equal(t, tt.want.code, response.StatusCode())
			assert.Equal(t, tt.want.response, strings.Trim(string(response.Body()), "\n"))
			assert.Equal(t, tt.want.contentType, response.Header().Get("Content-Type"))
		})
	}
}
