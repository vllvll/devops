package dictionaries

const GaugeType = "gauge"
const CounterType = "counter"

const GaugeRandomValue = "RandomValue"
const CounterPollCount = "PollCount"

const GaugeTotalMemoryValue = "TotalMemory"
const GaugeFreeMemoryValue = "FreeMemory"
const GaugeCPUutilization1Value = "CPUutilization1"

type Constants struct {
	names []string
}

type DictionaryInterface interface {
	In(str string) bool
}

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

func NewMemConstants() DictionaryInterface {
	return &Constants{
		names: gaugeNames,
	}
}

func (m Constants) In(str string) bool {
	for _, name := range m.names {
		if name == str {
			return true
		}
	}

	return false
}
