package sender

import (
	"context"
	"fmt"
	"time"

	"github.com/maxnrm/teleflood/config"
	"github.com/maxnrm/teleflood/pkg/message"
	"golang.org/x/time/rate"
	tele "gopkg.in/telebot.v3"
)

type UserRateLimiter struct {
	ChatId      string
	RateLimiter *rate.Limiter
	LastMsgSent time.Time
}

type Sender struct {
	b *tele.Bot
	// TODO: delete items unused for some time
	globalRateLimiter *rate.Limiter
	userRateLimiters  map[string]*UserRateLimiter
}

func New(bot *tele.Bot, gl *rate.Limiter) (*Sender, error) {
	return &Sender{
		b:                bot,
		userRateLimiters: make(map[string]*UserRateLimiter),
	}, nil
}

func (s *Sender) Send(to tele.Recipient, fm *message.FloodMessage, so *tele.SendOptions) error {

	ctx := context.Background()
	chatId := fm.Recipient.Recipient()

	// check if we can send complyting to user ratelimit
	rl, ok := s.userRateLimiters[chatId]
	if !ok {
		limit := rate.Every(time.Second / time.Duration(config.USER_RATE_LIMIT))
		rl = &UserRateLimiter{
			ChatId:      chatId,
			RateLimiter: rate.NewLimiter(limit, config.USER_RATE_LIMIT),
		}
		s.userRateLimiters[chatId] = rl
	}
	rl.RateLimiter.Wait(ctx)

	// check if we can send complying to global rate limit
	s.globalRateLimiter.Wait(ctx)

	_, err := s.b.Send(to, fm, so)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
