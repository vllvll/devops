package services

import (
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	"sync"
	"time"
)

type Sender struct {
	Client *resty.Client
	signer Signer
}

func NewSendClient(AgentConfig *conf.AgentConfig, signer Signer) *Sender {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	client := resty.New().
		SetBaseURL(AgentConfig.AddressWithHTTP()).
		SetHeader("Content-Type", "application/json")

	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	return &Sender{
		Client: client,
		signer: signer,
	}
}

func (c Sender) Prepare(gaugesIn <-chan types.Gauges, countersIn <-chan types.Counters, metricCh chan<- types.Metrics) {
	go func() {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()
			for gauges := range gaugesIn {
				for key, value := range gauges {
					gaugeValue := float64(value)

					metricCh <- types.Metrics{
						ID:    key,
						MType: dictionaries.GaugeType,
						Value: &gaugeValue,
						Hash:  c.signer.GetHashGauge(key, gaugeValue),
					}
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for counters := range countersIn {
				for key, value := range counters {
					var counterValue = int64(value)

					metricCh <- types.Metrics{
						ID:    key,
						MType: dictionaries.CounterType,
						Delta: &counterValue,
						Hash:  c.signer.GetHashCounter(key, counterValue),
					}
				}
			}
		}()

		wg.Wait()
		close(metricCh)
	}()
}

func (c Sender) Send(metricCh <-chan types.Metrics, reportTick <-chan time.Time, errCh chan<- error) {
	var metrics []types.Metrics
	for {
		select {
		case <-reportTick:
			err := c.push(metrics)
			if err != nil {
				errCh <- err
			}

		case metric, ok := <-metricCh:
			metrics = append(metrics, metric)

			if !ok {

				return
			}
		}
	}
}

func (c Sender) push(metrics []types.Metrics) error {
	_, err := c.Client.R().
		SetBody(metrics).
		Post("/updates/")

	if err != nil {
		return err
	}

	return nil
}
