package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

type Handler func(config *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error
