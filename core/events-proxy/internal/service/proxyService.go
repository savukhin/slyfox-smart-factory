package service

import (
	"context"
	"crypto/aes"
	"errors"
	"eventsproxy/internal/service/producer"
	"eventsproxy/internal/service/repo"

	"github.com/pquerna/otp/totp"
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
	user, err := svc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return err
	}

	c, err := aes.NewCipher([]byte(user.AesKey))
	if err != nil {
		return err
	}

	decrypted := make([]byte, len(hashedPassword))
	c.Decrypt(decrypted, []byte(hashedPassword))
	password := string(decrypted[:otpLength])

	isValid := totp.Validate(password, user.TotpKey)
	if !isValid {
		return ErrPasswordInvalid
	}

	return nil
}

func (svc *proxyService) Publish(ctx context.Context, topic, message string) error {
	return svc.natsProducer.Publish(ctx, topic, message)
}
