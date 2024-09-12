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
		Source:    "",
		Timestamp: time.Now(),
	}
}

type BasicMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	Kind      string    `json:"type",msgpack:"-"`
	Source    string    `json:"source",msgpack:"-"`
	Timestamp time.Time `json:"timestamp",msgpack:"timestamp"`
}

func (m *BasicMessage) GetKind() string {
	return "BasicMessage"
}
func (m *BasicMessage) GetSource() string {
	return m.Source
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
	return fmt.Sprintf("{ kind=\"%s\", src=\"%s\", timestamp=\"%s\"}", m.Kind, m.Source, m.Timestamp)
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
	m.Source = msg.Subject
	return &m, nil
}
func MsgPackBasicMessageDecoder(msg *nats.Msg) (Message, error) {
	var m BasicMessage
	err := msgpack.Unmarshal(msg.Data, &m)
	m.Source = msg.Subject
	return &m, err
}
func JsonBasicMessageDecoder(msg *nats.Msg) (Message, error) {
	var m BasicMessage
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		return nil, err
	}
	m.Source = msg.Subject
	return &m, nil
}

// endregion

// endregion
// region - text message

func NewTextMessage(message string) Message {
	return &TextMessage{
		BasicMessage: BasicMessage{
			Kind:      "TextMessage",
			Source:    "",
			Timestamp: time.Now(),
		},
		Text: message,
	}
}

type TextMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	BasicMessage
	Text string `json:"text",msgpack:"text"`
}

func (m *TextMessage) GetKind() string {
	return "TextMessage"
}
func (m *TextMessage) GetSource() string {
	return m.Source
}
func (m *TextMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *TextMessage) GetData() interface{} {
	return m.Text
}

func (m *TextMessage) String() string {
	return fmt.Sprintf("{ kind=\"%s\", src=\"%s\", timestamp=\"%s\", text=\"%v\"}", m.Kind, m.Source, m.Timestamp, m.Text)
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
	m.Source = msg.Subject
	return &m, nil
}
func MsgPackTextMessageDecoder(msg *nats.Msg) (Message, error) {
	var m TextMessage
	err := msgpack.Unmarshal(msg.Data, &m)
	m.Source = msg.Subject
	return &m, err
}
func JsonTextMessageDecoder(msg *nats.Msg) (Message, error) {
	var m TextMessage
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		return nil, err
	}
	m.Source = msg.Subject
	return &m, nil
}

// endregion

// endregion
// region - object message

func NewObjectMessage(object map[string]any) Message {
	return &ObjectMessage{
		BasicMessage: BasicMessage{
			Kind:      "ObjectMessage",
			Timestamp: time.Now(),
		},
		Object: object,
	}
}

type ObjectMessage struct {
	_msgpack struct{} `msgpack:",as_array"`

	BasicMessage
	Object map[string]any `json:"object",msgpack:"object"`
}

func (m *ObjectMessage) GetKind() string {
	return "ObjectMessage"
}
func (m *ObjectMessage) GetSource() string {
	return m.Source
}
func (m *ObjectMessage) GetTimestamp() time.Time {
	return m.Timestamp
}
func (m *ObjectMessage) GetData() interface{} {
	return m.Object
}

func (m *ObjectMessage) String() string {
	return fmt.Sprintf("{ kind=\"%s\", src=\"%s\", timestamp=\"%s\", object=%v}", m.Kind, m.Source, m.Timestamp, m.Object)
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
	m.Source = msg.Subject
	return &m, nil
}
func MsgPackObjectMessageDecoder(msg *nats.Msg) (Message, error) {
	var m ObjectMessage
	err := msgpack.Unmarshal(msg.Data, &m)
	m.Source = msg.Subject
	return &m, err
}
func JsonObjectMessageDecoder(msg *nats.Msg) (Message, error) {
	var m ObjectMessage
	err := json.Unmarshal(msg.Data, &m)
	if err != nil {
		return nil, err
	}
	m.Source = msg.Subject
	return &m, nil
}

// endregion

// endregion
