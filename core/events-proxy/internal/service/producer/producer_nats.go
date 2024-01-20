package producer

import (
	"context"

	"github.com/nats-io/nats.go"
)

type NatsProducer interface {
	Publish(ctx context.Context, topic, message string) error
}

type natsProducer struct {
	js nats.JetStream
}

func NewNastProducer(js nats.JetStream) natsProducer {
	return natsProducer{
		js: js,
	}
}

func (producer *natsProducer) Publish(ctx context.Context, topic, message string) error {
	_, err := producer.js.Publish(topic, []byte(message))
	return err
}
