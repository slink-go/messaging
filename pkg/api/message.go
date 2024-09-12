package api

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/goccy/go-json"
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
	return fmt.Sprintf("{ kind=\"%s\", timestamp=\"%s\"}", m.Kind, m.Timestamp)
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
	return fmt.Sprintf("{ kind=\"%s\", timestamp=\"%s\", text=\"%v\"}", m.Kind, m.Timestamp, m.Text)
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
// region - object message

func NewObjectMessage(object map[string]any) Message {
	return &ObjectMessage{
		Kind:      "ObjectMessage",
		Timestamp: time.Now(),
		Object:    object,
	}
}

type ObjectMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	Kind      string         `json:"type",msgpack:"-"`
	Timestamp time.Time      `json:"timestamp",msgpack:"timestamp"`
	Object    map[string]any `json:"object",msgpack:"object"`
}

func (m *ObjectMessage) GetKind() string {
	return "ObjectMessage"
}
func (m *ObjectMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *ObjectMessage) GetData() interface{} {
	return m.Object
}

func (m *ObjectMessage) String() string {
	return fmt.Sprintf("{ kind=\"%s\", timestamp=\"%s\", object=%v}", m.Kind, m.Timestamp, m.Object)
}

// region ~ decoders

func GobObjectMessageDecoder(msg *nats.Msg) (Message, error) {
	r := bytes.NewReader(msg.Data)
	var m ObjectMessage
	dec := gob.NewDecoder(r)
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
func MsgPackObjectMessageDecoder(msg *nats.Msg) (Message, error) {
	var v ObjectMessage
	err := msgpack.Unmarshal(msg.Data, &v)
	return &v, err
}
func JsonObjectMessageDecoder(msg *nats.Msg) (Message, error) {
	var bm ObjectMessage
	err := json.Unmarshal(msg.Data, &bm)
	if err != nil {
		return nil, err
	}
	return &bm, nil
}

// endregion

// endregion
