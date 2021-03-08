package controllers

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

var message []string
var kbComplete tgbotapi.InlineKeyboardMarkup

func ChecklistHandleUpdate(cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		return checkItem(cfg, bot, update)
	}
	listItems := strings.Split(update.Message.CommandArguments(), ",")
	var kb [][]tgbotapi.InlineKeyboardButton
	for i, item := range listItems {
		message = append(message, strconv.Itoa(i+1)+". "+item+msgs.IconX+"\n")
		kb = append(kb, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i)+". "+item, "checklist:"+strconv.Itoa(i)),
		))
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "CHECKLIST\n"+strings.Join(message, ""))
	kbComplete.InlineKeyboard = kb
	msg.ReplyMarkup = kbComplete
	_, err := bot.Send(msg)
	return err
}

func checkItem(cfg *config.Config,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update) error {
	arrayPos, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	itemChecked := strings.Replace(message[arrayPos], msgs.IconX, msgs.IconOk, 1)
	message[arrayPos] = strings.Replace(message[arrayPos], message[arrayPos], itemChecked, 1)
	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, strings.Join(message, ""))
	msg.ReplyMarkup = &kbComplete
	_, err = bot.Send(msg)
	return err
}
