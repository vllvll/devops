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
)

func Example_ping() {
	client := resty.New()
	_, _ = client.R().Get("/ping/")
}

func TestHandler_Ping(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		repository := repositories.NewStatsMemoryRepository()
		signer := services.NewMetricSigner("")
		handler := NewHandler(repository, signer, nil)

		r := chi.NewRouter()
		r.Get("/ping", handler.Ping())

		client := resty.New()

		ts := httptest.NewServer(r)
		defer ts.Close()

		response, err := client.R().Get(ts.URL + "/ping")
		require.NoError(t, err)

		assert.Equal(t, 500, response.StatusCode())
		assert.Equal(t, "Internal Server Error", strings.Trim(string(response.Body()), "\n"))
		assert.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
	})
}
