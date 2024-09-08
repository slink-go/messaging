package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/logging"
	adapter "github.com/slink-go/messaging/pkg/adapter/nats"
	"github.com/slink-go/messaging/pkg/api"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const topic = "q1"

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

	c.AddMessageDecoder(adapter.MessageDecoder{
		MessageType: "BasicMessage",
		Encoding:    api.EncodingJson,
		Handler:     jsonBasicMessageDecoder,
	})
	c.AddMessageDecoder(adapter.MessageDecoder{
		MessageType: "BasicMessage",
		Encoding:    api.EncodingMsgPack,
		Handler:     msgPackBasicMessageDecoder,
	})
	c.AddMessageDecoder(adapter.MessageDecoder{
		MessageType: "BasicMessage",
		Encoding:    api.EncodingGob,
		Handler:     gobBasicMessageDecoder,
	})

	c.Connect()
	defer c.Close()

	b := adapter.NewMessageBus(c)

	messages1 := make(chan api.Message)
	messages2 := make(chan api.Message)

	closer1, err := b.Subscribe(topic, messages1)
	if err != nil {
		print(err)
	}
	closer2, err := b.Subscribe(topic, messages2)
	if err != nil {
		print(err)
	}

	go func() {
		for msg := range messages1 {
			logging.GetLogger("subscriber-1").Info("received message: %s", msg)
		}
	}()

	go func() {
		for msg := range messages2 {
			logging.GetLogger("subscriber-2").Info("received message: %s", msg)
		}
	}()

	time.Sleep(7 * time.Second)
	closer1()
	close(messages1)

	time.Sleep(7 * time.Second)
	closer2()
	close(messages2)

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
