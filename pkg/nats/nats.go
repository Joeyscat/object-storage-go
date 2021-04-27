package nats

import (
    "encoding/json"
    "github.com/nats-io/nats.go"
)

type Nats struct {
    nc *nats.Conn
}

func New(s string) *Nats {
    nc, err := nats.Connect("demo.nats.io", nats.Name("API PublishBytes Example"))
    if err != nil {
        panic(err)
    }

    return &Nats{nc}
}

func (n *Nats) Publish(sub string, body interface{}) {
    str, e := json.Marshal(body)
    if e != nil {
        panic(e)
    }

    if err := n.nc.Publish(sub, str); err != nil {
        panic(err)
    }
}

func (n *Nats) Consume() {

}

func (n *Nats) Close() {
    n.nc.Close()
}
