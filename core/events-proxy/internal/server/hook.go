package server

import (
	"bytes"
	"context"
	"encoding/json"
	"eventsproxy/internal/service"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var (
	totalOnConnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_proxy_total_on_connect",
		Help: "Total number of OnConnect messages",
	})

	totalOnConnectAuthenticate = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_proxy_total_on_connect_authenticate",
		Help: "Total number of OnConnectAuthenticate messages",
	})

	totalOnPublish = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_proxy_total_on_publish",
		Help: "Total number of OnPublish messages",
	})

	totalOnDisconnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_proxy_total_on_disconnect",
		Help: "Total number of OnDisconnect messages",
	})
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

func (hook *Hook) authCredentials(pk packets.Packet) (string, error) {
	username := pk.Connect.Username
	hashedPassword := pk.Connect.Password
	token, err := hook.svc.Auth(context.Background(), string(username), string(hashedPassword))
	if err != nil {
		return "", ErrUnauthorize
	}
	return token, nil
}

func (hook *Hook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	totalOnConnect.Inc()
	log.Info().Msg("Received OnConnect")
	_, err := hook.authCredentials(pk)
	if err != nil {
		log.Err(err).Msg("OnConnect error")
	}
	return err
}

func (hook *Hook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	totalOnConnectAuthenticate.Inc()
	log.Info().Msg("Received OnConnectAuthenticate")
	_, err := hook.authCredentials(pk)
	if err != nil {
		log.Info().Err(err).Msg("OnConnectAuthenticate Auth error")
	}
	return err != nil
}

type MessageToken struct {
	AuthToken string `json:"token"`
}

func (hook *Hook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	totalOnPublish.Inc()
	log.Info().Msg("Received OnPublish")
	var message MessageToken
	err := json.Unmarshal(pk.Payload, &message)
	if err != nil {
		log.Err(err).Msg("unmarshalling OnPublish json error")
		return pk, ErrUnauthorize
	}

	ctx := context.Background()
	_, err = hook.svc.VerifyToken(ctx, message.AuthToken)
	if err != nil {
		log.Err(err).Msg("verifying OnPublish token error")
		return pk, ErrUnauthorize
	}

	hook.svc.Publish(context.Background(), pk.TopicName, string(pk.Payload))
	return pk, nil
}

func (h *Hook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	totalOnDisconnect.Inc()
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
	return true
}
