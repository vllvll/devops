package metric

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

func (m *Repository) UpdateMetric(key string, value Gauge) {
	m.Gauges[key] = value

	//fmt.Println(key, value)
}

func (m *Repository) UpdateCount(key string, value Counter) {
	m.Counters[key] += value

	//fmt.Println(key, m.Counters[key])
}
