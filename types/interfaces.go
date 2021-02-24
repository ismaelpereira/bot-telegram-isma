package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

type HandlerFunc func(*config.Config, *tgbotapi.BotAPI, *tgbotapi.Update) error

type HandlerCallback func(*config.Config, *tgbotapi.BotAPI, *tgbotapi.Update) error
