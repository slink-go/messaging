package api

type MessageBus interface {
	Publish(topic string, message Message) error
	PublishEncoded(topic string, message Message, encoding Encoding) error
	Subscribe(topic string, channel chan Message, options ...SubscribeOption) error
	//SubscribeWithFilter(topic string, channel chan Message, filterFunc func(message Message) bool) error
}
