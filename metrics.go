package zetka

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	WebSocketConnectionCounter = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "",
		Subsystem: "zetka",
		Name:      "websocket_connection_count",
	})

	All = []prometheus.Collector{
		WebSocketConnectionCounter,
	}
)

// registerMetrics all metrics for zetka
func registerMetrics() error {
	for _, collector := range All {
		if err := prometheus.Register(collector); err != nil {
			return err
		}
	}

	return nil
}
