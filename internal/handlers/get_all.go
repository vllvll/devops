package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h Handler) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		gauges, counters := h.repository.GetAll()

		answer := "Gauges:\n"
		for key, value := range gauges {
			answer += fmt.Sprintf("%s - %s\n", key, strconv.FormatFloat(float64(value), 'f', 3, 64))
		}

		answer += "Counters:\n"
		for key, value := range counters {
			answer += fmt.Sprintf("%s - %s\n", key, strconv.FormatInt(int64(value), 10))
		}

		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(answer))
	}
}
