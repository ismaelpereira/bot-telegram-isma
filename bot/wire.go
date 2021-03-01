package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

func Wire(c *config.Config) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(c.Telegram.Key)
	if err != nil {
		return nil, err
	}
	bot := &Bot{
		API: api,
	}
	return bot, nil
}
