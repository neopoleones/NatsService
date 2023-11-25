package nats

import (
	"L0_task/internal/config"
	"L0_task/internal/subscriber"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"sync/atomic"
)

var (
	ErrNoConnection = errors.New("not connected to nats")
	ErrEmptyMessage = errors.New("got empty message")
)

type Subscriber struct {
	cfg   *config.Subscriber
	subID atomic.Int64

	nc *nats.Conn
	js nats.JetStreamContext
}

func (s *Subscriber) Subscribe() (subscriber.Entry, error) {
	se := SubscriberEntry{
		Subscriber: s,
	}

	// Form name for new subscriber
	entryName := fmt.Sprintf("reader0%d", s.subID.Load())
	subEntry, err := se.js.PullSubscribe(s.cfg.Stream, entryName)
	if err != nil {
		return nil, err
	}
	se.subscription = subEntry

	return &se, nil
}

type SubscriberEntry struct {
	*Subscriber
	subscription *nats.Subscription
}

func (s *SubscriberEntry) PullMessage() (*subscriber.Message, error) {
	msgBatch, err := s.subscription.Fetch(1, nats.MaxWait(s.cfg.ReconnectWait/2))
	if err != nil || len(msgBatch) < 1 {
		return nil, err
	}

	// Size of batch - one element
	msg := msgBatch[0]
	if msg == nil {
		return nil, ErrEmptyMessage
	}
	_ = msg.Ack()

	return &subscriber.Message{
		Subject: msg.Subject,
		Header:  msg.Header,
		Data:    msg.Data,
	}, nil
}

func NewSubscriber(cfg *config.Subscriber) (*Subscriber, error) {
	sub := Subscriber{
		cfg: cfg,
	}

	// Open connection as subscriber
	nc, err := nats.Connect(
		sub.cfg.Address,
		nats.ReconnectWait(sub.cfg.ReconnectWait),
		nats.MaxReconnects(sub.cfg.MaxReconnect),
	)

	if err != nil {
		return nil, err
	}

	sub.nc = nc
	sub.js, err = sub.nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &sub, nil
}
