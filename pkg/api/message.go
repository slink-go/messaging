package api

import (
	"time"
)

type Message interface {
	GetKind() string
	GetTimestamp() time.Time
	GetData() interface{}
}

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
