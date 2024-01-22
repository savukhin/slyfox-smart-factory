package service

import (
	"eventsproxy/internal/domain"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_proxyService_generateJWT(t *testing.T) {
	t.Parallel()

	svc, _, _ := setupProxyService(t)
	user := domain.GenerateUser()
	s, err := svc.generateJWT(user)
	require.NoError(t, err)

	userParsed, err := svc.parseJWT(s)
	require.NoError(t, err)
	if !reflect.DeepEqual(user, userParsed) {
		require.Failf(t, "Expected deep equal", "user1 %v\nuserParsed %v", user, userParsed)
	}
}
