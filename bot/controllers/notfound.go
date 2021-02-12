package controllers

import (
	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//NotFoundHandlerUpdate send the message if is not a command
func NotFoundHandlerUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgCantUnderstand)
	_, err := bot.Send(msg)
	return err
}
