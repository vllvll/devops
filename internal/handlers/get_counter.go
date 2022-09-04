package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetCounter Получение значения типа Counter по ключу
func (h Handler) GetCounter() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		if key == "" {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		value, err := h.repository.GetCounterByKey(key)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(strconv.FormatInt(int64(value), 10)))
	}
}
