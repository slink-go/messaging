package main

import (
	"github.com/slink-go/logging"
	adapter "github.com/slink-go/messaging/pkg/adapter/nats"
	"github.com/slink-go/messaging/pkg/api"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const topic = "q1"

func main() {

	os.Setenv("GO_ENV", "dev")
	c := adapter.NewNatsClient()

	c.Connect()
	defer c.Close()

	b := adapter.NewMessageBus(c)

	messages1 := make(chan api.Message)
	messages2 := make(chan api.Message)

	closer1, err := b.Subscribe(topic, messages1)
	if err != nil {
		print(err)
	}
	_, err = b.Subscribe(topic, messages2)
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
