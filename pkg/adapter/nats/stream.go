package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.slink.ws/logging"
	"go.slink.ws/messaging/pkg/api"
)

type NatMessageStream struct {
	client         Client
	config         api.StreamConfig
	logger         logging.Logger
	defaultEncoder Encoder
}

func (s *NatMessageStream) Create() error {
	if err := validateClient(s.client); err != nil {
		return err
	}
	js, err := s.getJsContext()
	if err != nil {
		return err
	}
	return s.createStream(js)
}
func (s *NatMessageStream) Delete() error {
	if err := validateClient(s.client); err != nil {
		return err
	}
	js, err := s.getJsContext()
	if err != nil {
		return err
	}
	return js.DeleteStream(s.config.Name) // TODO: options
}
func (s *NatMessageStream) Flush() error {
	if err := validateClient(s.client); err != nil {
		return err
	}
	js, err := s.getJsContext()
	if err != nil {
		return err
	}
	return js.PurgeStream(s.config.Name) // TODO: options
}

func (s *NatMessageStream) Publish(subject string, message api.Message, encoding ...api.Encoding) error {
	if encoding != nil && len(encoding) > 0 {
		enc, err := GetEncoder(encoding[0])
		if err != nil {
			return err
		}
		return s.publish(subject, message, enc)
	}
	return s.publish(subject, message, s.client.encoder)
}
func (s *NatMessageStream) publish(subject string, message api.Message, encoder Encoder) error {

	// encode message
	msg, err := encoder.Encode(subject, message)
	if err != nil {
		return err
	}

	// get nats JetStream context
	jsctx, err := s.getJsContext()
	if err != nil {
		return err
	}

	// create stream if needed
	err = s.createStream(jsctx)
	if err != nil {
		return err
	}

	_, err = jsctx.PublishMsg(msg)
	if err != nil {
		return err
	}

	return err
}

func (s *NatMessageStream) Subscribe(subject string, channel chan api.Message, options ...api.SubscribeOption) (func(), error) {
	if err := validateClient(s.client); err != nil {
		return nil, err
	}

	// get nats JetStream context
	jsctx, err := s.getJsContext()
	if err != nil {
		return nil, err
	}

	sub, err := jsctx.Subscribe(subject, func(msg *nats.Msg) {
		m, e := s.client.decode(msg)
		if e != nil {
			s.logger.Warning("could not decode incoming message: %v", e)
			return
		}
		if shouldAcceptMessage(m, options...) {
			//s.logger.Trace("new message received in %s: %v", msg.Subject, m)
			if err := s.send(channel, m); err != nil {
				s.logger.Warning("could not send message: %v", err)
			}
		}
	},
		// опции подписки на stream
		nats.OrderedConsumer(),
		nats.DeliverAll(),
	)

	closer := func() {
		s.logger.Info("stream unsubscribe from %s", subject)
		_ = sub.Unsubscribe()
	}

	return closer, err
}

func (s *NatMessageStream) getJsContext() (nats.JetStreamContext, error) {
	if err := validateClient(s.client); err != nil {
		return nil, err
	}
	jsctx, err := s.client.conn.JetStream(nats.PublishAsyncMaxPending(256)) // TODO: replace with constant
	if err != nil {
		return nil, err
	}
	return jsctx, nil
}
func (s *NatMessageStream) createStream(jetStream nats.JetStreamContext) error {
	stream, err := jetStream.StreamInfo(s.config.Name)
	// stream not found, create it
	if stream == nil {
		s.logger.Info("Creating stream: %s\n", s.config.Name)
		_, err = jetStream.AddStream(&nats.StreamConfig{
			Name:        s.config.Name,
			Description: s.config.Description,
			Subjects:    s.config.Subjects,
			Compression: nats.S2Compression,
			//Storage:     nats.FileStorage(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NatMessageStream) send(channel chan api.Message, message api.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("message send error: %v", r)
		}
	}()
	channel <- message
	return
}
