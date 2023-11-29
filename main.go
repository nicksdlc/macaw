package main

import (
	"log"
	"strings"
	"time"

	"github.com/nicksdlc/macaw/config"
	"github.com/nicksdlc/macaw/context"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	cfg := readConfig()

	ctx, err := context.BuildContext(cfg)
	if err != nil {
		log.Panic(err.Error())
	}
	if cfg.DumpMetrics.Enabled {
		collectMetrics(cfg.DumpMetrics.Frequency)
	}

	ctx.Run()
}

func readConfig() config.Configuration {
	return config.Read("config")
}

func collectMetrics(frequency int) {
	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	go func() {
		for range ticker.C {
			logPrometheusMetrics()
		}
	}()
}

func logPrometheusMetrics() {
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		log.Printf("Error gathering Prometheus metrics: %v\n", err)
		return
	}

	for _, mf := range metricFamilies {
		if strings.Contains(*mf.Name, "macaw") {
			for _, m := range mf.Metric {
				log.Printf(" [metrics] %s: %v\n", *mf.Name, m.Gauge.GetValue())
			}
		}
	}
}
