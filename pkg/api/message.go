package api

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type Message interface {
	GetKind() string
	GetTimestamp() time.Time
	GetData() interface{}
}

// region - basic message

func NewBasicMessage() Message {
	return &BasicMessage{
		Kind:      "BasicMessage",
		Timestamp: time.Now(),
	}
}

type BasicMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	Kind      string    `json:"type",msgpack:"-"`
	Timestamp time.Time `json:"timestamp",msgpack:"timestamp"`
}

func (m *BasicMessage) GetKind() string {
	return "BasicMessage"
}
func (m *BasicMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *BasicMessage) GetData() interface{} {
	return nil
}

func (m *BasicMessage) String() string {
	//v, e := /*easy*/ json.Marshal(m)
	//if e != nil {
	//	return "{}"
	//}
	//return string(v)
	return m.Kind + "{ timestamp=\"" + m.Timestamp.String() + "\" }"
}

// region ~ decoders

func GobBasicMessageDecoder(msg *nats.Msg) (Message, error) {
	r := bytes.NewReader(msg.Data)
	var m BasicMessage
	dec := gob.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func MsgPackBasicMessageDecoder(msg *nats.Msg) (Message, error) {
	var v BasicMessage
	err := msgpack.Unmarshal(msg.Data, &v)
	return &v, err
}
func JsonBasicMessageDecoder(msg *nats.Msg) (Message, error) {
	var bm BasicMessage
	err := json.Unmarshal(msg.Data, &bm)
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

// endregion

// endregion
// region - text message

func NewTextMessage(message string) Message {
	return &TextMessage{
		Kind:      "TextMessage",
		Timestamp: time.Now(),
		Text:      message,
	}
}

type TextMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	Kind      string    `json:"type",msgpack:"-"`
	Timestamp time.Time `json:"timestamp",msgpack:"timestamp"`
	Text      string    `json:"text",msgpack:"text"`
}

func (m *TextMessage) GetKind() string {
	return "TextMessage"
}
func (m *TextMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *TextMessage) GetData() interface{} {
	return m.Text
}

func (m *TextMessage) String() string {
	return m.Kind + "{ timestamp=\"" + m.Timestamp.String() + "\", text=\"" + m.Text + "\" }"
}

// region ~ decoders

func GobTextMessageDecoder(msg *nats.Msg) (Message, error) {
	r := bytes.NewReader(msg.Data)
	var m TextMessage
	dec := gob.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func MsgPackTextMessageDecoder(msg *nats.Msg) (Message, error) {
	var v TextMessage
	err := msgpack.Unmarshal(msg.Data, &v)
	return &v, err
}
func JsonTextMessageDecoder(msg *nats.Msg) (Message, error) {
	var bm TextMessage
	err := json.Unmarshal(msg.Data, &bm)
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

// endregion

// endregion
