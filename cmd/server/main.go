package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	conf "github.com/vllvll/devops/internal/config"
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
	config, err := conf.CreateServerConfig()
	if err != nil {
		panic("Конфиг не загружен")
	}

	metricRepository := metric.NewRepository()
	metricHandler := metric.NewHandler(metricRepository)

	start(config, metricRepository)

	var storeTick = time.Tick(config.StoreInterval)

	r := routerChi.CreateRouter()

	r.Get("/", metricHandler.GetAll())
	r.Route("/value/", func(r chi.Router) {
		r.Post("/", metricHandler.GetMetricJSON())
		r.Get("/gauge/{key:[A-Za-z0-9]+}", metricHandler.GetGauge())
		r.Get("/counter/{key:[A-Za-z0-9]+}", metricHandler.GetCounter())
	})
	r.Post("/update/{format:[A-Za-z]+}/{key:[A-Za-z0-9]+}/{value:[A-Za-z0-9.]+}", metricHandler.SaveMetric())
	r.Post("/update/", metricHandler.SaveMetricJSON())

	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			// graceful shutdown
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if err := httpServer.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			cancel()

			save(config, metricRepository)

			return
		case <-storeTick:
			save(config, metricRepository)
		}
	}
}

func save(config *conf.ServerConfig, repository metric.RepositoryInterface) {
	var metrics []metric.Metrics

	fsProducer, err := metric.NewProducer(config.StoreFile)
	if err != nil {
		panic("Filesystem producer не загружен")
	}

	gauges, counters := repository.GetAll()

	for key, value := range gauges {
		flValue := float64(value)

		metrics = append(metrics, metric.Metrics{
			ID:    key,
			MType: metric.GaugeType,
			Value: &flValue,
		})
	}

	for key, value := range counters {
		iValue := int64(value)

		metrics = append(metrics, metric.Metrics{
			ID:    key,
			MType: metric.CounterType,
			Delta: &iValue,
		})
	}

	for _, m := range metrics {
		err := fsProducer.WriteMetric(&m)
		if err != nil {
			panic("can't write metric")
		}
	}

	fsProducer.Close()
}

func start(config *conf.ServerConfig, metricRepository metric.RepositoryInterface) {
	if config.Restore {
		fsConsumer, err := metric.NewConsumer(config.StoreFile)
		if err != nil {
			panic("Filesystem consumer не загружен")
		}
		defer fsConsumer.Close()

		for {
			readMetric, err := fsConsumer.ReadMetric()
			if err != nil {
				return
			}

			switch readMetric.MType {
			case metric.GaugeType:
				metricRepository.UpdateGauge(readMetric.ID, metric.Gauge(*readMetric.Value))

			case metric.CounterType:
				metricRepository.UpdateCount(readMetric.ID, metric.Counter(*readMetric.Delta))
			}
		}
	}
}
