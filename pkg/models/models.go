package models

import (
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
	Audio     MessageType = "Audio"
	Document  MessageType = "Document"
	Photo     MessageType = "Photo"
	Sticker   MessageType = "Sticker"
	Voice     MessageType = "Voice"
	VideoNote MessageType = "MessageType"
	Video     MessageType = "Video"
	Animation MessageType = "Animation"
	Contact   MessageType = "Contact"
	Location  MessageType = "Location"
	Venue     MessageType = "Venue"
	Poll      MessageType = "Poll"
	Game      MessageType = "Game"
	Dice      MessageType = "Dice"
)

type FloodMessage struct {
	Type      MessageType `json:"message_type"`
	BotToken  string      `json:"bot_token"`
	Recipient Recipient   `json:"recipient"`

	Audio     *tele.Audio     `json:"audio"`
	Document  *tele.Document  `json:"document"`
	Photo     *tele.Photo     `json:"photo"`
	Sticker   *tele.Sticker   `json:"sticker"`
	Voice     *tele.Voice     `json:"voice"`
	VideoNote *tele.VideoNote `json:"video_note"`
	Video     *tele.Video     `json:"video"`
	Animation *tele.Animation `json:"animation"`
	Contact   *tele.Contact   `json:"contact"`
	Location  *tele.Location  `json:"location"`
	Venue     *tele.Venue     `json:"venue"`
	Poll      *tele.Poll      `json:"poll"`
	Game      *tele.Game      `json:"game"`
	Dice      *tele.Dice      `json:"dice"`

	SendOptions tele.SendOptions `json:"send_options"`
}
