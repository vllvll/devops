package metric

type Сonstants struct {
	names []string
}

var metricNames = []string{
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
	"RandomValue",
	"PollCount",
}

func NewConstants() *Сonstants {
	return &Сonstants{
		names: metricNames,
	}
}

func (m Сonstants) In(str string) bool {
	for _, metric := range m.names {
		if metric == str {
			return true
		}
	}

	return false
}
