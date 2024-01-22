package server

import (
	"encoding/json"
	"eventsproxy/internal/config"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

func NewStatusServer(cfg config.StatusServerConfig, cfgProj config.ProjectConfig) (*http.Server, *atomic.Bool) {
	statusAddr := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	isReady := &atomic.Bool{}
	isReady.Store(false)

	mux := http.DefaultServeMux

	mux.Handle(cfg.LivelinessPath, LivenessHandler(cfgProj))
	mux.Handle(cfg.ReadinessPath, ReadinessHandler(isReady, cfgProj))

	statusServer := &http.Server{
		Addr:              statusAddr,
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
	}

	return statusServer, isReady
}

func VersionHandler(w http.ResponseWriter, r *http.Request, cfg config.ProjectConfig) {
	data := map[string]interface{}{
		"version": cfg.Version,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Error().Err(err).Msg("StatusServer: Data encoding error")
	}
}

// LivenessHandler handles liveness probe requests and responds with an HTTP status OK (200).
func LivenessHandler(cfg config.ProjectConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		VersionHandler(w, r, cfg)
	})
}

// ReadinessHandler returns an HTTP handler function that checks the server's readiness based on the provided atomic boolean value.
func ReadinessHandler(isReady *atomic.Bool, cfg config.ProjectConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isReady == nil || !isReady.Load() {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)

			return
		}

		VersionHandler(w, r, cfg)
	})
}
