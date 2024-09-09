package main

import (
	"fmt"
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/adapter/nats"
	"github.com/slink-go/messaging/pkg/api"
	"math/rand"
	"os"
	"time"
)

var topics = []string{"q1", "system-bus", "event-bus"}

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

var idx int

func publish(b api.MessageBus) {
	var m api.Message
	if rand.Intn(100) > 50 {
		m = api.NewBasicMessage()
	} else {
		idx += 1
		m = api.NewTextMessage(fmt.Sprintf("message #%d", idx))
	}
	for _, topic := range topics {
		if e := b.Publish(topic, m); e != nil {
			logging.GetLogger("bus publisher").Warning("publish error: %v", e)
		} else {
			logging.GetLogger("bus publisher").Info("new message published: %v", m)
		}
	}
}
