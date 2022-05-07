package metric

import (
	"fmt"
)

type Repository struct {
	Gauges   Gauges
	Counters Counters
}

type RepositoryInterface interface {
	UpdateGauge(key string, value Gauge)
	UpdateCount(key string, value Counter)
	GetAll() (map[string]Gauge, map[string]Counter)
	GetGaugeByKey(key string) (Gauge, error)
	GetCounterByKey(key string) (Counter, error)
}

func NewRepository() RepositoryInterface {
	return &Repository{
		Gauges:   Gauges{},
		Counters: Counters{},
	}
}

func (r *Repository) UpdateGauge(key string, value Gauge) {
	r.Gauges[key] = value
}

func (r *Repository) UpdateCount(key string, value Counter) {
	r.Counters[key] += value
}

func (r *Repository) GetAll() (map[string]Gauge, map[string]Counter) {
	return r.Gauges, r.Counters
}

func (r *Repository) GetGaugeByKey(key string) (Gauge, error) {
	if value, ok := r.Gauges[key]; ok {
		return value, nil
	}

	return Gauge(0), fmt.Errorf("%s key doesn't exists", key)
}

func (r *Repository) GetCounterByKey(key string) (Counter, error) {
	if value, ok := r.Counters[key]; ok {
		return value, nil
	}

	return Counter(0), fmt.Errorf("%s key doesn't exists", key)
}
