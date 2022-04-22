package metric

import (
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

		if len(path) != 4 {
			http.Error(rw, "Not found", http.StatusNotFound)

			return
		}

		entrypoint := path[0]
		format := path[1]
		metricName := path[2]
		value := path[3]

		if entrypoint != "update" {
			http.Error(rw, "Not found", http.StatusNotFound)
		}

		switch format {
		case "gauge":
			f, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(rw, "Value of metric incorrect", http.StatusBadRequest)
			}

			h.repository.UpdateMetric(metricName, Gauge(f))

		case "counter":
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				http.Error(rw, "Value of metric incorrect", http.StatusBadRequest)
			}

			h.repository.UpdateCount(metricName, Counter(i))

		default:
			http.Error(rw, "Not implemented", http.StatusNotImplemented)

			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
