package repositories

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/vllvll/devops/internal/dictionaries"
	"github.com/vllvll/devops/internal/types"
	"math/rand"
	"reflect"
	"runtime"
)

type Mem struct {
	mem       runtime.MemStats
	constants dictionaries.DictionaryInterface
}

type MemRepository interface {
	GetGauges(outGauges chan<- types.Gauges, errCh chan<- error)
	GetAdditionalGauges(outGauges chan<- types.Gauges, errCh chan<- error)
}

func NewMemRepository(constants dictionaries.DictionaryInterface) MemRepository {
	return &Mem{
		constants: constants,
	}
}

func (m *Mem) GetGauges(outGauges chan<- types.Gauges, errCh chan<- error) {
	defer func() {
		if err := recover(); err != nil {
			errCh <- fmt.Errorf("panic: %v", err)

			m.GetGauges(outGauges, errCh)
		}
	}()

	var gauges = types.Gauges{}

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
				errCh <- fmt.Errorf("error with get mem by key: %s", memReflect.Field(i).Kind())
			}

			gauges[memName] = memValue
		}
	}

	gauges[dictionaries.GaugeRandomValue] = types.Gauge(rand.Float64())

	outGauges <- gauges
}

func (m *Mem) GetAdditionalGauges(outGauges chan<- types.Gauges, errCh chan<- error) {
	defer func() {
		if err := recover(); err != nil {
			errCh <- fmt.Errorf("panic: %v", err)

			m.GetAdditionalGauges(outGauges, errCh)
		}
	}()

	var gauges = types.Gauges{}
	memory, err := mem.VirtualMemory()
	if err != nil {
		errCh <- err
	}

	cpu, err := load.Avg()
	if err != nil {
		errCh <- err
	}

	gauges[dictionaries.GaugeTotalMemoryValue] = types.Gauge(memory.Total)
	gauges[dictionaries.GaugeFreeMemoryValue] = types.Gauge(memory.Free)
	gauges[dictionaries.GaugeCPUutilization1Value] = types.Gauge(cpu.Load1)

	outGauges <- gauges
}
