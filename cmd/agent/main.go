package main

import (
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/types"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := conf.CreateAgentConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	var pollTick = time.Tick(config.PollInterval)
	var reportTick = time.Tick(config.ReportInterval)

	var pollCount types.Counter
	var gauges = types.Gauges{}

	signer := services.NewMetricSigner(config.Key)
	sender := services.NewSendClient(config, signer)
	constants := dictionaries.NewMemConstants()
	memRepository := repositories.NewMemRepository(constants)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			log.Println("Graceful shutdown")

			return
		case <-pollTick:
			pollCount++

			gauges = memRepository.GetGauges()

		case <-reportTick:
			err := sender.Send(gauges, pollCount)
			if err != nil {
				log.Printf("Error with send report: %v\n", err)
			}

			pollCount = 0
		}
	}
}
