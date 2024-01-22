package service

import (
	"context"
	"errors"
	"eventsproxy/internal/config"
	"eventsproxy/internal/domain"
	"eventsproxy/internal/service/producer"
	"eventsproxy/internal/service/repo"
	"fmt"
	"time"

	twoqueue "github.com/floatdrop/2q"

	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"
)

var (
	otpLength = 6

	ErrPasswordInvalid = errors.New("password in invalid")
)

type ProxyService interface {
	Auth(ctx context.Context, username, hashedPassword string) (string, error)
	VerifyToken(ctx context.Context, token string) (domain.User, error)
	Publish(ctx context.Context, topic, message string) error
}

type proxyService struct {
	userRepo     repo.UserRepo
	natsProducer producer.NatsProducer
	jwtSecret    []byte
	jwtDuration  time.Duration
	tokenCache   *twoqueue.TwoQueue[string, domain.User]
}

func NewProxyService(userRepo repo.UserRepo, natsProducer producer.NatsProducer, cfg config.JwtConfig) proxyService {
	return proxyService{
		userRepo:    userRepo,
		tokenCache:  twoqueue.New[string, domain.User](300),
		jwtSecret:   cfg.Secret,
		jwtDuration: time.Duration(cfg.DurationMin) * time.Minute,
	}
}

func (svc *proxyService) Auth(ctx context.Context, username, hashedPassword string) (string, error) {
	log.Info().Str("username", username).Str("hashedPassword", hashedPassword).Msg("Service Auth")
	user, err := svc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("error calling repo: %v", err)
	}

	password, err := decryptMessage([]byte(user.AesKey), hashedPassword)
	if err != nil {
		return "", fmt.Errorf("error decrypting AES: %v", err)
	}
	fmt.Println("password", password)

	isValid := totp.Validate(password, user.TotpKey)
	if !isValid {
		return "", fmt.Errorf("error validating OTP: %v", ErrPasswordInvalid)
	}

	token, err := svc.generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("error generation JWT: %v", err)
	}
	return token, nil
}

func (svc *proxyService) VerifyToken(ctx context.Context, token string) (domain.User, error) {
	user, err := svc.parseJWT(token)
	if err != nil {
		return user, fmt.Errorf("parsing token: %v", err)
	}
	return user, nil
}

func (svc *proxyService) Publish(ctx context.Context, topic, message string) error {
	return svc.natsProducer.Publish(ctx, topic, message)
}
