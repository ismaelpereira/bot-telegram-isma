package controllers

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

var moneyAPI clients.MoneyAPI

func init() {
	var cfg *config.Config
	var err error
	cfg, err = config.Wire()
	if err != nil {
		panic(err)
	}
	moneyAPI, err = clients.NewMoneyAPI(cfg.MoneyAcessKey.Key)
	if err != nil {
		panic(err)
	}
}

// MoneyHandleUpdate send the money message
func MoneyHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	command := strings.ToUpper(update.Message.CommandArguments())
	commandSplit := strings.Fields(command)
	if command == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMoney)
		_, err := bot.Send(msg)
		return err
	}
	if len(commandSplit) != 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Você digitou o comando errado. Não foi possível completar a solicitação")
		_, err := bot.Send(msg)
		return err
	}
	commandValue := commandSplit[0]
	currencyToConvert := commandSplit[1]
	currencyConverted := commandSplit[2]
	if currencyToConvert == currencyConverted {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			commandValue+" "+currencyConverted+" to "+currencyConverted+" --> "+commandValue)
		_, err := bot.Send(msg)
		return err
	}
	amount, err := strconv.ParseFloat(commandValue, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Parece que você digitou o comando errado, tente colocar espaços. Ex: '/money 1 usd brl")
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}
	moneyResults, err := moneyAPI.GetCurrencies()
	if err != nil {
		return err
	}

	if !strings.EqualFold(commandSplit[1], "EUR") && !strings.EqualFold(commandSplit[2], "EUR") {
		currency := ((1 / moneyResults.Rates[commandSplit[1]]) * moneyResults.Rates[commandSplit[2]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	if strings.EqualFold(commandSplit[1], "EUR") {
		currency := moneyResults.Rates[commandSplit[2]] * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	if strings.EqualFold(commandSplit[2], "EUR") {
		currency := (1 / moneyResults.Rates[commandSplit[1]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	return nil
}
