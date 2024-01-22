package server

import (
	"eventsproxy/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetricsServer(cfg config.MetricsServerConfig) *http.Server {
	addr := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
	}

	return metricsServer
}
