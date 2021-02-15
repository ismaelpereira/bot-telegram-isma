package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/cache"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var apiCache cache.Cache

//MoneyHandleUpdate send the money message
func MoneyHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
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
		if err != nil {
			return err
		}
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
	var moneyAPI *types.MoneySearchResult
	if temp := apiCache.Get(time.Now()); temp != nil {
		moneyAPI = temp.(*types.MoneySearchResult)
		return nil
	}
	if moneyAPI == nil {
		apiKey := c.MoneyAcessKey.Key
		moneyResult, err := http.Get("http://data.fixer.io/api/latest?access_key=" + url.QueryEscape(apiKey))
		if err != nil {
			return err
		}
		defer moneyResult.Body.Close()
		moneyValues, err := ioutil.ReadAll(moneyResult.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(moneyValues, &moneyAPI)
		if err != nil {
			return err
		}
		apiCache.Set(time.Now().Add(time.Hour), moneyAPI)
	}
	if !strings.EqualFold(commandSplit[1], "EUR") && !strings.EqualFold(commandSplit[2], "EUR") {
		currency := ((1 / moneyAPI.Rates[commandSplit[1]]) * moneyAPI.Rates[commandSplit[2]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	if strings.EqualFold(commandSplit[1], "EUR") {
		currency := moneyAPI.Rates[commandSplit[2]] * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	if strings.EqualFold(commandSplit[2], "EUR") {
		currency := (1 / moneyAPI.Rates[commandSplit[1]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		return err
	}
	return nil
}
