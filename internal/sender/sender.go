package sender

import (
	"errors"

	"github.com/maxnrm/teleflood/config"
	m "github.com/maxnrm/teleflood/pkg/message"
	"go.uber.org/ratelimit"
	tele "gopkg.in/telebot.v3"
)

type Sender struct {
	b *tele.Bot
	// TODO: delete items unused for some time
	userRateLimit map[string]ratelimit.Limiter
}

func New(botToken string) (*Sender, error) {
	b, err := tele.NewBot(tele.Settings{
		Token: botToken,
	})

	if err != nil {
		return nil, err
	}

	return &Sender{
		b: b,
	}, nil
}

func (s *Sender) Send(grl ratelimit.Limiter, fm *m.FloodMessage) error {

	var object tele.Sendable

	switch fm.Type {
	case m.Text:
		_, err := s.b.Send(&fm.Recipient, fm.Text, fm.SendOptions)
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

	chatId := fm.Recipient.Recipient()

	// check if we can send complyting to user ratelimit
	if rl, ok := s.userRateLimit[chatId]; !ok {
		s.userRateLimit[chatId] = ratelimit.New(config.USER_RATE_LIMIT)
		rl := s.userRateLimit[chatId]
		rl.Take()
	} else {
		rl.Take()
	}

	// check if we can send complying to global rate limit
	grl.Take()
	_, err := s.b.Send(&fm.Recipient, object, fm.SendOptions)

	return err
}
