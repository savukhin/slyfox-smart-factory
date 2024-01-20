package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-mqtt/mqtt"
)

// Publish is a method from mqtt.Client.
var Publish func(quit <-chan struct{}, message []byte, topic string) error

var wg sync.WaitGroup

func mock() {
	// defer func() {
	// 	r := recover()
	// 	fmt.Println("recovered with ", r)
	// }()
	client, err := mqtt.VolatileSession("some-client", &mqtt.Config{
		UserName:     "username1",
		Password:     []byte("password2"),
		Dialer:       mqtt.NewDialer("tcp", "localhost:1883"),
		KeepAlive:    6,
		PauseTimeout: 4 * time.Second,
	})

	fmt.Println(client, err)
	defer client.Close()

	// err = client.Publish(q, []byte("msg"), "topic")
	// m, topic, err := client.ReadSlices()
	// fmt.Println(m, topic, err)
	go func() {
		var big *mqtt.BigMessage
		for {
			message, topic, err := client.ReadSlices()
			switch {
			case err == nil:
				// do something with inbound message
				log.Printf("📥 %q: %q", topic, message)

			case errors.As(err, &big):
				log.Printf("📥 %q: %d byte message omitted", big.Topic, big.Size)

			case errors.Is(err, mqtt.ErrClosed):
				log.Print(err)
				return // terminated

			case mqtt.IsConnectionRefused(err):
				log.Print(err) // explains rejection
				// mqtt.ErrDown for a while
				time.Sleep(15 * time.Minute)

			default:
				log.Print("broker unavailable: ", err)
				// mqtt.ErrDown during backoff
				time.Sleep(2 * time.Second)
			}
		}
	}()

	quit := make(<-chan struct{})
	for i := 0; i < 2; i++ {
		topic := "sometopic" + string(i)
		for i, c := range topic {
			if c == 0 {
				fmt.Printf("i(%v) symb is 0\n", i)
			}
		}
		err = client.Publish(quit, []byte("some message "+string(i)), topic[:len(topic)-1])
		fmt.Println("publish err", err)
		time.Sleep(200 * time.Millisecond)
	}

	err = client.Disconnect(quit)
	fmt.Println("dicsonnect", err)

	wg.Done()
}

func main() {
	parallels := 1
	for i := 0; i < parallels; i++ {
		wg.Add(1)
		go mock()
	}
	wg.Wait()
	fmt.Println("OK")
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()
	// // launch read-routine
	// go func() {
	// 	var big *mqtt.BigMessage
	// 	for {
	// 		message, topic, err := client.ReadSlices()
	// 		switch {
	// 		case err == nil:
	// 			// do something with inbound message
	// 			log.Printf("📥 %q: %q", topic, message)

	// 		case errors.As(err, &big):
	// 			log.Printf("📥 %q: %d byte message omitted", big.Topic, big.Size)

	// 		case errors.Is(err, mqtt.ErrClosed):
	// 			log.Print(err)
	// 			return // terminated

	// 		case mqtt.IsConnectionRefused(err):
	// 			log.Print(err) // explains rejection
	// 			// mqtt.ErrDown for a while
	// 			time.Sleep(15 * time.Minute)

	// 		default:
	// 			log.Print("broker unavailable: ", err)
	// 			// mqtt.ErrDown during backoff
	// 			time.Sleep(2 * time.Second)
	// 		}
	// 	}
	// }()

	// // Install each method in use as a package variable. Such setup is
	// // compatible with the tools proveded from the mqtttest subpackage.
	// Publish = client.Publish
}
