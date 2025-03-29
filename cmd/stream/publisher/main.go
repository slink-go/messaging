package main

import (
	"go.slink.ws/logging"
	"go.slink.ws/messaging/pkg/adapter/nats"
	"go.slink.ws/messaging/pkg/api"
	"os"
	"time"
)

const topic = "test-stream"

func main() {

	os.Setenv("GO_ENV", "dev")

	c := nats.NewNatsClient(api.EncodingJson)
	c.Connect()
	defer c.Close()

	var streamConfig = api.StreamConfig{
		Name:        "test-stream",
		Description: "Test Stream",
		Subjects:    []string{"subject-1", "subject-2", "subject-3"},
	}

	s := nats.NewMessageStream(c, streamConfig)
	//s.Delete()
	//s.Create()
	s.Flush()
	for {
		time.Sleep(2 * time.Second)
		publish(s, "subject-1")
		time.Sleep(2 * time.Second)
		publish(s, "subject-2")
		time.Sleep(2 * time.Second)
		publish(s, "subject-3")
	}

}

func publish(s api.MessageStream, subject string) {
	m := api.NewBasicMessage()
	if e := s.Publish(subject, m); e != nil {
		logging.GetLogger("stream publisher").Warning("publish error: %v", e)
	} else {
		logging.GetLogger("stream publisher").Info("new message published: %v", m)
	}
}
