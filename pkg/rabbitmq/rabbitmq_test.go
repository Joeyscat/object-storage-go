package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

const host = "amqp://jojo:jojo@192.168.50.186:5672"

func TestPublish(t *testing.T) {
	q := New(host)
	defer q.Close()
	q.Bind("test")

	q2 := New(host)
	defer q2.Close()
	q2.Bind("test")

	q3 := New(host)
	defer q3.Close()

	expect := "test"
	q3.Publish("test2", "any")
	q3.Publish("test", expect)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
	if msg.ReplyTo != q3.Name {
		t.Error(msg)
	}

	c2 := q2.Consume()
	msg = <-c2
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}
	if msg.ReplyTo != q3.Name {
		t.Error(msg)
	}
	q2.Send(msg.ReplyTo, "test3")
	c3 := q3.Consume()
	msg = <-c3
	if string(msg.Body) != `"test3"` {
		t.Error(string(msg.Body))
	}
}

func TestSend(t *testing.T) {
	q := New(host)
	defer q.Close()

	q2 := New(host)
	defer q2.Close()

	expect := "test"
	expect2 := "test2"
	q2.Send(q.Name, expect)
	q2.Send(q2.Name, expect2)

	c := q.Consume()
	msg := <-c
	var actual interface{}
	err := json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect {
		t.Errorf("expected %s, actual %s", expect, actual)
	}

	c2 := q2.Consume()
	msg = <-c2
	err = json.Unmarshal(msg.Body, &actual)
	if err != nil {
		t.Error(err)
	}
	if actual != expect2 {
		t.Errorf("expected %s, actual %s", expect2, actual)
	}
}

func TestA(t *testing.T) {
	go sub("haha")

	for i := 0; i < 100; i++ {
		go func(index int) {
			pub(fmt.Sprintf("%d", index))
		}(i)
	}

	// for j := 0; j < 2; j++ {
	// 	go func(index int) {
	// 		sub(fmt.Sprintf("%d", index))
	// 	}(j)
	// }

	time.Sleep(5 * time.Second)
}

func pub(key string) {
	q := New(host)
	q.Publish("test002", key)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()

	locateInfo := make(map[string]string)
	for i := 0; i < 1; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info Message
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Key] = info.Addr
	}
	log.Printf("key: %s, localInfo: %+v\n", key, locateInfo)
}

func sub(addr string) {
	q := New(host)
	defer q.Close()
	q.Bind("test002")
	c := q.Consume()

	for msg := range c {
		key, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		log.Printf("key: %s", key)
		q.Send(msg.ReplyTo, Message{
			Addr: addr,
			Key:  key,
		})
	}
}

type Message struct {
	Addr string
	Key  string
}
