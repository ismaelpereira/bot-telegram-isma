package handler

import (
	"log"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/controllers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//var Commands []string = []string{"help", "admiral", "anime", "manga", "money", "movie"}
// var Commands map[string]interface{} = map[string]interface{}{
// 	"help": controllers.HelpHandlerUpdate(bot, update),
// }

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
			log.Println(err)
			return err
		}
	}
	return nil
}
