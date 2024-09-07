package nats

import (
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
)

func NewMessageBus(client Client, opts ...interface{}) api.MessageBus {
	return &NatMessageBus{
		client:   client,
		encoding: getDefaultEncoding(opts...),
		logger:   logging.GetLogger("nats-message-bus"),
	}
}
func NewMessageStream(client Client, opts ...interface{}) api.MessageStream {
	return &NatMessageStream{
		client:   client,
		encoding: getDefaultEncoding(opts...),
		logger:   logging.GetLogger("nats-message-stream"),
	}
}

func getDefaultEncoding(opts ...interface{}) api.Encoding {
	encoding := api.EncodingJson
	if opts != nil || len(opts) > 0 {
		v, ok := opts[0].(api.Encoding)
		if ok {
			encoding = v
		}
	}
	return encoding
}
