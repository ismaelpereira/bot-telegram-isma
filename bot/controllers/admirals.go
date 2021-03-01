package controllers

import (
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

var admiralDecoded []types.Admiral
var admiralAPI clients.AdmiralJSON

func init() {
	var cfg *config.Config
	var err error
	cfg, err = config.Wire()
	if err != nil {
		panic(err)
	}
	admiralAPI, err = clients.NewAdmiral(cfg.AdmiralPath.Path)
	if err != nil {
		panic(err)
	}
}

//AdmiralHandleUpdate is a function for admiral work
func AdmiralsHandleUpdate(cfg *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	admiralName := strings.TrimSpace(update.Message.CommandArguments())
	if admiralName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAdmirals)
		_, err := bot.Send(msg)
		return err
	}
	admiralDecoded, err := admiralAPI.GetAdmiral(cfg.AdmiralPath.Path)
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
	adMessage := tgbotapi.NewPhotoShare(update.Message.Chat.ID, ad.ProfilePicture)
	adMessage.Caption = "Nome real: " + ad.RealName +
		"\nNome de almirante: " + ad.AdmiralName +
		"\nIdade: " + strconv.Itoa(ad.Age) +
		"\nData de nascimento: " + ad.BirthDate +
		"\nSigno: " + ad.Sign +
		"\nAltura: " + strconv.FormatFloat(ad.Height, 'f', 2, 64) +
		"\nAkuma no Mi: " + ad.AkumaNoMi +
		"\nAnimal: " + ad.Animal +
		"\nPoder: " + ad.Power +
		"\nInspirado em: " + ad.ActorWhoInspire
	_, err = bot.Send(adMessage)
	return err
}
