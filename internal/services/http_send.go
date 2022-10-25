// Package services содержит вспомогательные сервисы
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
)

type Sender struct {
	Client  *resty.Client // HTTP клиент
	signer  Signer        // Сервис для подписи данных
	encrypt Encrypt       // Сервис для ассиметричного шифрования
}

// NewSendClient Создание сервиса для отправки данных из агента на сервер
func NewSendClient(AgentConfig *conf.AgentConfig, signer Signer, encrypt Encrypt) (*Sender, error) {
	ip, err := AgentConfig.GetServiceIP()
	if err != nil {
		return nil, fmt.Errorf("IP адрес не найден")
	}

	client := resty.New().
		SetBaseURL(AgentConfig.AddressWithHTTP()).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Real-IP", ip)

	return &Sender{
		Client:  client,
		signer:  signer,
		encrypt: encrypt,
	}, nil
}

// Prepare Подготовка метрик для отправки на сервер
func (c Sender) Prepare(ctx context.Context, gaugesIn <-chan types.Gauges, countersIn <-chan types.Counters, metricCh chan<- types.Metrics, errCh chan<- error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errCh <- fmt.Errorf("panic: %v", err)

				c.Prepare(ctx, gaugesIn, countersIn, metricCh, errCh)
			}
		}()

		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case gauges, ok := <-gaugesIn:
					if !ok {
						return
					}

					for key, value := range gauges {
						gaugeValue := float64(value)

						metricCh <- types.Metrics{
							ID:    key,
							MType: dictionaries.GaugeType,
							Value: &gaugeValue,
							Hash:  c.signer.GetHashGauge(key, gaugeValue),
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case counters, ok := <-countersIn:
					if !ok {
						return
					}

					for key, value := range counters {
						var counterValue = int64(value)

						metricCh <- types.Metrics{
							ID:    key,
							MType: dictionaries.CounterType,
							Delta: &counterValue,
							Hash:  c.signer.GetHashCounter(key, counterValue),
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Wait()
		close(metricCh)
	}()
}

// Send Отправка метрик на сервер с определенной периодичностью
func (c Sender) Send(ctx context.Context, metricCh <-chan types.Metrics, reportTick <-chan time.Time, errCh chan<- error) {
	defer func() {
		if err := recover(); err != nil {
			errCh <- fmt.Errorf("panic: %v", err)

			c.Send(ctx, metricCh, reportTick, errCh)
		}
	}()

	var metrics = make([]types.Metrics, 0, 100)
	for {
		select {
		case <-reportTick:
			err := c.push(&metrics)
			if err != nil {
				errCh <- err
			}

			metrics = metrics[:0]

		case metric, ok := <-metricCh:
			if !ok {
				return
			}

			metrics = append(metrics, metric)

		case <-ctx.Done():
			return
		}
	}
}

// Внутренний метод для отправки метрик на сервер
func (c Sender) push(metrics *[]types.Metrics) error {
	content, err := json.Marshal(*metrics)
	if err != nil {
		return err
	}

	if c.encrypt != nil {
		content, err = c.encrypt.Encrypt(content)
		if err != nil {
			return err
		}
	}

	_, err = c.Client.R().
		SetBody(content).
		Post("/updates/")

	if err != nil {
		return err
	}

	return nil
}
