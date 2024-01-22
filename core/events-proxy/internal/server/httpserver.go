package server

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

type HttpServerMux struct {
	servers []*http.Server
	isReady *atomic.Bool
	notify  chan error
	ctx     context.Context
	cancel  func()
}

func NewHttpServerMux(isReady *atomic.Bool, servers ...*http.Server) *HttpServerMux {
	ctx, cancel := context.WithCancel(context.Background())
	return &HttpServerMux{
		isReady: isReady,
		servers: servers,
		notify:  make(chan error),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (srv *HttpServerMux) serverRunRoutine(server *http.Server) {
	log.Info().Str("address", server.Addr).Msg("server is running")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error().Err(err).Msg("Failed running server")
		srv.notify <- err
	}
}

func (srv *HttpServerMux) Run() error {
	select {
	case <-srv.ctx.Done():
		return ErrStopped
	default:
	}

	for _, server := range srv.servers {
		go srv.serverRunRoutine(server)
	}

	go func(ready *atomic.Bool) {
		time.Sleep(2 * time.Second)
		ready.Store(true)
		log.Info().Msg("Service is ready")
	}(srv.isReady)

	return nil
}

func (srv HttpServerMux) Notify() chan error {
	return srv.notify
}

func (srv HttpServerMux) Close() {
	srv.cancel()
}
