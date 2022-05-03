package metric

import (
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
)

type Sender struct {
	Client *resty.Client
}

func NewClient() *Sender {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	client := resty.New().
		SetBaseURL("http://127.0.0.1:8080").
		SetHeader("Content-Type", "application/json")

	client.JSONMarshal = json.Marshal
	client.JSONUnmarshal = json.Unmarshal

	return &Sender{
		Client: client,
	}
}

func (c Sender) Send(gauges Gauges, pollCount Counter) error {
	for key, value := range gauges {
		var gaugeValue = float64(value)

		err := c.push(Metrics{
			ID:    key,
			MType: GaugeType,
			Value: &gaugeValue,
		})

		if err != nil {
			return err
		}
	}

	var counterValue = int64(pollCount)

	err := c.push(Metrics{
		ID:    CounterPollCount,
		MType: CounterType,
		Delta: &counterValue,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c Sender) push(metric Metrics) error {
	_, err := c.Client.R().
		SetBody(metric).
		Post("/update/")

	if err != nil {
		return err
	}

	return nil
}
