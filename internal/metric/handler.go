package metric

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	repository *Repository
	constants  *Сonstants
}

func NewHandler(repository *Repository, constants *Сonstants) *Handler {
	return &Handler{
		repository: repository,
		constants:  constants,
	}
}

func (h Handler) SaveMetric() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		if path[1] != "gauge" && path[1] != "counter" {
			fmt.Println("not found", path)
			http.Error(rw, "Not found", http.StatusNotFound)

			return
		}

		if !h.constants.In(path[2]) {
			http.Error(rw, "Name of metric incorrect", http.StatusBadRequest)
		}

		metricName := fmt.Sprintf("%s", path[2])

		switch path[1] {
		case "gauge":
			f, err := strconv.ParseFloat(path[3], 64)
			if err != nil {
				http.Error(rw, "Value of metric incorrect", http.StatusBadRequest)
			}

			h.repository.UpdateMetric(metricName, Gauge(f))

		case "counter":
			i, err := strconv.ParseInt(path[3], 10, 64)
			if err != nil {
				http.Error(rw, "Value of metric incorrect", http.StatusBadRequest)
			}

			h.repository.UpdateCount(Counter(i))
		}

		rw.WriteHeader(http.StatusOK)
	}
}
