package main

import (
	"fmt"
	conf "github.com/vllvll/devops/internal/config"
	"github.com/vllvll/devops/internal/metric"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"
)

func main() {
	config, err := conf.CreateConfig()
	if err != nil {
		panic("Конфиг не загружен")
	}

	var pollTick = time.Tick(config.PollInterval)
	var reportTick = time.Tick(config.ReportInterval)

	var mem runtime.MemStats
	var pollCount metric.Counter
	gauges := metric.Gauges{}

	sender := metric.NewClient(config)
	fields := metric.NewConstants()

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

			memReflect := reflect.ValueOf(&mem).Elem()

			for i := 0; i < memReflect.NumField(); i++ {
				var memValue metric.Gauge
				memName := memReflect.Type().Field(i).Name

				if fields.In(memName) {
					switch memReflect.Field(i).Kind() {
					case reflect.Uint64:
						memValue = metric.Gauge(memReflect.Field(i).Interface().(uint64))
					case reflect.Uint32:
						memValue = metric.Gauge(memReflect.Field(i).Interface().(uint32))
					case reflect.Float64:
						memValue = metric.Gauge(memReflect.Field(i).Interface().(float64))
					default:
						panic("Это ключ не имеет обработанного типа")
					}

					gauges[memName] = memValue
				}
			}

			gauges[metric.GaugeRandomValue] = metric.Gauge(rand.Float64())

		case <-reportTick:
			err := sender.Send(gauges, pollCount)
			if err != nil {
				log.Printf("can't send report: %v\n", err)
			}

			pollCount = 0
		}
	}
}
