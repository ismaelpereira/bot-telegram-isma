package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/controllers"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/handler"
)

type Bot struct {
	API    *tgbotapi.BotAPI
	Update *tgbotapi.Update
	wg     sync.WaitGroup
}

func (t *Bot) Start() {
	t.wg.Add(4)
	go t.updateSink()
	go controllers.ReminderCheck(t.API)
}

func (t *Bot) updateSink() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := t.API.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}
	cfg, err := config.Wire()
	if err != nil {
		panic(err)
	}
	log.Println("bot started")
	for {
		select {
		case update := <-updates:
			{
				if err := t.handle(cfg, &update); err != nil {
					if err != nil {
						fmt.Errorf("%w", err)
					}
				}
			}
		}
	}
}

func (t *Bot) Workers() {}

func (t *Bot) handle(cfg *config.Config, update *tgbotapi.Update) (err error) {

	if update.CallbackQuery != nil {
		log.Printf("got callback query\n")
		return handler.CallbackActions(cfg, t.API, update)
	}
	if update.Message == nil || !update.Message.IsCommand() {
		return nil
	}
	return handler.VerifyAndExecuteCommand(cfg, t.API, update, strings.ToLower(update.Message.Command()))

}

func (t *Bot) Wait() error {
	t.wg.Wait()
	return nil
}
