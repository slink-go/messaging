package main

import (
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/adapter/nats"
	"github.com/slink-go/messaging/pkg/api"
	"os"
	"time"
)

func main() {

	os.Setenv("GO_ENV", "dev")

	c := nats.NewNatsClient()
	c.Connect()
	defer c.Close()

	b := nats.NewMessageBus(c, api.EncodingMsgPack)
	for {
		time.Sleep(2 * time.Second)
		publish(b)
	}

}

func publishEncoded(b api.MessageBus, encoding api.Encoding) {
	m := api.NewBasicMessage()
	if e := b.PublishEncoded("test", m, encoding); e != nil {
		logging.GetLogger("publisher").Warning("publish error: %v", e)
	} else {
		logging.GetLogger("publisher").Info("new message published (as %s): %v", encoding, m)
	}
}
func publish(b api.MessageBus) {
	m := api.NewBasicMessage()
	if e := b.Publish("test", m); e != nil {
		logging.GetLogger("publisher").Warning("publish error: %v", e)
	} else {
		logging.GetLogger("publisher").Info("new message published: %v", m)
	}
}
