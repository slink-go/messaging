package nats

import (
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
)

type NatMessageStream struct {
	client   Client
	encoding api.Encoding
	logger   logging.Logger
}

func (b *NatMessageStream) Create(topic string) error {
	return nil
}
func (b *NatMessageStream) Delete(topic string) error {
	return nil
}
func (b *NatMessageStream) Flush(topic string) error {
	return nil
}

func (b *NatMessageStream) Publish(topic string, message api.Message) error {
	return nil
}
func (b *NatMessageStream) Subscribe(topic string, channel chan api.Message, options ...api.SubscribeOption) error {
	return nil
}
