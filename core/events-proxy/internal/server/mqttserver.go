package server

import (
	"context"
	"errors"
	"eventsproxy/internal/config"
	"eventsproxy/internal/service"
	"log"
	"strings"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
)

var (
	ErrStopped     = errors.New("server is already stopped")
	ErrUnauthorize = errors.New("failed to auth")
)

type MqttServer struct {
	notify     chan error
	ctx        context.Context
	cancel     func()
	mqttServer *mqtt.Server
}

func NewMqttServer(cfg config.MqttServerConfig, svc service.ProxyService) (srv MqttServer, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	server := mqtt.New(nil)

	err = server.AddHook(newHook(svc), nil)
	if err != nil {
		cancel()
		return
	}

	var builder strings.Builder
	builder.Grow(len(cfg.Host) + 1 + len(cfg.Port))
	builder.WriteString(cfg.Host)
	builder.WriteRune(':')
	builder.WriteString(cfg.Port)

	tcp := listeners.NewTCP(cfg.Id, builder.String(), nil)

	err = server.AddListener(tcp)
	if err != nil {
		cancel()
		return
	}

	srv = MqttServer{
		notify:     make(chan error),
		ctx:        ctx,
		cancel:     cancel,
		mqttServer: server,
	}

	return
}

func (srv *MqttServer) Run() error {
	select {
	case <-srv.ctx.Done():
		return ErrStopped
	default:
	}

	go func() {
		err := srv.mqttServer.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (srv MqttServer) Notify() chan error {
	return srv.notify
}

func (srv MqttServer) Close() {
	srv.cancel()
}
