package nats

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/messaging/pkg/api"
	"github.com/vmihailenco/msgpack/v5"
)

type Encoder struct {
}

func (e *Encoder) Encode(topic string, message api.Message, encoding api.Encoding) (*nats.Msg, error) {

	msg := e.createMessage(topic, encoding, message)

	switch encoding {

	case api.EncodingJson:
		msg.Data = []byte(fmt.Sprintf("%s", message))
		return msg, nil

	case api.EncodingMsgPack:
		data, err := e.msgPackEncode(message)
		if err != nil {
			return nil, err
		}
		msg.Data = data
		return msg, nil

	case api.EncodingGob:
		data, err := e.gobEncode(message)
		if err != nil {
			return nil, err
		}
		msg.Data = data
		return msg, nil

	}

	return nil, fmt.Errorf("unsupported encoding: %s", encoding)
}
func (e *Encoder) createMessage(topic string, encoding api.Encoding, message api.Message) *nats.Msg {
	msg := nats.NewMsg(topic)
	msg.Header.Add(msgEncodingHeader, encoding.String())
	msg.Header.Add(msgTypeHeader, fmt.Sprintf("%s", message.GetKind()))
	return msg
}
func (e *Encoder) msgPackEncode(message interface{}) ([]byte, error) {
	return msgpack.Marshal(message)
}
func (e *Encoder) gobEncode(message interface{}) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(message)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
