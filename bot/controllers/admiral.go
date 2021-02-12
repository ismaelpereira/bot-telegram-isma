package controllers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//AdmiralHandleUpdate is a function for admiral work
func AdmiralHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	admiralName := strings.TrimSpace(update.Message.CommandArguments())
	if admiralName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAdmiral)
		_, err := bot.Send(msg)
		return err
	}
	admiralPath, err := config.GetAdmiralPath()
	if err != nil {
		return err
	}
	file, err := os.Open(admiralPath)
	if err != nil {
		return err
	}
	defer file.Close()
	admiralArchive, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var admiralDecoded []types.Admiral
	err = json.Unmarshal(admiralArchive, &admiralDecoded)
	if err != nil {
		return err
	}
	for _, admiral := range admiralDecoded {
		if strings.EqualFold(admiral.AdmiralName, admiralName) || strings.EqualFold(admiral.RealName, admiralName) {
			msgs.GetAdmiralPictureAndSendMessage(admiral, update, bot)
		}
	}
	return nil
}
