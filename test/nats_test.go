package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

var natsAddr = "nats://me.io:4222"

var subject = "hello"

func TestNats(t *testing.T) {
	go sub("127.0.0.1")
	go sub("127.0.0.2")

	time.Sleep(time.Millisecond * 100)

	nc, err := nats.Connect(natsAddr, nats.Name("pub--------"))
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	for i := 0; i < 20; i++ {
		go func(index int) {
			pub(nc, fmt.Sprintf("ok-%d", index))
		}(i)
	}

	time.Sleep(time.Second * 3)
}

func pub(nc *nats.Conn, key string) {
	replyTo := "xxx" + key
	sub, err := nc.SubscribeSync(replyTo)
	if err != nil {
		log.Fatal(err)
	}
	nc.Flush()

	// Send the request
	nc.PublishRequest(subject, replyTo, []byte(key))

	// Wait for a single response
	max := 500 * time.Millisecond
	start := time.Now()
	responses := make([]string, 0)
	i := 0
	for i < 20 && time.Since(start) < max {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}
		i++
		responses = append(responses, string(msg.Data))
	}
	sub.Unsubscribe()

	log.Printf("[%s] received: %v\n", key, responses)
}

func sub(addr string) {
	nc, err := nats.Connect(natsAddr, nats.Name("sub--------"))
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		key := string(msg.Data)
		m := Message{addr, key}
		bs, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		msg.Respond(bs)
	})
	if err != nil {
		panic(err)
	}

	nc.Flush()

	c := make(chan os.Signal, 1)
	// 监听指定信号
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

type Message struct {
	Addr string
	Key  string
}
