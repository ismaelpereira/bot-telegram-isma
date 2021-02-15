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
	c := config.Load()

	telegramKey := c.Telegram.Key
	bot, err := tgbotapi.NewBotAPI(telegramKey)
	if err != nil {
		log.Println(err)
		return err
	}
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if err := handleUpdate(c, bot, &update); err != nil {
			continue
		}
	}
	return nil
}

func handleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, "movie:") {
			update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "movie:")
			return controllers.MovieHandleUpdate(c, bot, update)
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "serie:") {
			update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "serie:")
			return controllers.SeriesHandleUpdate(c, bot, update)
		}
	}
	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}
	return handler.VerifyAndExecuteCommand(c, bot, update, strings.ToLower(update.Message.Command()))
}
