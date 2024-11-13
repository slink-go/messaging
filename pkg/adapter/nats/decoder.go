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

	t := msg.Header.Get(msgTypeHeader)
	if t == "" {
		return nil, fmt.Errorf("message type is missing: %v", string(msg.Data))
	}

	dec, err := d.getDecoder(encoding, t, decoders)
	if err != nil {
		return nil, err
	}
	return dec(msg)
}

func (d *Decoder) getDecoder(encoding api.Encoding, msgType string, decoders MessageHandlers) (MessageDecoderFunc, error) {
	dec, ok := decoders[encoding.String()+":"+msgType]
	if ok {
		return dec, nil
	}
	dec, ok = decoders[encoding.String()+":*"]
	if ok {
		return dec, nil
	}
	return nil, fmt.Errorf("no '%s' decoder specified for message type '%s'", encoding, msgType)
}
