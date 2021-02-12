package main

import (
	"log"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	telegramPath, err := config.GetTelegramKey()
	if err != nil {
		log.Println(err)
		return err
	}
	bot, err := tgbotapi.NewBotAPI(telegramPath)
	if err != nil {
		log.Println(err)
		return err
	}
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if err := handleUpdate(bot, &update); err != nil {
			log.Println(err)
			continue
		}
	}
	return nil
}

func handleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		controllers.MovieHandleUpdate(bot, update)
	}
	if update.Message == nil {
		return nil
	}
	if update.Message.IsCommand() {
		err := handler.VerifyAndExecuteCommand(strings.ToLower(update.Message.Command()), bot, update)
		if err != nil {
			return err
		}
	}
	return nil
}
