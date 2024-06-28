package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"teleflood/config"
	"teleflood/internal/models"
	"teleflood/internal/nats"
	"teleflood/internal/sendlimiter"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	tele "gopkg.in/telebot.v3"
)

var wg sync.WaitGroup
var ctx = context.Background()
var sl = sendlimiter.Init(ctx, config.RATE_LIMIT_GLOBAL, config.RATE_LIMIT_BURST_GLOBAL)

var nc *nats.NatsClient = nats.Init(nats.NatsSettings{
	Ctx: ctx,
	URL: config.NATS_URL,
})

var botSender, _ = tele.NewBot(tele.Settings{
	Token:  config.BOT_TOKEN,
	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
})

var streamConfig = jetstream.StreamConfig{
	Name:      config.NATS_MESSAGES_STREAM,
	Subjects:  []string{config.NATS_MESSAGES_SUBJECT},
	Retention: jetstream.WorkQueuePolicy,
	Storage:   jetstream.FileStorage,
}

var consumerConfig = jetstream.ConsumerConfig{
	Name:          config.NATS_MESSAGES_CONSUMER,
	Durable:       config.NATS_MESSAGES_CONSUMER,
	FilterSubject: config.NATS_MESSAGES_SUBJECT,
	AckWait:       2 * time.Second,
	MaxAckPending: 60,
	MemoryStorage: true,
}

func main() {
	nc.CreateStream(streamConfig)

	cons := nc.CreateConsumer(streamConfig.Name, consumerConfig)
	messageHandler := createConsumeHandler(ctx, botSender, sl)

	wg.Add(4)

	go sl.RemoveOldUserRateLimitersCache(20)

	cons.Consume(messageHandler)

	fmt.Println("Consuming...")

	wg.Wait()
}

func createConsumeHandler(ctx context.Context, bot *tele.Bot, limiter *sendlimiter.SendLimiter) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		var sendableMessage models.SendableMessage

		err := json.Unmarshal(msg.Data(), &sendableMessage)
		if err != nil {
			fmt.Println("Error while unmarshalling sendableMessage from json:", err)
			return
		}

		msg.DoubleAck(ctx)

		go sendableMessage.Send(bot, limiter)
	}
}
