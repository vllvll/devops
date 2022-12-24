// Модуль agent отправляет информацию о состоянии
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/services"
	"github.com/vllvll/devops/internal/types"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

const BuildTemplate = `
Build version: {{ .version }}
Build date: {{ .date }}
Build commit: {{ .commit }}
`

func main() {
	t := template.Must(template.New("build").Parse(BuildTemplate))
	err := t.Execute(os.Stdout, map[string]string{
		"version": buildVersion,
		"date":    buildDate,
		"commit":  buildCommit,
	})
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	config, err := conf.CreateAgentConfig()
	if err != nil {
		log.Fatalf("Error with config: %v", err)
	}

	var pollCount types.Counter

	crypt, err := services.NewMetricEncrypt(config.CryptoKey)
	if err != nil {
		log.Fatalf("Ошибка с инициализацией сервиса шифрования: %v", err)
	}

	signer := services.NewMetricSigner(config.Key)
	sender, err := services.NewGRPCSendClient(config, signer, crypt)
	if err != nil {
		log.Fatalf("Ошибка с иницализацией сервиса http клиента: %v", err)
	}
	constants := dictionaries.NewMemConstants()
	memRepository := repositories.NewMemRepository(constants)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	errCh := make(chan error)
	gaugesCh := make(chan types.Gauges)
	counterCh := make(chan types.Counters)
	metricCh := make(chan types.Metrics)

	var pollTick = time.Tick(config.PollInterval)
	var reportTick = time.Tick(config.ReportInterval)
	var reportMain = time.Tick(config.ReportInterval)

	go sender.Prepare(ctx, gaugesCh, counterCh, metricCh, errCh)
	go sender.Send(ctx, metricCh, reportTick, errCh)

	for {
		select {
		case <-ctx.Done():
			close(gaugesCh)
			close(counterCh)
			close(errCh)

			return

		case <-pollTick:
			go memRepository.GetGauges(ctx, gaugesCh, errCh)
			go memRepository.GetAdditionalGauges(ctx, gaugesCh, errCh)

			pollCount++

		case <-reportMain:
			counterCh <- types.Counters{dictionaries.CounterPollCount: pollCount}

			pollCount = 0

		case <-errCh:
			log.Printf("Error: %v\n", err)
		}
	}
}
