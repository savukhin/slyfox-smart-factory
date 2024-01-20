package main

import (
	"eventsproxy/internal/config"
	"eventsproxy/internal/connections"
	"eventsproxy/internal/server"
	"eventsproxy/internal/service"
	"eventsproxy/internal/service/producer"
	"eventsproxy/internal/service/repo"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const (
	HOST = "localhost"
	PORT = "1883"
	TYPE = "tcp"
)

func main() {
	cfg := config.ReadConfig()

	db, err := connections.CreatePostgres(cfg.Postgres)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	nc, js, err := connections.CreateNatsJetstream(cfg.Nats)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	producer := producer.NewNastProducer(js)
	proxyRepo := repo.NewUserRepo(db)
	proxySvc := service.NewProxyService(&proxyRepo, &producer)

	serverCfg := server.NewMqttServerConfig()
	server, err := server.NewMqttServer(serverCfg, &proxySvc)

	if err != nil {
		panic(err)
	}
	err = server.Run()
	if err != nil {
		panic(err)
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("app - Run - signal: " + s.String())
	case err := <-server.Notify():
		fmt.Println(err)
	}
}
