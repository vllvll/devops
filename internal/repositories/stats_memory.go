package repositories

import (
	"fmt"

	"github.com/vllvll/devops/internal/types"
)

type StatsMemory struct {
	Gauges   types.Gauges
	Counters types.Counters
}

type StatsRepository interface {
	UpdateGauge(key string, value types.Gauge)
	UpdateCount(key string, value types.Counter)
	GetAll() (map[string]types.Gauge, map[string]types.Counter)
	GetGaugeByKey(key string) (types.Gauge, error)
	GetCounterByKey(key string) (types.Counter, error)
	UpdateAll(gauges types.Gauges, counters types.Counters) error
}

// NewStatsMemoryRepository Создание репозитория, который отвечает за хранение метрик в оперативной памяти
func NewStatsMemoryRepository() StatsRepository {
	return &StatsMemory{
		Gauges:   types.Gauges{},
		Counters: types.Counters{},
	}
}

// UpdateGauge Обновить значение метрики с типом Gauge в оперативной памяти
func (s *StatsMemory) UpdateGauge(key string, value types.Gauge) {
	s.Gauges[key] = value
}

// UpdateCount Обновить значение метрики с типом Counter в оперативной памяти
func (s *StatsMemory) UpdateCount(key string, value types.Counter) {
	s.Counters[key] += value
}

// GetAll Получение всех метрик из оперативной памяти
func (s *StatsMemory) GetAll() (map[string]types.Gauge, map[string]types.Counter) {
	return s.Gauges, s.Counters
}

// GetGaugeByKey Получить значение метрики типа Gauge по ключу из оперативной памяти
func (s *StatsMemory) GetGaugeByKey(key string) (types.Gauge, error) {
	if value, ok := s.Gauges[key]; ok {
		return value, nil
	}

	return types.Gauge(0), fmt.Errorf("%s key doesn't exists", key)
}

// GetCounterByKey Получить значение метрики типа Counter по ключу из оперативной памяти
func (s *StatsMemory) GetCounterByKey(key string) (types.Counter, error) {
	if value, ok := s.Counters[key]; ok {
		return value, nil
	}

	return types.Counter(0), fmt.Errorf("%s key doesn't exists", key)
}

// UpdateAll Обновление всех значений типов Gauge и Counter в оперативной памяти
func (s *StatsMemory) UpdateAll(gauges types.Gauges, counters types.Counters) error {
	for key, value := range gauges {
		s.Gauges[key] = value
	}

	for key, value := range counters {
		s.Counters[key] += value
	}

	return nil
}
