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
	var reportMain = time.Tick(config.ReportInterval)

	var pollCount types.Counter

	signer := services.NewMetricSigner(config.Key)
	sender := services.NewSendClient(config, signer)
	constants := dictionaries.NewMemConstants()
	memRepository := repositories.NewMemRepository(constants)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	errCh := make(chan error)

	gaugesCh := make(chan types.Gauges)
	counterCh := make(chan types.Counters)

	metricCh := make(chan types.Metrics)
	go sender.Prepare(gaugesCh, counterCh, metricCh, errCh)
	go sender.Send(metricCh, reportTick, errCh)

	for {
		select {
		case <-c:
			log.Println("Graceful shutdown")

			close(gaugesCh)
			close(counterCh)

			return

		case <-pollTick:
			go memRepository.GetGauges(gaugesCh, errCh)
			go memRepository.GetAdditionalGauges(gaugesCh, errCh)

			pollCount++

		case <-reportMain:
			counterCh <- types.Counters{dictionaries.CounterPollCount: pollCount}

			pollCount = 0

		case <-errCh:
			log.Printf("Error: %v\n", err)
		}
	}
}
