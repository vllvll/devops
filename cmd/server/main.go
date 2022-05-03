package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/vllvll/devops/internal/metric"
	routerChi "github.com/vllvll/devops/pkg/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	metricRepository := metric.NewRepository()
	metricConstants := metric.NewConstants()
	metricHandler := metric.NewHandler(metricRepository, metricConstants)

	r := routerChi.CreateRouter()

	r.Get("/", metricHandler.GetAll())
	r.Route("/value/", func(r chi.Router) {
		r.Post("/", metricHandler.GetMetric())
		r.Get("/gauge/{key:[A-Za-z0-9]+}", metricHandler.GetGauge())
		r.Get("/counter/{key:[A-Za-z0-9]+}", metricHandler.GetCounter())
	})
	r.Post("/update/{format:[A-Za-z]+}/{key:[A-Za-z0-9]+}/{value:[A-Za-z0-9.]+}", metricHandler.SaveMetric())
	r.Post("/update/", metricHandler.SaveMetricJson())

	httpServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
