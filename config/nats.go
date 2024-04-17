package config

import "strings"

type Nats struct {
	URL      string
	Callback Callback
	Nkey     string
}

type Callback struct {
	Prefix string
	Topics string
	Final  string
	Items  int
}

type Max struct {
	Tx int
}

func NatsCallbackTopics() []string {
	return strings.Split(cfg.Nats.Callback.Topics, ",")
}

func NatsCallbackPrefix() string {
	return cfg.Nats.Callback.Prefix
}

func NatsCallbackItems() int {
	return cfg.Nats.Callback.Items
}

func NatsCallbackFinalTopic() string {
	return cfg.Nats.Callback.Final
}

func NatsNkey() string {
	return cfg.Nats.Nkey
}
func NatsURL() string {
	return cfg.Nats.URL
}
