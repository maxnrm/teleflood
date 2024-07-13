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
	cons jetstream.Consumer
	mc   jetstream.MessagesContext
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
	cons := natsClient.CreateConsumer(streamConfig.Name, consumerConfig)
	mc, _ := cons.Messages()

	return &Provider{
		cons: cons,
		mc:   mc,
	}
}

func (p *Provider) Next(ctx context.Context) (*m.FloodMessageWithToken, error) {
	msg, err := p.mc.Next()
	if err != nil {
		return nil, err
	}

	var floodMsg m.FloodMessageWithToken

	err = json.Unmarshal(msg.Data(), &floodMsg)
	if err != nil {
		return nil, err
	}

	// TODO: ack only after message been sent without error
	msg.DoubleAck(ctx)
	return &floodMsg, nil
}
