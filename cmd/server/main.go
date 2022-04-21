package main

import (
	"context"
	"fmt"
	"github.com/vllvll/devops/internal/metric"
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

	httpServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: metricHandler.SaveMetric(),
	}

	http.HandleFunc("/update/", metricHandler.SaveMetric())

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case <-c:
		fmt.Println("Graceful shutdown")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
