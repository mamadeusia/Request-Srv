package nats

import (
	"context"
	"errors"

	"go-micro.dev/v4/events"
)

var (
	ErrPullBasedConsumerNotFound = errors.New("nats handler : pull-based consumer not found")
	ErrPushBasedConsumerNotFound = errors.New("nats handler : push-based consumer not found")
)

type EventHandler func(ctx context.Context, e []*events.Event) error

type Topic interface {
	StartConsume(context.Context) error
}

type Handler struct {
	pullBasedConsumer events.Store
	pushBasedConsumer events.Stream
	topics            []Topic
}

type HandlerConfiguration func(a *Handler) error

func NewHandler(cfgs ...HandlerConfiguration) (*Handler, error) {
	handler := &Handler{}
	for _, cfg := range cfgs {
		err := cfg(handler)
		if err != nil {
			return nil, err
		}
	}
	return handler, nil
}

func (handler *Handler) Start(ctx context.Context) error {

	for _, topic := range handler.topics {
		if err := topic.StartConsume(ctx); err != nil {
			return err
		}
	}

	return nil
}
