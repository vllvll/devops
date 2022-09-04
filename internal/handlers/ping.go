package handlers

import "net/http"

func (h Handler) Ping() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := h.db.Ping()
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(http.StatusText(http.StatusOK)))
	}
}
