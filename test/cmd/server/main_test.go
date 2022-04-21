package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vllvll/devops/internal/metric"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSaveMetricHandler(t *testing.T) {
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
				code:        404,
				response:    "Not found",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:   "key not exists",
			format: "gauge",
			key:    "Alloce",
			value:  "0.01",
			want: want{
				code:        400,
				response:    "Name of metric incorrect",
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
				response:    "Value of metric incorrect",
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
				response:    "Value of metric incorrect",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/update/%s/%s/%s", tt.format, tt.key, tt.value),
				nil,
			)

			constants := metric.NewConstants()
			repository := metric.NewRepository()
			handler := metric.NewHandler(repository, constants)

			w := httptest.NewRecorder()
			h := handler.SaveMetric()
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, w.Code)

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.want.response, strings.Trim(string(resBody), "\n"))
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
		})
	}
}
