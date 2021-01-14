package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//AnimeHandleUpdate is a function for anime work
func AnimeHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	animeName := strings.TrimSpace(update.Message.CommandArguments())
	if animeName == "" {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAnime)
	}
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/anime?q=" + url.QueryEscape(animeName) + "&page=1&limit=3&type=tv")
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgServerError)
		log.Println(err)
	}
	defer apiResult.Body.Close()
	readAnimes, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgServerError)
		log.Println(err)
	}
	var searchResults types.AnimeResponse
	err = json.Unmarshal(readAnimes, &searchResults)
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgServerError)
		log.Println(err)
	}
	for _, anime := range searchResults.Results {
		msgs.GetAnimePictureAndSendMessage(anime, update, bot)
	}
	return nil
}
