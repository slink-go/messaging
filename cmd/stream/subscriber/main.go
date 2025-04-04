package main

import (
	"bytes"
	"encoding/gob"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"go.slink.ws/logging"
	adapter "go.slink.ws/messaging/pkg/adapter/nats"
	"go.slink.ws/messaging/pkg/api"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const subject = "subject-1"

func gobBasicMessageDecoder(msg *nats.Msg) (api.Message, error) {
	r := bytes.NewReader(msg.Data)
	var m api.BasicMessage
	dec := gob.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func msgPackBasicMessageDecoder(msg *nats.Msg) (api.Message, error) {
	var v api.BasicMessage
	err := msgpack.Unmarshal(msg.Data, &v)
	return &v, err
}
func jsonBasicMessageDecoder(msg *nats.Msg) (api.Message, error) {
	var bm api.BasicMessage
	err := json.Unmarshal(msg.Data, &bm)
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

func main() {

	os.Setenv("GO_ENV", "dev")
	c := adapter.NewNatsClient(api.EncodingString)

	//c.AddMessageDecoder(adapter.MessageDecoder{
	//	MessageType: "BasicMessage",
	//	Encoding:    api.EncodingJson,
	//	Handler:     jsonBasicMessageDecoder,
	//})
	//c.AddMessageDecoder(adapter.MessageDecoder{
	//	MessageType: "BasicMessage",
	//	Encoding:    api.EncodingMsgPack,
	//	Handler:     msgPackBasicMessageDecoder,
	//})
	//c.AddMessageDecoder(adapter.MessageDecoder{
	//	MessageType: "BasicMessage",
	//	Encoding:    api.EncodingGob,
	//	Handler:     gobBasicMessageDecoder,
	//})

	c.Connect()
	defer c.Close()

	var streamConfig = api.StreamConfig{
		Name:        "test-stream",
		Description: "Test Stream",
		Subjects:    []string{"subject-1", "subject-2", "subject-3"},
	}

	s := adapter.NewMessageStream(c, streamConfig)

	messages := make(chan api.Message)

	s.Subscribe(subject, messages)

	go func() {
		for msg := range messages {
			logging.GetLogger("subscriber").Info("received message: %s", msg)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	sigChn := make(chan os.Signal, 1)
	signal.Notify(sigChn,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		_ = <-sigChn
		wg.Done()
	}()

	wg.Wait()

}
