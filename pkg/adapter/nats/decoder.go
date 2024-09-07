package nats

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/messaging/pkg/api"
)

type Decoder struct {
}

func (d *Decoder) Decode(msg *nats.Msg, decoders MessageHandlers) (api.Message, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	encoding := api.ParseEncoding(msg.Header.Get(msgEncodingHeader))
	if encoding == api.EncodingUndefined {
		encoding = api.EncodingJson
	}

	t := msg.Header.Get(msgTypeHeader)
	if t == "" {
		return nil, errors.New("message type is missing")
	}

	switch encoding {
	default:
		key := encoding.String() + ":" + t
		dec, ok := decoders[key]
		if !ok {
			return nil, fmt.Errorf("no decoder specified for message type %s", t)
		}
		return dec(msg)
	}
}
