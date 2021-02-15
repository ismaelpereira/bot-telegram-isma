package controllers

import (
	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//HelpHandlerUpdate send the help message
func HelpHandlerUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgHelp)
	_, err := bot.Send(msg)
	return err
}
