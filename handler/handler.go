package handler

import (
	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//VerifyAndExecuteCommand is a function to pick the right command and execute the respective function
func VerifyAndExecuteCommand(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update, cmd string) error {
	commands := map[string]func(*config.Config, *tgbotapi.BotAPI, *tgbotapi.Update) error{
		"help":    controllers.HelpHandlerUpdate,
		"admiral": controllers.AdmiralHandleUpdate,
		"anime":   controllers.AnimeHandleUpdate,
		"manga":   controllers.MangaHandleUpdate,
		"money":   controllers.MoneyHandleUpdate,
		"movie":   controllers.MovieHandleUpdate,
	}
	if f, ok := commands[cmd]; ok {
		err := f(c, bot, update)
		return err
	}
	err := controllers.NotFoundHandlerUpdate(c, bot, update)
	return err
}
