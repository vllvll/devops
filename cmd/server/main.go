package main

import (
	"context"
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/handlers"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/routes"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/storage"
	"github.com/vllvll/devops/internal/storage/file"
	"github.com/vllvll/devops/pkg/postgres"
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
		log.Fatalf("Error with config: %v", err)
	}

	db, err := postgres.ConnectDatabase(config.DatabaseDsn)
	if err != nil {
		log.Fatalf("Error with database: %v", err)
	}
	defer db.Close()

	var storeTick = time.Tick(config.StoreInterval)

	statsRepository := repositories.NewStatsDatabaseRepository(db)
	if config.DatabaseDsn == "" {
		statsRepository = repositories.NewStatsMemoryRepository()
	}

	signer := services.NewMetricSigner(config.Key)
	handler := handlers.NewHandler(statsRepository, signer, db)
	router := routes.NewRouter(*handler)
	router.RegisterHandlers()

	consumer, err := file.NewFileConsumer(config.StoreFile)
	if err != nil {
		log.Fatalf("Error with file consumer: %v", err)
	}

	producer, err := file.NewFileProducer(config.StoreFile)
	if err != nil {
		log.Fatalf("Error with file producer: %v", err)
	}
	defer producer.Close()

	fileStorage := storage.NewStatsStorage(config, consumer, producer)

	defer fileStorage.Save(statsRepository)

	statsRepository, err = fileStorage.Start(statsRepository)
	if err != nil {
		log.Fatalf("Error with file file storage: %v", err)
	}

	httpServer := &http.Server{
		Addr:    config.Address,
		Handler: router.Router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error with HTTP server ListenAndServe: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			log.Println("Graceful shutdown")

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
