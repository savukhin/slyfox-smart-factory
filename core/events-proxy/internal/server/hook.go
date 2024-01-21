package server

import (
	"bytes"
	"context"
	"eventsproxy/internal/service"
	"fmt"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/rs/zerolog/log"
)

type Hook struct {
	mqtt.HookBase
	svc service.ProxyService
}

func newHook(svc service.ProxyService) *Hook {
	return &Hook{
		svc: svc,
	}
}

// ID returns the ID of the hook.
func (hook *Hook) ID() string {
	return "slyfox-hook"
}

// Provides indicates which hook methods this hook provides.
func (hook *Hook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnConnectAuthenticate,
		mqtt.OnSubscribe,
		mqtt.OnUnsubscribe,
		mqtt.OnPacketRead,
		mqtt.OnPublish,
		mqtt.OnDisconnect,
		mqtt.OnACLCheck,
	}, []byte{b})
}

func (hook *Hook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	log.Info().Msg("Received OnConnect")
	return nil
}

func (hook *Hook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	log.Info().Msg("Received OnConnectAuthenticate")
	fmt.Printf("OnConnectAuthenticate cl = %v pk = %v\n", cl, pk)
	username := pk.Connect.Username
	hashedPassword := pk.Connect.Password
	err := hook.svc.Auth(context.Background(), string(username), string(hashedPassword))
	// fmt.Printf("Auth result cl = %v pk = %v\n", cl, pk)
	log.Info().Err(err).Msg("Auth result")
	return err != nil
}

func (hook *Hook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	log.Info().Msg("Received OnPublish")
	fmt.Printf("OnPublish cl = %v pk = %v\n", cl, pk)
	hook.svc.Publish(context.Background(), pk.TopicName, string(pk.Payload))
	return pk, nil
}

func (h *Hook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
}

// OnPacketRead is called when a packet is received.
func (h *Hook) OnPacketRead(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	return pk, nil
}

func (h *Hook) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnUnsubscribe is called when a client unsubscribes from one or more filters.
func (h *Hook) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnACLCheck returns true/allowed for all checks.
func (h *Hook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	log.Info().Msg("Received OnACLCheck")
	return true
}
