package natsmq

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func SubscribeWithReply(nc *nats.Conn, subject string, h MsgHandler) error {
	_, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		response, err := h(msg)
		if err != nil {
			return
		}

		msg.Respond(response)
	})
	if err != nil {
		return err
	}

	nc.Flush()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	return nil
}

type MsgHandler func(msg *nats.Msg) ([]byte, error)

func PublichAndWaitForReply(nc *nats.Conn, subject string, data []byte, timeout time.Duration, replyCount int) ([]*nats.Msg, error) {
	replyTo := fmt.Sprintf("unique-reply-subject-%d", time.Now().UnixNano()) // TODO 确保唯一

	sub, err := nc.SubscribeSync(replyTo)
	if err != nil {
		return nil, err
	}
	nc.Flush()

	nc.PublishRequest(subject, replyTo, data)

	start := time.Now()
	msgs := make([]*nats.Msg, 0)
	i := 0
	for i < replyCount && time.Since(start) < timeout {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}
		i++
		msgs = append(msgs, msg)
	}
	sub.Unsubscribe()

	return msgs, nil
}

type Response struct {
	Data []byte
}

var once sync.Once
var singleton *nats.Conn

func GetSingletonNats(url string, opts ...nats.Option) (*nats.Conn, error) {
	var err error

	once.Do(func() {
		singleton, err = nats.Connect(url, opts...)
	})

	return singleton, err
}

func CloseSingletonNats() {
	if singleton != nil {
		singleton.Close()
	}
}
