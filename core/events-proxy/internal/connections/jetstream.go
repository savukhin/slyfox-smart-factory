package connections

import (
	"eventsproxy/internal/config"

	"github.com/nats-io/nats.go"
)

func CreateNatsJetstream(cfg config.NatsConfig) (*nats.Conn, nats.JetStream, error) {
	nc, err := nats.Connect(cfg.Urls)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	return nc, js, err
}
