package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
	"os"
)

// region - message handlers

type MessageDecoderFunc func(msg *nats.Msg) (api.Message, error)
type MessageHandlers map[string]MessageDecoderFunc

type MessageDecoder struct {
	MessageType string
	Encoding    api.Encoding
	Handler     MessageDecoderFunc
}

func (h *MessageDecoder) key() string {
	return h.Encoding.String() + ":" + h.MessageType
}

// endregion
// region - nats client

type Client struct {
	url             string
	conn            *nats.Conn
	logger          logging.Logger
	messageHandlers MessageHandlers
	encoder         Encoder
	decoder         Decoder
}

func NewNatsClient(opts ...interface{}) Client {
	l := logging.GetLogger("nats-client")
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
		l.Info("using default Nats URL: %s", url)
	}
	enc := getDefaultEncoder(api.EncodingJson, opts...)
	return Client{
		url:             url,
		logger:          l,
		messageHandlers: make(MessageHandlers),
		encoder:         enc,
	}
}
func getDefaultEncoder(defaultEncoding api.Encoding, opts ...interface{}) Encoder {
	if opts != nil || len(opts) > 0 {
		v, ok := opts[0].(api.Encoding)
		if ok {
			if enc := getEncoder(v); enc != nil {
				return enc
			}
		}
	}
	return getEncoder(defaultEncoding)
}
func getEncoder(encoding api.Encoding) Encoder {
	enc, err := GetEncoder(encoding)
	if err != nil {
		logging.GetLogger("nats-client").Warning("could not get default encoder: %v", err)
		return nil
	}
	return enc
}

func (c *Client) Connect() error {
	nc, err := nats.Connect(c.url)
	if err != nil {
		return err
	}
	c.conn = nc
	nc.SetClosedHandler(func(cc *nats.Conn) {
		c.logger.Info("connection closed: %v", cc)
	})
	return nil
}
func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Drain()
	}
}

func (c *Client) AddMessageDecoder(handler MessageDecoder) {
	c.messageHandlers[handler.key()] = handler.Handler
}

func (c *Client) encode(topic string, message api.Message, encoding api.Encoding) (*nats.Msg, error) {
	return c.encoder.Encode(topic, message)
}

func (c *Client) decode(msg *nats.Msg) (api.Message, error) {
	return c.decoder.Decode(msg, c.messageHandlers)
}

// endregion
