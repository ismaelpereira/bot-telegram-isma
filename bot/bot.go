package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/controllers"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/handler"
	r "github.com/ismaelpereira/telegram-bot-isma/redis"
)

type Bot struct {
	API    *tgbotapi.BotAPI
	Update *tgbotapi.Update
}

func (t *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.API.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	cfg, err := config.Wire()
	conn, err := r.Wire(cfg)
	if err != nil {
		panic(err)
	}
	go controllers.ReminderCheck(t.API, conn)
	if err != nil {
		panic(err)
	}
	log.Println("bot started")
	for update := range updates {
		if err := t.handle(cfg, conn, &update); err != nil {
			if err != nil {
				fmt.Errorf("%w", err)
			}
		}
	}
}

func (t *Bot) handle(cfg *config.Config,
	redis *redis.Client,
	update *tgbotapi.Update,
) (err error) {
	if update.CallbackQuery != nil {
		log.Printf("got callback query\n")
		return handler.CallbackActions(cfg, redis, t.API, update)
	}
	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}
	return handler.VerifyAndExecuteCommand(cfg, redis, t.API, update, strings.ToLower(update.Message.Command()))

}
