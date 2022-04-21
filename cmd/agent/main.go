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
	metrics := metric.Metrics{}

	client := metric.NewClient()

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

			metrics["Alloc"] = metric.Gauge(mem.Alloc)
			metrics["BuckHashSys"] = metric.Gauge(mem.BuckHashSys)
			metrics["Frees"] = metric.Gauge(mem.Frees)
			metrics["GCCPUFraction"] = metric.Gauge(mem.GCCPUFraction)
			metrics["GCSys"] = metric.Gauge(mem.GCSys)
			metrics["HeapAlloc"] = metric.Gauge(mem.HeapAlloc)
			metrics["HeapIdle"] = metric.Gauge(mem.HeapIdle)
			metrics["HeapInuse"] = metric.Gauge(mem.HeapInuse)
			metrics["HeapObjects"] = metric.Gauge(mem.HeapObjects)
			metrics["HeapReleased"] = metric.Gauge(mem.HeapReleased)
			metrics["HeapSys"] = metric.Gauge(mem.HeapSys)
			metrics["LastGC"] = metric.Gauge(mem.LastGC)
			metrics["Lookups"] = metric.Gauge(mem.Lookups)
			metrics["MCacheInuse"] = metric.Gauge(mem.MCacheInuse)
			metrics["MCacheSys"] = metric.Gauge(mem.MCacheSys)
			metrics["MSpanInuse"] = metric.Gauge(mem.MSpanInuse)
			metrics["MSpanSys"] = metric.Gauge(mem.MSpanSys)
			metrics["Mallocs"] = metric.Gauge(mem.Mallocs)
			metrics["NextGC"] = metric.Gauge(mem.NextGC)
			metrics["NumForcedGC"] = metric.Gauge(mem.NumForcedGC)
			metrics["NumGC"] = metric.Gauge(mem.NumGC)
			metrics["OtherSys"] = metric.Gauge(mem.OtherSys)
			metrics["PauseTotalNs"] = metric.Gauge(mem.PauseTotalNs)
			metrics["StackInuse"] = metric.Gauge(mem.StackInuse)
			metrics["StackSys"] = metric.Gauge(mem.StackSys)
			metrics["Sys"] = metric.Gauge(mem.Sys)
			metrics["TotalAlloc"] = metric.Gauge(mem.TotalAlloc)
			metrics["RandomValue"] = metric.Gauge(rand.Float64())

		case <-reportTick:
			err := client.Send(metrics, pollCount)
			if err != nil {
				log.Printf("can't send report: %v\n", err)
			}

			pollCount = 0
		}
	}
}
