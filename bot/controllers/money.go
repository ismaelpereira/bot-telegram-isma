package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

func MoneyHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	command := strings.ToUpper(update.Message.CommandArguments())
	commandSplit := strings.Fields(command)
	if command == "" {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMoney)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	}
	commandValue := commandSplit[0]
	currencyToConvert := commandSplit[1]
	currencyConverted := commandSplit[2]
	amount, err := strconv.ParseFloat(commandValue, 64)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"Parece que você digitou o comando errado, tente colocar espaços. Ex: '/money 1 usd brl")
		bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(err)
		return err
	}
	var moneyAPI *types.MoneySearchResult
	if temp := apiCache.Get(time.Now()); temp != nil {
		moneyAPI = temp.(*types.MoneySearchResult)
	}
	if moneyAPI == nil {
		apiKey, err := config.GetMoneyApiKey()
		if err != nil {
			log.Println(err)
			return err
		}
		moneyResult, err := http.Get("http://data.fixer.io/api/latest?access_key=" + url.QueryEscape(apiKey))
		if err != nil {
			log.Println(err)
			return err
		}
		defer moneyResult.Body.Close()
		moneyValues, err := ioutil.ReadAll(moneyResult.Body)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(moneyValues, &moneyAPI)
		if err != nil {
			log.Println(err)
			return err
		}
		apiCache.Set(time.Now().Add(time.Hour), moneyAPI)
	}

	if strings.EqualFold(commandSplit[1], "EUR") {
		currency := moneyAPI.Rates[commandSplit[2]] * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if strings.EqualFold(commandSplit[2], "EUR") {
		currency := (1 / moneyAPI.Rates[commandSplit[1]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if strings.EqualFold(commandSplit[1], "EUR") == false && strings.EqualFold(commandSplit[2], "EUR") == false {
		currency := ((1 / moneyAPI.Rates[commandSplit[1]]) * moneyAPI.Rates[commandSplit[2]]) * amount
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, commandValue+" "+currencyToConvert+
			" to "+currencyConverted+" --> "+strconv.FormatFloat(currency, 'f', 2, 64))
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return err
}
