package api

import "strings"

type Encoding int

const (
	EncodingUndefined Encoding = iota
	EncodingString
	EncodingJson
	EncodingGob
	EncodingMsgPack
)

func (e Encoding) String() string {
	switch e {
	case EncodingString:
		return "string"
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
	case "string":
		return EncodingString
	case "json":
		return EncodingJson
	case "msgpack":
		return EncodingMsgPack
	case "gob":
		return EncodingGob
	}
	return EncodingUndefined
}
