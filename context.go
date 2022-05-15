package tgbot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageConfigOption func(c *tgbotapi.MessageConfig)

type Context struct {
	context.Context

	*tgbotapi.BotAPI

	update *tgbotapi.Update
}

func (ctx *Context) Command() string {
	if message := ctx.Message(); message != nil {
		return message.Command()
	}
	return ""
}

func (ctx *Context) CommandArgs() string {
	if message := ctx.Message(); message != nil {
		return message.CommandArguments()
	}
	return ""
}

func (ctx *Context) Message() *tgbotapi.Message {
	switch {
	case ctx.update.Message != nil:
		return ctx.update.Message
	case ctx.update.EditedMessage != nil:
		return ctx.update.EditedMessage
	case ctx.update.ChannelPost != nil:
		return ctx.update.ChannelPost
	case ctx.update.EditedChannelPost != nil:
		return ctx.update.EditedChannelPost
	default:
		return nil
	}
}

func (ctx *Context) Update() *tgbotapi.Update {
	return ctx.update
}

func (ctx *Context) SentFrom() *tgbotapi.User {
	return ctx.update.SentFrom()
}

func (ctx *Context) FromChat() *tgbotapi.Chat {
	return ctx.update.FromChat()
}

func (ctx *Context) ReplyText(text string, opts ...MessageConfigOption) error {
	return ctx.reply(text, opts...)
}

func (ctx *Context) ReplyMarkdown(text string, opts ...MessageConfigOption) error {
	return ctx.reply(text, mergeOpts(func(c *tgbotapi.MessageConfig) {
		c.ParseMode = tgbotapi.ModeMarkdown
	}, opts)...)
}

func (ctx *Context) ReplyHTML(text string, opts ...MessageConfigOption) error {
	return ctx.reply(text, mergeOpts(func(c *tgbotapi.MessageConfig) {
		c.ParseMode = tgbotapi.ModeHTML
	}, opts)...)
}

func (ctx *Context) reply(text string, opts ...MessageConfigOption) error {
	msg := tgbotapi.NewMessage(ctx.update.FromChat().ID, text)
	msg.DisableWebPagePreview = true
	for _, o := range opts {
		o(&msg)
	}
	_, err := ctx.Send(msg)
	return err
}

func mergeOpts(a MessageConfigOption, b []MessageConfigOption) []MessageConfigOption {
	return append([]MessageConfigOption{a}, b...)
}

// WithEnableWebPagePreview enable web page preview
func WithEnableWebPagePreview() MessageConfigOption {
	return func(c *tgbotapi.MessageConfig) {
		c.DisableWebPagePreview = false
	}
}

// WithChatId set message chat id
func WithChatId(chatId int64) MessageConfigOption {
	return func(c *tgbotapi.MessageConfig) {
		c.ChatID = chatId
	}
}
