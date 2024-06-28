package models

import (
	"fmt"
	"teleflood/config"
	"teleflood/internal/sendlimiter"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Recipient struct {
	ChatID string
}

func (r *Recipient) Recipient() string {
	return r.ChatID
}

type SendableMessage struct {
	Text        *string           `json:"text,omitempty"`
	Photo       *Photo            `json:"photo,omitempty"`
	SendOptions *tele.SendOptions `json:"send_options,omitempty"`
	Variant     int               `json:"variant"`
	Recipient   *Recipient        `json:"recipient"`
}

type Photo struct {
	tele.File

	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Caption string `json:"caption,omitempty"`
}

func (sm *SendableMessage) createWhat() interface{} {
	var what interface{}

	if sm.Text != nil {
		what = *sm.Text
	} else {
		what = &tele.Photo{File: tele.FromURL(sm.Photo.FileURL), Caption: sm.Photo.Caption}
	}

	return what
}

func (sm *SendableMessage) getSendOptions() *tele.SendOptions {
	if sm.SendOptions != nil {
		return sm.SendOptions
	}

	return &tele.SendOptions{}

}

func (sm *SendableMessage) sendWithLimit(bot *tele.Bot, limiter *sendlimiter.SendLimiter) error {
	chatID := sm.Recipient.Recipient()

	userRateLimiter := limiter.GetUserRateLimiter(chatID)
	if userRateLimiter == nil {
		limiter.AddUserRateLimiter(chatID, config.RATE_LIMIT_USER, config.RATE_LIMIT_BURST_USER)
		userRateLimiter = limiter.GetUserRateLimiter(chatID)
	}

	err := userRateLimiter.RateLimiter.Wait(limiter.Ctx)
	if err != nil {
		return err
	}

	err = limiter.GlobalRateLimiter.Wait(limiter.Ctx)
	if err != nil {
		return err
	}

	fmt.Println("Tokens left:", limiter.GlobalRateLimiter.Tokens())

	what := sm.createWhat()
	opts := sm.getSendOptions()

	_, err = bot.Send(sm.Recipient, what, opts)
	if err != nil {
		fmt.Println(err)
		return err
	}

	userRateLimiter.LastMsgSent = time.Now()

	return nil

}

func (sm *SendableMessage) Send(bot *tele.Bot, limiter *sendlimiter.SendLimiter) error {
	return sm.sendWithLimit(bot, limiter)
}
