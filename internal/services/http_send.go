package services

import (
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
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

func (c Sender) Send(gauges types.Gauges, pollCount types.Counter) error {
	for key, value := range gauges {
		var gaugeValue = float64(value)

		err := c.push(types.Metrics{
			ID:    key,
			MType: dictionaries.GaugeType,
			Value: &gaugeValue,
			Hash:  c.signer.GetHashGauge(key, gaugeValue),
		})

		if err != nil {
			return err
		}
	}

	var counterValue = int64(pollCount)

	err := c.push(types.Metrics{
		ID:    dictionaries.CounterPollCount,
		MType: dictionaries.CounterType,
		Delta: &counterValue,
		Hash:  c.signer.GetHashCounter(dictionaries.CounterPollCount, counterValue),
	})

	if err != nil {
		return err
	}

	return nil
}

func (c Sender) push(metric types.Metrics) error {
	_, err := c.Client.R().
		SetBody(metric).
		Post("/update/")

	if err != nil {
		return err
	}

	return nil
}
