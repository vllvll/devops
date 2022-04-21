package metric

type metrics map[string]Gauge

type Repository struct {
	Metrics   metrics
	PollCount Counter
}

func NewRepository() *Repository {
	return &Repository{
		Metrics: metrics{},
	}
}

func (m *Repository) UpdateMetric(key string, value Gauge) {
	m.Metrics[key] = value
}

func (m *Repository) UpdateCount(value Counter) {
	m.PollCount += value
}
