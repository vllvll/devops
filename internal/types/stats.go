package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type Gauge float64

func (g Gauge) Value() (driver.Value, error) {
	return strconv.FormatFloat(float64(g), 'f', -1, 64), nil
}

func (g *Gauge) Scan(value interface{}) error {
	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("cannot scan value. %w", err)
	}

	v, err := strconv.ParseFloat(sv.(string), 64)
	if err != nil {
		return fmt.Errorf("cannot scan value. cannot convert value to string")
	}

	*g = Gauge(v)

	return nil
}

type Counter int64

func (c Counter) Value() (driver.Value, error) {
	return strconv.FormatInt(int64(c), 10), nil
}

func (c *Counter) Scan(value interface{}) error {
	sv, err := driver.String.ConvertValue(value)
	if err != nil {
		return fmt.Errorf("cannot scan value. %w", err)
	}

	v, err := strconv.ParseInt(sv.(string), 10, 64)

	if err != nil {
		return fmt.Errorf("cannot scan value. cannot convert value to string")
	}

	*c = Counter(v)

	return nil
}

type Counters map[string]Counter
type Gauges map[string]Gauge

type Metrics struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // Значение хеш-функции
}
