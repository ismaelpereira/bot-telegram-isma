package bot

import (
	"log"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

func Wire(
	cfg *config.Config,
	redis *redis.Client,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.Telegram.Key)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	bot := &Bot{
		API: api,
	}
	return bot, nil
}
