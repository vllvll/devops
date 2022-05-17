package repositories

import (
	"fmt"
	"github.com/vllvll/devops/internal/types"
)

type Stats struct {
	Gauges   types.Gauges
	Counters types.Counters
}

type StatsRepository interface {
	UpdateGauge(key string, value types.Gauge)
	UpdateCount(key string, value types.Counter)
	GetAll() (map[string]types.Gauge, map[string]types.Counter)
	GetGaugeByKey(key string) (types.Gauge, error)
	GetCounterByKey(key string) (types.Counter, error)
}

func NewStatsRepository() StatsRepository {
	return &Stats{
		Gauges:   types.Gauges{},
		Counters: types.Counters{},
	}
}

func (r *Stats) UpdateGauge(key string, value types.Gauge) {
	r.Gauges[key] = value
}

func (r *Stats) UpdateCount(key string, value types.Counter) {
	r.Counters[key] += value
}

func (r *Stats) GetAll() (map[string]types.Gauge, map[string]types.Counter) {
	return r.Gauges, r.Counters
}

func (r *Stats) GetGaugeByKey(key string) (types.Gauge, error) {
	if value, ok := r.Gauges[key]; ok {
		return value, nil
	}

	return types.Gauge(0), fmt.Errorf("%s key doesn't exists", key)
}

func (r *Stats) GetCounterByKey(key string) (types.Counter, error) {
	if value, ok := r.Counters[key]; ok {
		return value, nil
	}

	return types.Counter(0), fmt.Errorf("%s key doesn't exists", key)
}
