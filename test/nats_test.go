package test

import (
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

var natsAddr = "nats://192.168.50.186:4222"

func TestNats(t *testing.T) {
	sub("127.0.0.1")

	time.Sleep(time.Millisecond * 100)

	pub("ok")

	time.Sleep(time.Second * 1)
}

func pub(key string) {
	nc, err := nats.Connect(natsAddr, nats.Name("pub--------"))
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	// nc.Publish("hello", []byte("xx"))
	r, err := nc.Request("hello", []byte("xx"), time.Second*1)
	if err != nil {
		panic(err)
	}
	log.Printf("request return: %v\n", string(r.Data))
}

func sub(addr string) {
	nc, err := nats.Connect(natsAddr, nats.Name("sub--------"))
	if err != nil {
		panic(err)
	}
	defer nc.Close()
	nc.Subscribe("hello", func(msg *nats.Msg) {
		key, err := strconv.Unquote(string(msg.Data))
		if err != nil {
			panic(err)
		}
		m := Message{addr, key}
		log.Printf("msg: %+v\n", m)
		bs, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		msg.Respond(bs)
	})
}

type Message struct {
	Addr string
	Key  string
}
