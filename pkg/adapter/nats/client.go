package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/slink-go/logging"
	"github.com/slink-go/messaging/pkg/api"
	"os"
)

// region - message handlers

type MessageHandlerFunc func(msg *nats.Msg) (api.Message, error)
type MessageHandlers map[string]MessageHandlerFunc

type MessageHandler struct {
	MessageType string
	Encoding    api.Encoding
	Handler     MessageHandlerFunc
}

func (h *MessageHandler) key() string {
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

func NewNatsClient() Client {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
		logging.GetLogger("nats-client").Info("using default Nats URL: %s", url)
	}
	return Client{
		url:             url,
		logger:          logging.GetLogger("nats-client"),
		messageHandlers: make(MessageHandlers),
		encoder:         Encoder{},
		decoder:         Decoder{},
	}
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

func (c *Client) AddMessageHandler(handler MessageHandler) {
	c.messageHandlers[handler.key()] = handler.Handler
}

func (c *Client) encode(topic string, message api.Message, encoding api.Encoding) (*nats.Msg, error) {
	return c.encoder.Encode(topic, message, encoding)
}

func (c *Client) decode(msg *nats.Msg) (api.Message, error) {
	return c.decoder.Decode(msg, c.messageHandlers)
}

// endregion
