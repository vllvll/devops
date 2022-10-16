// Модуль agent отправляет информацию о состоянии
package main

import (
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
	buildVersion string
	buildDate    string
	buildCommit  string
)

const BuildTemplate = `
Build version: {{if .version}}{{ .version }}{{ else }}N/A{{ end }}
Build date: {{if .date}}{{ .date }}{{ else }}N/A{{ end }}
Build commit: {{if .commit}}{{ .commit }}{{ else }}N/A{{ end }}
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

	var pollTick = time.Tick(config.PollInterval)
	var reportTick = time.Tick(config.ReportInterval)
	var reportMain = time.Tick(config.ReportInterval)

	var pollCount types.Counter

	crypt, err := services.NewMetricEncrypt(config.CryptoKey)
	if err != nil {
		log.Fatalf("Ошибка с инициализацией сервиса шифрования: %v", err)
	}

	signer := services.NewMetricSigner(config.Key)
	sender, err := services.NewSendClient(config, signer, crypt)
	if err != nil {
		log.Fatalf("Ошибка с иницализацией сервиса http клиента: %v", err)
	}
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
