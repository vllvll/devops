package main

import (
	"context"
	"fmt"
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/handlers"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/routes"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/storage"
	"github.com/vllvll/devops/internal/storage/file"
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

	var storeTick = time.Tick(config.StoreInterval)

	statsRepository := repositories.NewStatsRepository()
	signer := services.NewMetricSigner(config.Key)
	handler := handlers.NewHandler(statsRepository, signer)
	router := routes.NewRouter(*handler)
	router.RegisterHandlers()

	consumer, err := file.NewFileConsumer(config.StoreFile)
	if err != nil {
		panic("Консьюмер не загружен")
	}

	producer, err := file.NewFileProducer(config.StoreFile)
	if err != nil {
		panic("Продюсер не загружен")
	}

	fileStorage := storage.NewStatsStorage(config, consumer, producer)
	defer fileStorage.Save(statsRepository)

	statsRepository, err = fileStorage.Start(statsRepository)
	if err != nil {
		panic("Загрузка из файла произошла с ошибкой")
	}

	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: router.Router,
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
			fmt.Println("Graceful shutdown")

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if err := httpServer.Shutdown(ctx); err != nil {
				log.Println(err)
			}

			cancel()

			fileStorage.Save(statsRepository)

			return
		case <-storeTick:
			fileStorage.Save(statsRepository)
		}
	}
}
