package handler

import (
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//VerifyAndExecuteCommand is a function to pick the right command and execute the respective function
func VerifyAndExecuteCommand(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update, cmd string) error {
	Commands := map[string]types.HandlerFunc{
		"help":     controllers.HelpHandlerUpdate,
		"admirals": controllers.AdmiralsHandleUpdate,
		"animes":   controllers.AnimesHandleUpdate,
		"mangas":   controllers.MangasHandleUpdate,
		"money":    controllers.MoneyHandleUpdate,
		"movies":   controllers.MoviesHandleUpdate,
		"series":   controllers.SeriesHandleUpdate,
		"reminder": controllers.TimerHandleUpdate,
		"now":      controllers.TimerHandleUpdate,
	}
	if f, ok := Commands[cmd]; ok {
		return f(c, bot, update)
	}
	return controllers.NotFoundHandlerUpdate(c, bot, update)
}

func CallbackActions(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if strings.HasPrefix(update.CallbackQuery.Data, "tvshows:") {
		update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "tvshows:")
		return controllers.SeriesHandleUpdate(c, bot, update)
	}
	if strings.HasPrefix(update.CallbackQuery.Data, "movies:") {
		update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "movies:")
		return controllers.MoviesHandleUpdate(c, bot, update)
	}
	return nil
}
