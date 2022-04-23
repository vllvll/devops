package metric

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c Client) Send(metrics Metrics, pollCount Counter) error {
	for key, value := range metrics {
		err := c.sendGauge(key, value)
		if err != nil {
			return err
		}
	}

	err := c.sendCounter("PollCount", pollCount)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) sendGauge(name string, value Gauge) error {
	response, err := http.Post(
		fmt.Sprintf(
			"http://127.0.0.1:8080/update/gauge/%s/%s",
			name,
			strconv.FormatFloat(float64(value), 'f', 3, 64),
		),
		"text/plain",
		bytes.NewBuffer([]byte("")),
	)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}

func (c Client) sendCounter(name string, value Counter) error {
	response, err := http.Post(
		fmt.Sprintf(
			"http://127.0.0.1:8080/update/counter/%s/%s",
			name,
			strconv.FormatInt(int64(value), 10),
		),
		"text/plain",
		bytes.NewBuffer([]byte("")),
	)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
