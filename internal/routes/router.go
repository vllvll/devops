// Package routes содержит функционал для работы роутингом
package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/vllvll/devops/internal/handlers"
	"github.com/vllvll/devops/internal/middlewares"
)

type Router struct {
	Router   chi.Router       // Роутер
	handlers handlers.Handler // Обработчики
}

// NewRouter Регистрируем middleware и возвращаем роутер
func NewRouter(handlers handlers.Handler, trustedSubnet string) Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.TrustedSubnet(trustedSubnet))

	// r.Mount("/debug", middleware.Profiler())

	return Router{
		Router:   r,
		handlers: handlers,
	}
}

// RegisterHandlers Регистрируем обработчики
func (ro *Router) RegisterHandlers() {
	ro.Router.Get("/", ro.handlers.GetAll())
	ro.Router.Get("/ping", ro.handlers.Ping())
	ro.Router.Route("/value/", func(r chi.Router) {
		r.Post("/", ro.handlers.GetMetricJSON())
		r.Get("/gauge/{key:[A-Za-z0-9]+}", ro.handlers.GetGauge())
		r.Get("/counter/{key:[A-Za-z0-9]+}", ro.handlers.GetCounter())
	})
	ro.Router.Post("/update/{format:[A-Za-z]+}/{key:[A-Za-z0-9]+}/{value:[A-Za-z0-9.]+}", ro.handlers.SaveMetric())
	ro.Router.Post("/update/", ro.handlers.SaveMetricJSON())
	ro.Router.Post("/updates/", ro.handlers.BulkSaveMetricJSON())
}
