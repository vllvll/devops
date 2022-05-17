package services

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vllvll/devops/internal/handlers"
)

type Router struct {
	Router   chi.Router
	handlers handlers.Handler
}

func NewRouter(handlers handlers.Handler) Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	return Router{
		Router:   r,
		handlers: handlers,
	}
}

func (ro *Router) RegisterHandlers() {
	ro.Router.Get("/", ro.handlers.GetAll())
	ro.Router.Route("/value/", func(r chi.Router) {
		r.Post("/", ro.handlers.GetMetricJSON())
		r.Get("/gauge/{key:[A-Za-z0-9]+}", ro.handlers.GetGauge())
		r.Get("/counter/{key:[A-Za-z0-9]+}", ro.handlers.GetCounter())
	})
	ro.Router.Post("/update/{format:[A-Za-z]+}/{key:[A-Za-z0-9]+}/{value:[A-Za-z0-9.]+}", ro.handlers.SaveMetric())
	ro.Router.Post("/update/", ro.handlers.SaveMetricJSON())
}
