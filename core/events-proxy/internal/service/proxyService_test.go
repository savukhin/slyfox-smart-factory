package service

import (
	"context"
	"eventsproxy/internal/config"
	"eventsproxy/internal/domain"
	mock_producer "eventsproxy/internal/service/producer/mocks"
	mock_repo "eventsproxy/internal/service/repo/mocks"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func setupProxyService(t *testing.T) (proxyService, mock_repo.MockUserRepo, mock_producer.MockNatsProducer) {
	ctrl := gomock.NewController(t)
	mockRepo := mock_repo.NewMockUserRepo(ctrl)
	mockProducer := mock_producer.NewMockNatsProducer(ctrl)
	cfg := config.JwtConfig{
		Secret:      []byte("somesecret"),
		DurationMin: 10,
	}
	svc := NewProxyService(mockRepo, mockProducer, cfg)

	return svc, *mockRepo, *mockProducer
}

func Test_proxyService_Auth(t *testing.T) {
	t.Parallel()

	svc, mockRepo, _ := setupProxyService(t)
	username := "username1"
	aesKey := strings.Repeat("a", 32)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "alice@example.com",
	})
	require.NoError(t, err)
	totpKey := key.Secret()

	code := generateTOTPCode(totpKey, time.Now().Unix())
	src := fmt.Sprintf("%d", code)

	dst, err := encryptMessage([]byte(aesKey), src)
	require.NoError(t, err)

	user := domain.User{
		Username: username,
		AesKey:   aesKey,
		TotpKey:  totpKey,
	}
	mockRepo.EXPECT().GetByUsername(gomock.Any(), username).Return(user, nil)

	_, err = svc.Auth(context.Background(), username, string(dst))
	require.NoError(t, err)

}
