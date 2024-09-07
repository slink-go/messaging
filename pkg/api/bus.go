package api

type MessageBus interface {

	// По умолчанию используется defaultEncoder (заданный при создании nats-клиента),
	// но при необходимости мы можем указать другую кодировку при публикации сообщения

	Publish(topic string, message Message, encoding ...Encoding) error

	Subscribe(topic string, channel chan Message, options ...SubscribeOption) error
}
