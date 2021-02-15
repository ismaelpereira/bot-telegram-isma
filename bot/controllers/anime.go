package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//AnimeHandleUpdate is a function for anime work
func AnimeHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	animeName := strings.TrimSpace(update.Message.CommandArguments())
	if animeName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAnime)
		_, err := bot.Send(msg)
		return err
	}
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/anime?q=" + url.QueryEscape(animeName) +
		"&page=1&limit=3")
	if err != nil {
		return err
	}
	defer apiResult.Body.Close()
	readAnimes, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return err
	}
	var searchResults types.AnimeResponse
	err = json.Unmarshal(readAnimes, &searchResults)
	if err != nil {
		return err
	}
	for _, anime := range searchResults.Results {
		getAnimePictureAndSendMessage(anime, update, bot)
	}
	return nil
}

func getAnimePictureAndSendMessage(an types.Anime, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	anPicture, err := http.Get(an.CoverPicture)
	if err != nil {
		return err
	}
	defer anPicture.Body.Close()
	anPictureData, err := ioutil.ReadAll(anPicture.Body)
	if err != nil {
		return err
	}
	anMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: anPictureData})
	var airing string
	if an.Airing == true {
		airing = "Sim"
	} else {
		airing = "Não"
	}
	animeEpisodes := strconv.Itoa(an.Episodes)
	if animeEpisodes == "0" {
		animeEpisodes = "?"
	}
	anMessage.Caption = "Título: " + an.Title + "\nNota: " + strconv.FormatFloat(an.Score, 'f', 2, 64) +
		"\nEpisódios: " + animeEpisodes + "\nPassando? " + airing
	_, err = bot.Send(anMessage)
	return err
}
