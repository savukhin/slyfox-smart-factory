package service

import (
	"context"
	"errors"
	"eventsproxy/internal/service/producer"
	"eventsproxy/internal/service/repo"
	"fmt"

	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
)

var (
	otpLength = 6

	ErrPasswordInvalid = errors.New("password in invalid")
)

type ProxyService interface {
	Auth(ctx context.Context, username, hashedPassword string) error
	Publish(ctx context.Context, topic, message string) error
}

type proxyService struct {
	userRepo     repo.UserRepo
	natsProducer producer.NatsProducer
}

func NewProxyService(userRepo repo.UserRepo, natsProducer producer.NatsProducer) proxyService {
	return proxyService{
		userRepo: userRepo,
	}
}

func (svc *proxyService) Auth(ctx context.Context, username, hashedPassword string) error {
	log.Info().Str("username", username).Str("hashedPassword", hashedPassword).Msg("Service Auth")
	user, err := svc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	password, err := decryptMessage([]byte(user.AesKey), hashedPassword)
	if err != nil {
		return err
	}
	fmt.Println("password", password)

	isValid := totp.Validate(password, user.TotpKey)
	if !isValid {
		return ErrPasswordInvalid
	}

	return nil
}

func (svc *proxyService) Publish(ctx context.Context, topic, message string) error {
	return svc.natsProducer.Publish(ctx, topic, message)
}
