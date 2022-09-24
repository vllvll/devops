// Package storage Функционал для работы с хранилищами
package storage

import (
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/repositories"
	"github.com/vllvll/devops/internal/storage/file"
	"github.com/vllvll/devops/internal/types"
)

type statsStorage struct {
	config   *conf.ServerConfig
	consumer file.ConsumerFile
	producer file.ProducerFile
}

// NewStatsStorage Создание обработчика для восстановления данных в памяти при инициализации
func NewStatsStorage(serverConfig *conf.ServerConfig, consumer file.ConsumerFile, producer file.ProducerFile) *statsStorage {
	return &statsStorage{
		config:   serverConfig,
		consumer: consumer,
		producer: producer,
	}
}

// Save Сохранение данных метрик перед отключением приложения
func (s *statsStorage) Save(statsRepository repositories.StatsRepository) {
	if s.config.DatabaseDsn == "" {
		var metrics []types.Metrics

		gauges, counters := statsRepository.GetAll()

		for key, value := range gauges {
			flValue := float64(value)

			metrics = append(metrics, types.Metrics{
				ID:    key,
				MType: dictionaries.GaugeType,
				Value: &flValue,
			})
		}

		for key, value := range counters {
			iValue := int64(value)

			metrics = append(metrics, types.Metrics{
				ID:    key,
				MType: dictionaries.CounterType,
				Delta: &iValue,
			})
		}

		for _, m := range metrics {
			err := s.producer.WriteMetric(&m)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Start Восстановление метрик перед инициализацией приложения
func (s *statsStorage) Start(statsRepository repositories.StatsRepository) (repositories.StatsRepository, error) {
	if s.config.DatabaseDsn == "" && s.config.Restore {
		for {
			readMetric, err := s.consumer.ReadMetric()
			if err != nil {
				return statsRepository, nil
			}

			switch readMetric.MType {
			case dictionaries.GaugeType:
				statsRepository.UpdateGauge(readMetric.ID, types.Gauge(*readMetric.Value))

			case dictionaries.CounterType:
				statsRepository.UpdateCount(readMetric.ID, types.Counter(*readMetric.Delta))
			}
		}
	}

	return statsRepository, nil
}
