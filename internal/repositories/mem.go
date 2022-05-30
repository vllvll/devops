package repositories

import (
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	"log"
	"math/rand"
	"reflect"
	"runtime"
)

type Mem struct {
	mem       runtime.MemStats
	gauges    types.Gauges
	constants dictionaries.DictionaryInterface
}

type MemRepository interface {
	GetGauges() types.Gauges
}

func NewMemRepository(constants dictionaries.DictionaryInterface) MemRepository {
	return &Mem{
		gauges:    types.Gauges{},
		constants: constants,
	}
}

func (m *Mem) GetGauges() types.Gauges {
	runtime.ReadMemStats(&m.mem)

	memReflect := reflect.ValueOf(&m.mem).Elem()

	for i := 0; i < memReflect.NumField(); i++ {
		var memValue types.Gauge
		memName := memReflect.Type().Field(i).Name

		if m.constants.In(memName) {
			switch memReflect.Field(i).Kind() {
			case reflect.Uint64:
				memValue = types.Gauge(memReflect.Field(i).Interface().(uint64))
			case reflect.Uint32:
				memValue = types.Gauge(memReflect.Field(i).Interface().(uint32))
			case reflect.Float64:
				memValue = types.Gauge(memReflect.Field(i).Interface().(float64))
			default:
				log.Fatalf("Error with get mem by key: %s", memReflect.Field(i).Kind())
			}

			m.gauges[memName] = memValue
		}
	}

	m.gauges[dictionaries.GaugeRandomValue] = types.Gauge(rand.Float64())

	return m.gauges
}
