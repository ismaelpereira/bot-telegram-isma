package main

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/controllers"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/handler"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	c := config.Load()
	bot, err := tgbotapi.NewBotAPI(c.Telegram.Key)
	if err != nil {
		log.Println(err)
		return err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	go controllers.ReminderCheck(bot)
	for update := range updates {
		if err := handleUpdate(c, bot, &update); err != nil {
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func handleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		handler.CallbackActions(c, bot, update)
	}
	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}
	return handler.VerifyAndExecuteCommand(c, bot, update, strings.ToLower(update.Message.Command()))
}
