package controllers

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gomodule/redigo/redis"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

func TimerHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.Message.Command() == "reminder" {
		if update.Message.CommandArguments() == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgReminder)
			_, err := bot.Send(msg)
			return err
		}
		return reminderHandler(bot, update)
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
	value := commandSplited[1]
	measureOfTime := commandSplited[2]
	var hour time.Time
	if operation == "plus" {
		if measureOfTime == "seconds" {
			duration, err := time.ParseDuration(value + "s")
			if err != nil {
				return err
			}
			hour = time.Now().Add(duration).Add(-time.Hour * 3)
		}
		if measureOfTime == "minutes" {
			duration, err := time.ParseDuration(value + "m")
			if err != nil {
				return err
			}
			hour = time.Now().Add(duration).Add(-time.Hour * 3)
		}
		if measureOfTime == "hours" {
			duration, err := time.ParseDuration(value + "h")
			if err != nil {
				return err
			}
			hour = time.Now().Add(duration).Add(-time.Hour * 3)
		}
	}
	if operation == "minus" {
		if measureOfTime == "seconds" {
			duration, err := time.ParseDuration(value + "s")
			if err != nil {
				return err
			}
			hour = time.Now().Add(-time.Second * duration).Add(-time.Hour * 3)
		}
		if measureOfTime == "minutes" {
			duration, err := time.ParseDuration(value + "m")
			if err != nil {
				return err
			}
			hour = time.Now().Add(-time.Minute * duration).Add(-time.Hour * 3)
		}
		if measureOfTime == "hours" {
			duration, err := time.ParseDuration(value + "s")
			if err != nil {
				return err
			}
			hour = time.Now().Add(-time.Hour * duration).Add(-time.Hour * 3)
		}
	}
	if operation != "plus" && operation != "minus" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Você digitou o comando errado. Não foi possível completar a solicitação")
		_, err := bot.Send(msg)
		return err
	}
	time := hour.Format("Monday, 2 January, 2006 - 15:04:05")
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, time)
	_, err := bot.Send(msg)
	return err
}

func reminderHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	//1 minute eat
	commandSplited := strings.SplitAfterN(strings.ToLower(strings.ToLower(update.Message.CommandArguments())), " ", 3)
	if len(commandSplited) < 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Você digitou o comando errado. Não foi possível completar a solicitação")
		_, err := bot.Send(msg)
		return err
	}

	value := strings.TrimSpace(commandSplited[0])
	t, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if t < 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Não é possível criar o lembrete com o tempo negativo")
		_, err := bot.Send(msg)
		return err
	}
	measureOfTime := strings.TrimSpace(commandSplited[1])
	message := commandSplited[2]
	var expireTime time.Time
	if measureOfTime == "seconds" {
		duration, err := time.ParseDuration(value + "s")
		if err != nil {
			return err
		}
		expireTime = time.Now().Add(duration)
	}
	if measureOfTime == "minutes" {
		duration, err := time.ParseDuration(value + "m")
		if err != nil {
			return err
		}
		expireTime = time.Now().Add(duration)
	}
	if measureOfTime == "hours" {
		duration, err := time.ParseDuration(value + "s")
		if err != nil {
			return err
		}
		expireTime = time.Now().Add(duration)
	}
	conn, err := config.StartRedis()
	_, err = conn.Do("HMSET", "telegram:reminder:"+expireTime.Format("2006:01:02:15:04:05"), "chatID", update.Message.Chat.ID, "text", message)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"Lembrete criado com sucesso! \nPara: "+expireTime.Format("2006:01:02:15:04:05")+"\nCom o texto: "+message)
	_, err = bot.Send(msg)
	return err
}

func reminderWorker(bot *tgbotapi.BotAPI) error {
	conn, err := config.StartRedis()
	keys, err := redis.Strings(conn.Do("KEYS", "telegram:reminder:*"))
	if err != nil {
		return err
	}
	spew.Dump(keys)
	if len(keys) != 0 {
		sort.Strings(keys)
		now := "telegram:reminder:" + time.Now().Format("2006:01:02:15:04:05")
		for _, key := range keys {
			if key <= now {
				data, err := redis.StringMap(conn.Do("HGETALL", key))
				if err != nil {
					return err
				}
				if data != nil && data["chatID"] != "" && data["text"] != "" {
					chatID, err := strconv.ParseInt(data["chatID"], 10, 64)
					if err != nil {
						return err
					}
					msg := tgbotapi.NewMessage(chatID, msgs.IconAlarmClock+data["text"])
					_, err = bot.Send(msg)
					if err != nil {
						return err
					}
					conn.Do("DEL", key)
					return nil
				}
			}
		}
	}
	return nil
}

func ReminderCheck(bot *tgbotapi.BotAPI) {
	for {
		reminderWorker(bot)

	}
}
