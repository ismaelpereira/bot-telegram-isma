package types

import (
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

type Handler func(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error

type RedisHandler func(
	[]byte,
) (interface{}, error)

type MediaTypeHandler func() (interface{}, error)
