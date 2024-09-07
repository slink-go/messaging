package nats

import "errors"

var ErrClientIsNil = errors.New("nats client is nil")
var ErrNatsConnectionIsNil = errors.New("nats connection is nil")
var ErrNatsNotConnected = errors.New("nats is not connected")
