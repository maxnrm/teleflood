package provider

import (
	"context"
	"encoding/json"
	"time"

	"github.com/maxnrm/teleflood/config"
	"github.com/maxnrm/teleflood/internal/nats"
	m "github.com/maxnrm/teleflood/pkg/message"
	"github.com/nats-io/nats.go/jetstream"
)

type Provider struct {
	client jetstream.Consumer
}

func New() *Provider {

	settings := nats.NatsSettings{
		URL: config.NATS_URL,
		Ctx: context.Background(),
	}

	natsClient := *nats.Init(settings)

	streamConfig := jetstream.StreamConfig{
		Name:      config.NATS_MESSAGES_STREAM,
		Subjects:  []string{config.NATS_MESSAGES_SUBJECT},
		Retention: jetstream.WorkQueuePolicy,
		Storage:   jetstream.FileStorage,
	}

	consumerConfig := jetstream.ConsumerConfig{
		Name:          config.NATS_MESSAGES_CONSUMER,
		Durable:       config.NATS_MESSAGES_CONSUMER,
		FilterSubject: config.NATS_MESSAGES_SUBJECT,
		AckWait:       2 * time.Second,
		MaxAckPending: 60,
		MemoryStorage: true,
	}

	natsClient.CreateStream(streamConfig)
	consumer := natsClient.CreateConsumer(streamConfig.Name, consumerConfig)

	return &Provider{
		client: consumer,
	}
}

func (p *Provider) Next() (*m.FloodMessage, error) {
	msg, err := p.client.Next()
	if err != nil {
		return nil, err
	}

	var floodMsg m.FloodMessage

	json.Unmarshal(msg.Data(), &floodMsg)

	return &floodMsg, nil
}
