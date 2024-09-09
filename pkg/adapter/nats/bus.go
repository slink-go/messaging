package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
)

type NatMessageBus struct {
	client Client
	logger logging.Logger
}

func (b *NatMessageBus) Publish(topic string, message api.Message, encoding ...api.Encoding) error {
	if encoding != nil && len(encoding) > 0 {
		enc, err := GetEncoder(encoding[0])
		if err != nil {
			return err
		}
		return b.publish(topic, message, enc)
	}
	return b.publish(topic, message, b.client.encoder)
}

func (b *NatMessageBus) publish(topic string, message api.Message, encoder Encoder) error {
	if err := validateClient(b.client); err != nil {
		return err
	}
	msg, err := encoder.Encode(topic, message)
	if err != nil {
		return err
	}
	return b.client.conn.PublishMsg(msg)
}

func (b *NatMessageBus) Subscribe(topic string, channel chan api.Message, options ...api.SubscribeOption) (func(), error) {
	if err := validateClient(b.client); err != nil {
		return nil, err
	}
	sub, err := b.client.conn.Subscribe(topic, func(msg *nats.Msg) {
		m, e := b.client.decode(msg)
		if e != nil {
			b.logger.Warning("could not decode incoming message: %v", e)
			return
		}
		if shouldAcceptMessage(m, options...) {
			if err := b.send(channel, m); err != nil {
				b.logger.Warning("could not send message: %v", err)
			}
		}
	})
	closer := func() {
		b.logger.Info("bus unsubscribe from %s", topic)
		_ = sub.Unsubscribe()
	}
	return closer, err
}

func (b *NatMessageBus) send(channel chan api.Message, message api.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("message send error: %v", r)
		}
	}()
	channel <- message
	return
}
