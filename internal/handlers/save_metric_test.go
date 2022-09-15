package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
)

func Example_saveMetric() {
	_, _ = http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/update/%s/%s/%d", "counter", "CounterPollCount", 10),
		nil,
	)
}

func TestHandler_SaveMetric(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	tests := []struct {
		name   string
		format string
		key    string
		value  string
		want   want
	}{
		{
			name:   "format not exists",
			format: "gaage",
			key:    "Alloc",
			value:  "0.01",
			want: want{
				code:        501,
				response:    "Not Implemented",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "can't parse Gauge value",
			format: "gauge",
			key:    "Alloc",
			value:  "ffff",
			want: want{
				code:        400,
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "can't parse Counter value",
			format: "counter",
			key:    "PollCount",
			value:  "ffff",
			want: want{
				code:        400,
				response:    "Bad Request",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "success counter",
			format: "counter",
			key:    "TestCounter",
			value:  "10",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:   "success gauge",
			format: "gauge",
			key:    "TestGauge",
			value:  "10.0001",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := repositories.NewStatsMemoryRepository()
			signer := services.NewMetricSigner("")
			handler := NewHandler(repository, signer, nil)

			r := chi.NewRouter()
			r.Post("/update/{format:[A-Za-z]+}/{key:[A-Za-z0-9]+}/{value:[A-Za-z0-9.]+}", handler.SaveMetric())

			ts := httptest.NewServer(r)
			defer ts.Close()

			request, err := http.NewRequest(
				http.MethodPost,
				ts.URL+fmt.Sprintf("/update/%s/%s/%s", tt.format, tt.key, tt.value),
				nil,
			)
			require.NoError(t, err)

			response, err := http.DefaultClient.Do(request)
			require.NoError(t, err)

			responseBody, err := ioutil.ReadAll(response.Body)
			require.NoError(t, err)
			defer response.Body.Close()

			assert.Equal(t, tt.want.code, response.StatusCode)
			assert.Equal(t, tt.want.response, strings.Trim(string(responseBody), "\n"))
			assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
		})
	}
}
