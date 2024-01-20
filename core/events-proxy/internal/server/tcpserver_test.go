package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-mqtt/mqtt"
	"github.com/stretchr/testify/require"
)

func TestTcpServer_Notify(t *testing.T) {
	srv := NewTcpServer(TcpServerConfig{})
	err := srv.Notify()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println(<-err)
		wg.Done()
	}()
	srv.notify <- errors.New("my err")
	wg.Wait()
	t.Fail()
}

func TestTcpServer_Connect(t *testing.T) {
	srv := NewTcpServer(NewTcpServerConfig())
	err := srv.Run(context.Background())
	require.NoError(t, err)

	defer srv.Close()

	client, err := mqtt.VolatileSession("some-client", &mqtt.Config{
		UserName:     "username",
		Password:     []byte("password"),
		Dialer:       mqtt.NewDialer("tcp", "localhost:1883"),
		PauseTimeout: 4 * time.Second,
	})
	m, topic, err := client.ReadSlices()
	fmt.Println(m, topic, err)
}

func Test_calculateVarLenFromShift(t *testing.T) {
	packet := []byte{0x10, 0x24, 0x00, 0x04, 0x4d, 0x51, 0x54, 0x54, 0x05, 0xc2, 0x00, 0x3c, 0x05, 0x11, 0x00, 0x00, 0x07, 0x08, 0x00, 0x04, 0x6d, 0x79, 0x50, 0x79, 0x00, 0x06, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x00, 0x0a, 0x70, 0x61, 0x73, 0x73}
	len, end := calculateVarLenFromShift(packet, 1)
	require.EqualValues(t, len, 0x24)
	require.EqualValues(t, end, 1)

	username, password, err := extractCredentials(packet)
	require.NoError(t, err)
	require.EqualValues(t, username, "client")
	require.EqualValues(t, password, "pass")

	packet = []byte{16, 43, 0, 4, 77, 81, 84, 84, 4, 192, 0, 6, 0, 11, 115, 111, 109, 101, 45, 99, 108, 105, 101, 110, 116, 0, 8, 117, 115, 101, 114, 110, 97, 109, 101, 0, 8, 112, 97, 115, 115, 119, 111, 114, 100}
	username, password, err = extractCredentials(packet)
	require.NoError(t, err)
	require.EqualValues(t, username, "username")
	require.EqualValues(t, password, "password")
}
