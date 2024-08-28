package message

import (
	"errors"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Recipient struct {
	ChatId string `json:"chat_id"`
}

func (r *Recipient) Recipient() string {
	return r.ChatId
}

type MessageType string

const (
	Text      MessageType = "Text"
	Audio     MessageType = "Audio"
	Document  MessageType = "Document"
	Photo     MessageType = "Photo"
	Sticker   MessageType = "Sticker"
	Voice     MessageType = "Voice"
	VideoNote MessageType = "MessageType"
	Video     MessageType = "Video"
	Animation MessageType = "Animation"
	Location  MessageType = "Location"
	Venue     MessageType = "Venue"
	Poll      MessageType = "Poll"
	Game      MessageType = "Game"
	Dice      MessageType = "Dice"
)

type FloodMessage struct {
	Type      MessageType `json:"message_type"`
	Recipient Recipient   `json:"recipient"`

	Text      *string         `json:"text,omitempty"`
	Audio     *tele.Audio     `json:"audio,omitempty"`
	Document  *tele.Document  `json:"document,omitempty"`
	Photo     *tele.Photo     `json:"photo,omitempty"`
	Sticker   *tele.Sticker   `json:"sticker,omitempty"`
	Voice     *tele.Voice     `json:"voice,omitempty"`
	VideoNote *tele.VideoNote `json:"video_note,omitempty"`
	Video     *tele.Video     `json:"video,omitempty"`
	Animation *tele.Animation `json:"animation,omitempty"`
	Location  *tele.Location  `json:"location,omitempty"`
	Venue     *tele.Venue     `json:"venue,omitempty"`
	Poll      *tele.Poll      `json:"poll,omitempty"`
	Game      *tele.Game      `json:"game,omitempty"`
	Dice      *tele.Dice      `json:"dice,omitempty"`

	SendOptions *tele.SendOptions `json:"send_options,omitempty"`
}

type WrappedMessage struct {
	BotToken     string       `json:"bot_token"`
	FloodMessage FloodMessage `json:"message"`
}

func (fm *FloodMessage) Send(b *tele.Bot, r tele.Recipient, so *tele.SendOptions) (*tele.Message, error) {

	var object tele.Sendable

	switch fm.Type {
	case Text:
		msg, err := b.Send(r, *fm.Text, so)
		if err != nil {
			return nil, err
		}
		return msg, nil
	case Audio:
		o := *fm.Audio
		object = &o
	case Document:
		o := *fm.Document
		object = &o
	case Photo:
		o := *fm.Photo
		object = &o
	case Sticker:
		o := *fm.Sticker
		object = &o
	case Voice:
		o := *fm.Voice
		object = &o
	case VideoNote:
		o := *fm.VideoNote
		object = &o
	case Video:
		o := *fm.Video
		object = &o
	case Animation:
		o := *fm.Animation
		object = &o
	case Location:
		o := *fm.Location
		object = &o
	case Venue:
		o := *fm.Venue
		object = &o
	case Poll:
		o := *fm.Poll
		object = &o
	case Game:
		o := *fm.Game
		object = &o
	case Dice:
		o := *fm.Dice
		object = &o
	default:
		return nil, errors.New("teleflood: now Sendable provided")
	}

	msg, err := b.Send(r, object, so)
	if err != nil {
		fmt.Println(err)
	}

	return msg, nil
}
