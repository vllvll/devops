// Package dictionaries содержит константы проекта
package dictionaries

// Названия типов, отправляемых значений
const (
	GaugeType   = "gauge"
	CounterType = "counter"
)

// GaugeRandomValue Ключ для случайного значения
const GaugeRandomValue = "RandomValue"

// CounterPollCount Какое количество раз мы получили данные о системе
const CounterPollCount = "PollCount"

// Список ключей, которые мы получаем из mem.VirtualMemory()
const (
	GaugeTotalMemoryValue     = "TotalMemory"
	GaugeFreeMemoryValue      = "FreeMemory"
	GaugeCPUutilization1Value = "CPUutilization1"
)

// Constants Структура для хранения списка ключей
type Constants struct {
	names []string
}

// DictionaryInterface Интерфейс для структуры Constants
type DictionaryInterface interface {
	In(str string) bool
}

// Список Gauge ключей, которые мы получаем из runtime.ReadMemStats
var gaugeNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
}

// NewMemConstants Получить структуру с ключами для Gauge значений
func NewMemConstants() DictionaryInterface {
	return &Constants{
		names: gaugeNames,
	}
}

// In Проверяем наличие ключа в словаре
func (m Constants) In(str string) bool {
	for _, name := range m.names {
		if name == str {
			return true
		}
	}

	return false
}
