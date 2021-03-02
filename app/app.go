//+build wireinject

package app

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
	"github.com/ismaelpereira/telegram-bot-isma/bot"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	r "github.com/ismaelpereira/telegram-bot-isma/redis"
)

type App struct {
	Config *config.Config
	Redis  *redis.Client
	Bot    *bot.Bot
}

func Wire(
	cfg *config.Config,
	client *redis.Client,
	bot *bot.Bot,
) (*App, error) {
	return &App{
		Config: cfg,
		Redis:  client,
		Bot:    bot,
	}, nil
}

func Build() (*App, error) {
	wire.Build(
		Wire,
		config.Wire,
		r.Wire,
		bot.Wire,
	)
	return &App{}, nil
}
