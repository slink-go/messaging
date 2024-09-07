package api

type SubscribeOption interface {

	// TODO: implement filters via options pattern

	Apply(message Message) bool
}
