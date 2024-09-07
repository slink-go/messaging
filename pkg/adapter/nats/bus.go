package nats

import (
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
)

type NatMessageBus struct {
	client   Client
	encoding api.Encoding
	logger   logging.Logger
}

func (b *NatMessageBus) Publish(topic string, message api.Message) error {
	return b.PublishEncoded(topic, message, b.encoding)
}

func (b *NatMessageBus) PublishEncoded(topic string, message api.Message, encoding api.Encoding) error {
	if b.client.conn == nil {
		return errors.New("nats connection is nil")
	}
	msg, err := b.client.encode(topic, message, encoding)
	if err != nil {
		return err
	}
	return b.client.conn.PublishMsg(msg)
}

func (b *NatMessageBus) Subscribe(topic string, channel chan api.Message, options ...api.SubscribeOption) error {
	b.client.conn.Subscribe(topic, func(msg *nats.Msg) {
		m, e := b.client.decode(msg)
		if e != nil {
			b.logger.Warning("could not decode incoming message: %v", e)
			return
		}
		channel <- m
	})
	return nil
}
