package metric

import "fmt"

type gauges map[string]Gauge
type counters map[string]Counter

type Repository struct {
	Gauges   gauges
	Counters counters
}

func NewRepository() *Repository {
	return &Repository{
		Gauges:   gauges{},
		Counters: counters{},
	}
}

func (r *Repository) UpdateMetric(key string, value Gauge) {
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
