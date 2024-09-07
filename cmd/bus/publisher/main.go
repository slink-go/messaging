package main

import (
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/adapter/nats"
	"github.com/slink-go/messaging/pkg/api"
	"os"
	"time"
)

const topic = "q1"

func main() {

	os.Setenv("GO_ENV", "dev")

	c := nats.NewNatsClient(api.EncodingJson)
	c.Connect()
	defer c.Close()

	b := nats.NewMessageBus(c)
	for {
		time.Sleep(2 * time.Second)
		publish(b)
	}

}

func publish(b api.MessageBus) {
	m := api.NewBasicMessage()
	if e := b.Publish(topic, m); e != nil {
		logging.GetLogger("bus publisher").Warning("publish error: %v", e)
	} else {
		logging.GetLogger("bus publisher").Info("new message published: %v", m)
	}
}
