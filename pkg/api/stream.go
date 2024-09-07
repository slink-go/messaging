package api

type StreamConfig struct {
	Name        string
	Description string
	Subjects    []string
	// TODO: add additional parameters: size, ttl, replicas, etc
}

type MessageStream interface {
	Create() error
	Delete() error
	Flush() error

	Publish(topic string, message Message, encoding ...Encoding) error
	Subscribe(topic string, channel chan Message, options ...SubscribeOption) error
}
