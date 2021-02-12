package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
			getAdmiralPictureAndSendMessage(admiral, update, bot)
		}
	}
	return nil
}

func getAdmiralPictureAndSendMessage(ad types.Admiral, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	adPicture, err := http.Get(ad.ProfilePicture)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgServerError)
		_, err := bot.Send(msg)
		return err
	}
	defer adPicture.Body.Close()
	adPictureData, err := ioutil.ReadAll(adPicture.Body)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgNotFound)
		_, err := bot.Send(msg)
		return err

	}

	adMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: adPictureData})
	adMessage.Caption = "Nome real: " + ad.RealName + "\nNome de almirante: " +
		ad.AdmiralName + "\nIdade: " + strconv.Itoa(ad.Age) + "\nData de nascimento: " +
		ad.BirthDate + "\nSigno: " + ad.Sign + "\nAltura: " + strconv.FormatFloat(ad.Height, 'f', 2, 64) +
		"\nAkuma no Mi: " + ad.AkumaNoMi + "\nAnimal: " + ad.Animal + "\nPoder: " + ad.Power + "\nInspirado em: " +
		ad.ActorWhoInspire
	_, err = bot.Send(adMessage)
	return err
}
