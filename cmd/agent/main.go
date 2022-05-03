package main

import (
	"fmt"
	"github.com/vllvll/devops/internal/metric"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var pollInterval = 2
var reportInterval = 10

var pollTick = time.Tick(time.Duration(pollInterval) * time.Second)
var reportTick = time.Tick(time.Duration(reportInterval) * time.Second)

func main() {
	var mem runtime.MemStats
	var pollCount metric.Counter
	gauges := metric.Gauges{}

	sender := metric.NewClient()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			fmt.Println("Graceful shutdown")

			return
		case <-pollTick:
			pollCount++
			runtime.ReadMemStats(&mem)

			gauges["Alloc"] = metric.Gauge(mem.Alloc)
			gauges["BuckHashSys"] = metric.Gauge(mem.BuckHashSys)
			gauges["Frees"] = metric.Gauge(mem.Frees)
			gauges["GCCPUFraction"] = metric.Gauge(mem.GCCPUFraction)
			gauges["GCSys"] = metric.Gauge(mem.GCSys)
			gauges["HeapAlloc"] = metric.Gauge(mem.HeapAlloc)
			gauges["HeapIdle"] = metric.Gauge(mem.HeapIdle)
			gauges["HeapInuse"] = metric.Gauge(mem.HeapInuse)
			gauges["HeapObjects"] = metric.Gauge(mem.HeapObjects)
			gauges["HeapReleased"] = metric.Gauge(mem.HeapReleased)
			gauges["HeapSys"] = metric.Gauge(mem.HeapSys)
			gauges["LastGC"] = metric.Gauge(mem.LastGC)
			gauges["Lookups"] = metric.Gauge(mem.Lookups)
			gauges["MCacheInuse"] = metric.Gauge(mem.MCacheInuse)
			gauges["MCacheSys"] = metric.Gauge(mem.MCacheSys)
			gauges["MSpanInuse"] = metric.Gauge(mem.MSpanInuse)
			gauges["MSpanSys"] = metric.Gauge(mem.MSpanSys)
			gauges["Mallocs"] = metric.Gauge(mem.Mallocs)
			gauges["NextGC"] = metric.Gauge(mem.NextGC)
			gauges["NumForcedGC"] = metric.Gauge(mem.NumForcedGC)
			gauges["NumGC"] = metric.Gauge(mem.NumGC)
			gauges["OtherSys"] = metric.Gauge(mem.OtherSys)
			gauges["PauseTotalNs"] = metric.Gauge(mem.PauseTotalNs)
			gauges["StackInuse"] = metric.Gauge(mem.StackInuse)
			gauges["StackSys"] = metric.Gauge(mem.StackSys)
			gauges["Sys"] = metric.Gauge(mem.Sys)
			gauges["TotalAlloc"] = metric.Gauge(mem.TotalAlloc)
			gauges["RandomValue"] = metric.Gauge(rand.Float64())

		case <-reportTick:
			err := sender.Send(gauges, pollCount)
			if err != nil {
				log.Printf("can't send report: %v\n", err)
			}

			pollCount = 0
		}
	}
}
