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

	"github.com/pressly/goose"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Info().Any("cfg", cfg).Msg("Starting service...")

	db, err := connections.CreatePostgres(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect database")
		return
	}
	defer db.Close()

	if err = goose.Up(db.DB, cfg.Database.MigrationsFolder); err != nil {
		log.Fatal().Err(err).Msg("Cannot migrate")
		return
	}

	nc, js, err := connections.CreateNatsJetstream(cfg.Nats)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create nats jetstream")
		return
	}
	defer nc.Close()

	producer := producer.NewNastProducer(js)
	proxyRepo := repo.NewUserRepo(db)
	proxySvc := service.NewProxyService(&proxyRepo, &producer)

	server, err := server.NewMqttServer(cfg.MqttServer, &proxySvc)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create mqtt server")
		return
	}

	err = server.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot run mqtt server")
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Println("app - Run - signal: " + s.String())
	case err := <-server.Notify():
		fmt.Println(err)
	}

	server.Close()
	log.Info().Msg("Stopped app")
}
