package nats

import (
	"go-micro.dev/v4/events"
)

func WithPullBasedStore(es events.Store) HandlerConfiguration {
	return func(h *Handler) error {
		h.pullBasedConsumer = es
		return nil
	}
}

func WithPushBasedStream(es events.Stream) HandlerConfiguration {
	return func(h *Handler) error {
		h.pushBasedConsumer = es
		return nil
	}
}

func WithPullBasedTopic(pbt *PullBasedTopic) HandlerConfiguration {
	return func(h *Handler) error {
		if h.topics == nil {
			h.topics = []Topic{}
		}
		if h.pullBasedConsumer == nil {
			return ErrPullBasedConsumerNotFound
		}
		//here we need to set the consumer of the topic
		pbt.pullBasedConsumer = h.pullBasedConsumer
		h.topics = append(h.topics, pbt)
		return nil
	}
}

func WithPushBasedTopic(pbt *PushBasedTopic) HandlerConfiguration {
	return func(h *Handler) error {
		if h.topics == nil {
			h.topics = []Topic{}
		}
		if h.pushBasedConsumer == nil {
			return ErrPushBasedConsumerNotFound
		}
		//here we need to set the consumer of the topic
		pbt.pushBasedConsumer = h.pushBasedConsumer
		h.topics = append(h.topics, pbt)
		return nil
	}
}
