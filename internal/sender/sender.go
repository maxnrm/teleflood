package sender

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/maxnrm/teleflood/config"
	m "github.com/maxnrm/teleflood/pkg/message"
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
	userRateLimit map[string]*UserRateLimiter
}

func New(botToken string) (*Sender, error) {
	b, err := tele.NewBot(tele.Settings{
		Token: botToken,
	})

	if err != nil {
		return nil, err
	}

	return &Sender{
		b:             b,
		userRateLimit: make(map[string]*UserRateLimiter),
	}, nil
}

func (s *Sender) Send(ctx context.Context, grl *rate.Limiter, fm *m.FloodMessage) error {

	var object tele.Sendable

	var sendOptions *tele.SendOptions
	if fm.SendOptions == nil {
		sendOptions = &tele.SendOptions{Protected: true}
	} else {
		sendOptions = fm.SendOptions
	}

	chatId := fm.Recipient.Recipient()

	// check if we can send complyting to user ratelimit
	rl, ok := s.userRateLimit[chatId]
	if !ok {
		limit := rate.Every(time.Second / time.Duration(config.USER_RATE_LIMIT))
		rl = &UserRateLimiter{
			ChatId:      chatId,
			RateLimiter: rate.NewLimiter(limit, config.USER_RATE_LIMIT),
		}
		s.userRateLimit[chatId] = rl
	}
	rl.RateLimiter.Wait(ctx)

	// check if we can send complying to global rate limit
	grl.Wait(ctx)

	switch fm.Type {
	case m.Text:
		_, err := s.b.Send(&fm.Recipient, *fm.Text, sendOptions)
		if err != nil {
			fmt.Println(err)
		}
		return err
	case m.Audio:
		object = fm.Audio
	case m.Document:
		object = fm.Document
	case m.Photo:
		object = fm.Photo
	case m.Sticker:
		object = fm.Sticker
	case m.Voice:
		object = fm.Voice
	case m.VideoNote:
		object = fm.VideoNote
	case m.Video:
		object = fm.Video
	case m.Animation:
		object = fm.Animation
	case m.Location:
		object = fm.Location
	case m.Venue:
		object = fm.Venue
	case m.Poll:
		object = fm.Poll
	case m.Game:
		object = fm.Game
	case m.Dice:
		object = fm.Dice
	default:
		return errors.New("teleflood: now Sendable provided")
	}

	_, err := s.b.Send(&fm.Recipient, object, sendOptions)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
