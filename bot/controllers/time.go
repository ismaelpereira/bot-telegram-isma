package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func TimerHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.Message.Command() == "reminder" {
		if update.Message.CommandArguments() == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgReminder)
			_, err := bot.Send(msg)
			return err
		}
		return nil
	}
	if update.Message.Command() == "now" {
		if strings.ToLower(update.Message.CommandArguments()) == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgNow)
			_, err := bot.Send(msg)
			return err
		}
		return nowHandler(bot, update)

	}
	return nil
}

func nowHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	commandSplited := strings.Fields(strings.ToLower(strings.ToLower(update.Message.CommandArguments())))
	if len(commandSplited) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Você digitou o comando errado. Não foi possível completar a solicitação")
		_, err := bot.Send(msg)
		return err
	}
	operation := commandSplited[0]
	value, err := strconv.Atoi(commandSplited[1])
	if err != nil {
		return err
	}
	measureOfTime := commandSplited[2]
	duration := time.Duration(int64(value))
	var hour time.Time
	if operation == "plus" {
		if measureOfTime == "seconds" {
			hour = time.Now().Add(time.Second * duration)
		}
		if measureOfTime == "minutes" {
			hour = time.Now().Add(time.Minute * duration)
		}
		if measureOfTime == "hours" {
			hour = time.Now().Add(time.Hour * duration)
		}
	}
	if operation == "minus" {
		if measureOfTime == "seconds" {
			hour = time.Now().Add(-time.Second * duration)
		}
		if measureOfTime == "minutes" {
			hour = time.Now().Add(-time.Minute * duration)
		}
		if measureOfTime == "hours" {
			hour = time.Now().Add(-time.Hour * duration)
		}
	}
	if operation != "plus" && operation != "minus" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Você digitou o comando errado. Não foi possível completar a solicitação")
		_, err = bot.Send(msg)
		return err
	}
	time := hour.Format("Monday, 2 January, 2006 - 15:04:05")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, time)
	_, err = bot.Send(msg)
	return err
}
