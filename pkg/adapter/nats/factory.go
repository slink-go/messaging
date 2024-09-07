package nats

import (
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
)

func NewMessageBus(client Client) api.MessageBus {
	return &NatMessageBus{
		client: client,
		logger: logging.GetLogger("nats-message-bus"),
	}
}
func NewMessageStream(client Client, config api.StreamConfig) api.MessageStream {

	var defaultName = "default"
	if config.Name == "" {
		config.Name = defaultName
	}

	var defaultSubjects = []string{defaultName}
	if config.Subjects == nil || len(config.Subjects) == 0 {
		config.Subjects = defaultSubjects
	}

	return &NatMessageStream{
		client: client,
		config: config,
		logger: logging.GetLogger("nats-message-stream"),
	}

}
