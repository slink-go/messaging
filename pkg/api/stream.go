package api

type MessageStream interface {
	Create(topic string) error
	Delete(topic string) error
	Flush(topic string) error

	Publish(topic string, message Message) error
	Subscribe(topic string, channel chan Message, options ...SubscribeOption) error
	//SubscribeWithFilter(topic string, channel chan interface{}) error
}
