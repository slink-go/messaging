package nats

// https://canopas.com/how-to-do-data-serialization-and-deserialization-in-golang-90c80b0b6d58

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/slink-go/messaging/pkg/api"
	"github.com/vmihailenco/msgpack/v5"
)

type Encoder interface {
	Encode(topic string, message api.Message) (*nats.Msg, error)
}

// region - factory

var encoders map[api.Encoding]Encoder

func GetEncoder(encoding api.Encoding) (Encoder, error) {
	v, ok := encoders[encoding]
	if !ok {
		enc, err := createEncoder(encoding)
		if err != nil {
			return nil, err
		}
		if encoders == nil {
			encoders = make(map[api.Encoding]Encoder)
		}
		encoders[encoding] = enc
		return enc, nil
	}
	return v, nil
}

func createEncoder(encoding api.Encoding) (Encoder, error) {
	switch encoding {
	case api.EncodingString:
		return &StringEncoder{}, nil
	case api.EncodingJson:
		return &JsonEncoder{}, nil
	case api.EncodingMsgPack:
		return &MsgPackEncoder{}, nil
	case api.EncodingGob:
		return &GobEncoder{}, nil
	default:
		return nil, fmt.Errorf("could not find suitable encoder for encoding %s (%s)", encoding, encoding)
	}
}

//endregion
// region - string encoder

type StringEncoder struct{}

func (j *StringEncoder) Encode(topic string, message api.Message) (*nats.Msg, error) {
	msg := createMessage(topic, api.EncodingString, message)
	msg.Data = []byte(fmt.Sprintf("%s", message))
	return msg, nil
}

// endregion
// region - json encoder

type JsonEncoder struct{}

func (j *JsonEncoder) Encode(topic string, message api.Message) (*nats.Msg, error) {
	msg := createMessage(topic, api.EncodingJson, message)
	b, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	msg.Data = b
	return msg, nil
}

// endregion
// region - message pack encoder

type MsgPackEncoder struct{}

func (j *MsgPackEncoder) Encode(topic string, message api.Message) (*nats.Msg, error) {
	msg := createMessage(topic, api.EncodingMsgPack, message)
	data, err := msgpack.Marshal(message)
	if err != nil {
		return nil, err
	}
	msg.Data = data
	return msg, nil
}

// endregion
// region - gob encoder

type GobEncoder struct{}

func (j *GobEncoder) Encode(topic string, message api.Message) (*nats.Msg, error) {
	msg := createMessage(topic, api.EncodingGob, message)
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(message)
	if err != nil {
		return nil, err
	}
	msg.Data = buff.Bytes()
	return msg, nil
}

// endregion
// region - common

func createMessage(topic string, encoding api.Encoding, message api.Message) *nats.Msg {
	msg := nats.NewMsg(topic)
	msg.Header.Add(msgEncodingHeader, encoding.String())
	msg.Header.Add(msgTypeHeader, fmt.Sprintf("%s", message.GetKind()))
	return msg
}

// endregion
