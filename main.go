package main

import (
	"log"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
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
	if update.Message == nil {
		return nil
	}
	if update.Message.IsCommand() {
		argument := update.Message.Command()
		switch strings.ToLower(argument) {
		case "help":
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgHelp)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
					return err
				}

			}
		case "admiral":
			{
				err := controllers.AdmiralHandleUpdate(bot, update)
				if err != nil {
					log.Println(err)
					return err
				}

			}
		case "anime":
			{
				err := controllers.AnimeHandleUpdate(bot, update)
				if err != nil {
					log.Println(err)
					return err
				}
			}
		case "manga":
			{
				err := controllers.MangaHandleUpdate(bot, update)
				if err != nil {
					log.Println(err)
					return err
				}
			}
		default:
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgCantUnderstand)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}

			}
		}
	}
	return nil
}
