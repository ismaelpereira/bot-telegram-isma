package controllers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

//HelpHandlerUpdate send the help message
func HelpHandlerUpdate(cfg *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgHelp)
	_, err := bot.Send(msg)
	return err
}
