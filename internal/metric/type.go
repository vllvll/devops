package metric

type Gauge float64
type Counter int64

type Metrics map[string]Gauge
