package server

import (
	"context"
	"errors"
	"eventsproxy/internal/service"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
)

var (
	ErrStopped     = errors.New("server is already stopped")
	ErrUnauthorize = errors.New("failed to auth")

	// ConnackMsg = []byte{0x20, 0}
	ConnackMsg = []byte{0x20, 2, 0, 0}
	// PubackMsg = []byte{0x40, 2, 0, 0}
	PubackMsg = []byte{0x40, 2, 0x80, 0x10}
)

type TcpServerConfig struct {
	Type    string
	Host    string
	Port    string
	Workers int
}

type TcpServer struct {
	cfg      TcpServerConfig
	listener net.Listener
	notify   chan error
	stop     chan struct{}
	proxySvc service.ProxyService
}

func NewTcpServerConfig() TcpServerConfig {
	return TcpServerConfig{
		Type:    "tcp",
		Host:    "localhost",
		Port:    "1883",
		Workers: 1,
	}
}

var Counter atomic.Int32

func NewTcpServer(cfg TcpServerConfig) TcpServer {
	return TcpServer{
		cfg:    cfg,
		notify: make(chan error),
	}
}

func (srv TcpServer) runRoutine(ctx context.Context, wg *sync.WaitGroup, worker int) (err error) {
runLoop:
	for {
		select {
		case <-srv.stop:
			break runLoop
		default:
			conn, err := srv.listener.Accept()
			c := Counter.Add(1)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			srv.handleConnection(ctx, conn, worker, c)
		}
	}
	wg.Done()
	return nil
}

func (srv TcpServer) handleConnection(ctx context.Context, conn net.Conn, worker int, counter int32) {
	var err error
	buffer := make([]byte, 1024)

	for err == nil {
		_, err = conn.Read(buffer)
		if err != nil {
			break
		}

		err = srv.handleRequest(ctx, conn, buffer, worker, counter)
	}

	conn.Close()
}

func createPuback(id byte) []byte {
	msg := PubackMsg
	msg[3] = byte(id)
	return msg
}

func (srv TcpServer) handleRequest(ctx context.Context, conn net.Conn, buffer []byte, worker int, counter int32) error {
	controlHeader := buffer[0]
	if isConnect(controlHeader) {
		username, hashedPassword, err := extractCredentials(buffer)
		err = srv.proxySvc.Auth(ctx, username, hashedPassword)
		if err != nil {
			return ErrUnauthorize
		}
		conn.Write(ConnackMsg)
	} else if isPublish(controlHeader) {
		fmt.Println("GOT PUBLISH")
		conn.Write(createPuback(byte(counter - 1)))
	} else {
		fmt.Println("unrecognized header", controlHeader)
	}
	return nil

}

func (srv TcpServer) Run(ctx context.Context) error {
	select {
	case <-srv.stop:
		return ErrStopped
	default:
	}

	Counter.Store(0)

	var err error
	srv.listener, err = net.Listen(srv.cfg.Type, srv.cfg.Host+":"+srv.cfg.Port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for i := 0; i < srv.cfg.Workers; i++ {
		wg.Add(1)
		go srv.runRoutine(ctx, &wg, i)
	}

	wg.Wait()
	return nil
}

func (srv TcpServer) Notify() chan error {
	return srv.notify
}

func (srv TcpServer) Close() {
	close(srv.stop)

}
