package handler

import (
	"log"
	"strings"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/controllers"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

//VerifyAndExecuteCommand is a function to pick the right command and execute the respective function
func VerifyAndExecuteCommand(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
	cmd string,
) error {
	log.Printf("got cmd %q\n", cmd)
	Commands := map[string]types.Handler{
		"help":      controllers.HelpHandlerUpdate,
		"admirals":  controllers.AdmiralsHandleUpdate,
		"animes":    controllers.AnimesHandleUpdate,
		"mangas":    controllers.MangasHandleUpdate,
		"money":     controllers.MoneyHandleUpdate,
		"movies":    controllers.MoviesHandleUpdate,
		"tvshows":   controllers.TVShowHandleUpdate,
		"reminder":  controllers.TimerHandleUpdate,
		"now":       controllers.TimerHandleUpdate,
		"checklist": controllers.ChecklistHandleUpdate,
	}
	if f, ok := Commands[cmd]; ok {
		return f(cfg, redis, bot, update)
	}
	return controllers.NotFoundHandlerUpdate(cfg, redis, bot, update)
}

func CallbackActions(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if strings.HasPrefix(update.CallbackQuery.Data, "tvshows:") {
		update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "tvshows:")
		return controllers.TVShowHandleUpdate(cfg, redis, bot, update)
	}
	if strings.HasPrefix(update.CallbackQuery.Data, "movies:") {
		update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "movies:")
		return controllers.MoviesHandleUpdate(cfg, redis, bot, update)
	}
	if strings.HasPrefix(update.CallbackQuery.Data, "checklist:") {
		update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "checklist:"+update.CallbackQuery.ID+":")
		return controllers.ChecklistHandleUpdate(cfg, redis, bot, update)
	}
	return nil
}
