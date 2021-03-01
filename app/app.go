//+build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/ismaelpereira/telegram-bot-isma/bot"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

type App struct {
	Config *config.Config
	Bot    *bot.Bot
}

func Wire(cfg *config.Config, bot *bot.Bot) (*App, error) {
	return &App{
		Config: cfg,
		Bot:    bot,
	}, nil
}

func Build() (*App, error) {
	wire.Build(
		Wire,
		config.Wire,
		bot.Wire,
	)
	return &App{}, nil
}
