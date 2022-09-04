package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h Handler) GetGauge() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		value, err := h.repository.GetGaugeByKey(key)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(strconv.FormatFloat(float64(value), 'f', 3, 64)))
	}
}
