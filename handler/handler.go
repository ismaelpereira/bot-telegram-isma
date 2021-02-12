package handler

import (
	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//VerifyAndExecuteCommand is a function to pick the right command and execute the respective function
func VerifyAndExecuteCommand(cmd string, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var commands map[string]func(*tgbotapi.BotAPI, *tgbotapi.Update) error = map[string]func(*tgbotapi.BotAPI, *tgbotapi.Update) error{
		"help":    controllers.HelpHandlerUpdate,
		"admiral": controllers.AdmiralHandleUpdate,
		"anime":   controllers.AnimeHandleUpdate,
		"manga":   controllers.MangaHandleUpdate,
		"money":   controllers.MoneyHandleUpdate,
		"movie":   controllers.MovieHandleUpdate,
	}
	if f, ok := commands[cmd]; ok {
		err := f(bot, update)
		if err != nil {
			return err
		}
		return nil
	}
	err := controllers.NotFoundHandlerUpdate(bot, update)
	if err != nil {
		return err
	}
	return nil
}
