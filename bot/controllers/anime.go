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
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAnime)
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}

		return nil
	}
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/anime?q=" + url.QueryEscape(animeName) + "&page=1&limit=3")
	if err != nil {
		log.Println(err)
		return err
	}
	defer apiResult.Body.Close()
	readAnimes, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	var searchResults types.AnimeResponse
	err = json.Unmarshal(readAnimes, &searchResults)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, anime := range searchResults.Results {
		msgs.GetAnimePictureAndSendMessage(anime, update, bot)
	}
	return nil
}
