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

	c := nats.NewNatsClient(api.EncodingGob)
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
	n := rand.Intn(100)
	if n > 66 {
		m = api.NewBasicMessage()
	} else if n > 33 {
		idx += 1
		obj := make(map[string]any)
		obj["key"] = "item"
		obj["value"] = fmt.Sprintf("value #%d", idx)
		obj["index"] = idx
		m = api.NewObjectMessage(obj)
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
