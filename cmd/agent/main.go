package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

type gauge float64
type counter int64

var pollInterval = 2
var reportInterval = 10

type Metrics map[string]gauge

func main() {
	var pollCount counter
	var mem runtime.MemStats
	metrics := Metrics{}

	var pollTick = time.Tick(time.Duration(pollInterval) * time.Second)
	var reportTick = time.Tick(time.Duration(reportInterval) * time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-c:
			fmt.Println("Break the loop")

			return
		case <-pollTick:
			fmt.Println("poll tick")

			pollCount++
			runtime.ReadMemStats(&mem)

			metrics["Alloc"] = gauge(mem.Alloc)
			metrics["BuckHashSys"] = gauge(mem.BuckHashSys)
			metrics["Frees"] = gauge(mem.Frees)
			metrics["GCCPUFraction"] = gauge(mem.GCCPUFraction)
			metrics["GCSys"] = gauge(mem.GCSys)
			metrics["HeapAlloc"] = gauge(mem.HeapAlloc)
			metrics["HeapIdle"] = gauge(mem.HeapIdle)
			metrics["HeapInuse"] = gauge(mem.HeapInuse)
			metrics["HeapObjects"] = gauge(mem.HeapObjects)
			metrics["HeapReleased"] = gauge(mem.HeapReleased)
			metrics["HeapSys"] = gauge(mem.HeapSys)
			metrics["LastGC"] = gauge(mem.LastGC)
			metrics["Lookups"] = gauge(mem.Lookups)
			metrics["MCacheInuse"] = gauge(mem.MCacheInuse)
			metrics["MCacheSys"] = gauge(mem.MCacheSys)
			metrics["MSpanInuse"] = gauge(mem.MSpanInuse)
			metrics["MSpanSys"] = gauge(mem.MSpanSys)
			metrics["Mallocs"] = gauge(mem.Mallocs)
			metrics["NextGC"] = gauge(mem.NextGC)
			metrics["NumForcedGC"] = gauge(mem.NumForcedGC)
			metrics["NumGC"] = gauge(mem.NumGC)
			metrics["OtherSys"] = gauge(mem.OtherSys)
			metrics["PauseTotalNs"] = gauge(mem.PauseTotalNs)
			metrics["StackInuse"] = gauge(mem.StackInuse)
			metrics["StackSys"] = gauge(mem.StackSys)
			metrics["Sys"] = gauge(mem.Sys)
			metrics["TotalAlloc"] = gauge(mem.TotalAlloc)
			metrics["RandomValue"] = gauge(rand.Float64())

		case <-reportTick:
			fmt.Println("report tick")

			for metric, value := range metrics {
				err := updateGauge(metric, value)
				if err != nil {
					log.Printf("can't send report: %v\n", err)
				}
			}

			err := updateCounter("PollCount", pollCount)
			if err != nil {
				log.Printf("can't send report: %v\n", err)
			}

			pollCount = 0
		}
	}
}

func updateGauge(name string, value gauge) error {
	_, err := http.Post(
		fmt.Sprintf(
			"http://127.0.0.1:8080/update/gauge/%s/%s",
			name,
			strconv.FormatFloat(float64(value), 'f', 6, 64),
		),
		"text/plain",
		bytes.NewBuffer([]byte("")),
	)

	if err != nil {
		return err
	}

	return nil
}

func updateCounter(name string, value counter) error {
	_, err := http.Post(
		fmt.Sprintf("http://127.0.0.1:8080/update/counter/%s/%d", name, value),
		"text/plain",
		bytes.NewBuffer([]byte("")),
	)

	if err != nil {
		return err
	}

	return nil
}
