package api

import "strings"

type Encoding int

const (
	EncodingUndefined Encoding = iota
	EncodingJson
	EncodingGob
	EncodingMsgPack
)

func (e Encoding) String() string {
	switch e {
	case EncodingJson:
		return "json"
	case EncodingMsgPack:
		return "msgpack"
	case EncodingGob:
		return "gob"
	}
	return "unknown"
}
func ParseEncoding(input string) Encoding {
	switch strings.ToLower(input) {
	case "json":
		return EncodingJson
	case "msgpack":
		return EncodingMsgPack
	case "gob":
		return EncodingGob
	}
	return EncodingUndefined
}
