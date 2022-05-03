package metric

const GaugeType = "gauge"
const CounterType = "counter"

const GaugeRandomValue = "RandomValue"
const CounterPollCount = "PollCount"

type Constants struct {
	names []string
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

func NewConstants() *Constants {
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
