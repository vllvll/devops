// Package services содержит вспомогательные сервисы
package services

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	pb "github.com/vllvll/devops/proto"
)

type GRPCSender struct {
	Client  pb.MetricsClient
	signer  Signer  // Сервис для подписи данных
	encrypt Encrypt // Сервис для ассиметричного шифрования
	ip      string  // IP клиента
}

// NewGRPCSendClient Создание сервиса для отправки данных из агента на сервер
func NewGRPCSendClient(AgentConfig *conf.AgentConfig, signer Signer, encrypt Encrypt) (*GRPCSender, error) {
	var ip string
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addresses {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			ip = ipNet.IP.String()
		}
	}

	if ip == "" {
		return nil, fmt.Errorf("ip адрес не найден")
	}

	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(AgentConfig.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	c := pb.NewMetricsClient(conn)

	return &GRPCSender{
		Client:  c,
		signer:  signer,
		encrypt: encrypt,
		ip:      ip,
	}, nil
}

// Prepare Подготовка метрик для отправки на сервер
func (c GRPCSender) Prepare(gaugesIn <-chan types.Gauges, countersIn <-chan types.Counters, metricCh chan<- types.Metrics, errCh chan<- error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errCh <- fmt.Errorf("panic: %v", err)

				c.Prepare(gaugesIn, countersIn, metricCh, errCh)
			}
		}()

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

// Send Отправка метрик на сервер с определенной периодичностью
func (c GRPCSender) Send(metricCh <-chan types.Metrics, reportTick <-chan time.Time, errCh chan<- error) {
	defer func() {
		if err := recover(); err != nil {
			errCh <- fmt.Errorf("panic: %v", err)

			c.Send(metricCh, reportTick, errCh)
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
			metrics = append(metrics, metric)

			if !ok {

				return
			}
		}
	}
}

// Внутренний метод для отправки метрик на сервер
func (c GRPCSender) push(metrics *[]types.Metrics) error {
	var bulkMetrics []*pb.Metric

	for _, metric := range *metrics {
		var metricType pb.Metric_Type
		switch metric.MType {
		case dictionaries.CounterType:
			metricType = pb.Metric_COUNTER
		case dictionaries.GaugeType:
			metricType = pb.Metric_GAUGE
		}

		bulkMetrics = append(bulkMetrics, &pb.Metric{
			Id:    metric.ID,
			MType: metricType,
			Delta: metric.Delta,
			Value: metric.Value,
			Hash:  &metric.Hash,
		})
	}

	request := pb.AddBulkMetricsRequest{}
	request.Metrics = &pb.BulkMetrics{Metrics: bulkMetrics}

	md := metadata.New(map[string]string{"ip": c.ip})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := c.Client.BulkSaveMetrics(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}
