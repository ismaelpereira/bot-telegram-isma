package controllers

import (
	"log"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HelpHandlerUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgHelp)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
