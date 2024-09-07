package nats

import "github.com/slink-go/messaging/pkg/api"

func validateClient(client Client) error {
	if client.conn == nil {
		return ErrNatsConnectionIsNil
	}
	if !client.conn.IsConnected() {
		return ErrNatsNotConnected
	}
	return nil
}
func shouldAcceptMessage(message api.Message, options ...api.SubscribeOption) bool {
	accept := true
	if options != nil && len(options) > 0 {
		accept := false
		for _, opt := range options {
			accept = accept || opt.Apply(message)
		}
	}
	return accept
}
