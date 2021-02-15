package types

import (
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HandlerFunc func(*config.Config, *tgbotapi.BotAPI, *tgbotapi.Update) error
