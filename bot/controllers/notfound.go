package controllers

import (
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

//NotFoundHandlerUpdate send the message if is not a command
func NotFoundHandlerUpdate(cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgCantUnderstand)
	_, err := bot.Send(msg)
	return err
}
