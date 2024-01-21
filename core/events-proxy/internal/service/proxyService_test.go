package service

import (
	"context"
	mock_producer "eventsproxy/internal/service/producer/mocks"
	"eventsproxy/internal/service/repo"
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
	svc := NewProxyService(mockRepo, mockProducer)

	return svc, *mockRepo, *mockProducer
}

func Test_proxyService_Auth(t *testing.T) {
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

	user := repo.UserRecord{
		Username: username,
		AesKey:   aesKey,
		TotpKey:  totpKey,
	}
	mockRepo.EXPECT().GetByUsername(gomock.Any(), username).Return(user, nil)

	err = svc.Auth(context.Background(), username, string(dst))
	require.NoError(t, err)
}
