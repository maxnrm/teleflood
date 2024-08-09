package sender

import (
	"fmt"
	"time"

	"github.com/maxnrm/teleflood/config"
	"github.com/maxnrm/teleflood/pkg/message"
	"go.uber.org/ratelimit"
	tele "gopkg.in/telebot.v3"
)

type UserRateLimiter struct {
	ChatId      string
	RateLimiter ratelimit.Limiter
	LastMsgSent time.Time
}

type Sender struct {
	b *tele.Bot
	// TODO: delete items unused for some time
	globalRateLimiter ratelimit.Limiter
	userRateLimiters  map[string]*UserRateLimiter
}

func New(bot *tele.Bot, gl ratelimit.Limiter) (*Sender, error) {
	return &Sender{
		b:                bot,
		userRateLimiters: make(map[string]*UserRateLimiter),
	}, nil
}

func (s *Sender) Send(to tele.Recipient, fm *message.FloodMessage, so *tele.SendOptions) error {
	chatId := fm.Recipient.Recipient()

	// check if we can send complyting to user ratelimit
	rl, ok := s.userRateLimiters[chatId]
	if !ok {
		rl = &UserRateLimiter{
			ChatId:      chatId,
			RateLimiter: ratelimit.New(config.USER_RATE_LIMIT, ratelimit.WithoutSlack),
		}
		s.userRateLimiters[chatId] = rl
	}
	rl.RateLimiter.Take()

	// check if we can send complying to global rate limit
	s.globalRateLimiter.Take()

	_, err := s.b.Send(to, fm, so)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
