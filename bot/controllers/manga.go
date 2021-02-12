package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MangaHandleUpdate is a function for manga work
func MangaHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	mangaName := strings.TrimSpace(update.Message.CommandArguments())
	if mangaName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgManga)
		_, err := bot.Send(msg)
		return err
	}
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/manga?q=" + url.QueryEscape(mangaName) + "&page=1&limit=3")
	if err != nil {
		return err
	}
	defer apiResult.Body.Close()
	readMangas, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return err
	}
	var searchResults types.MangaResponse
	err = json.Unmarshal(readMangas, &searchResults)
	if err != nil {
		return err
	}
	for _, manga := range searchResults.Results {
		msgs.GetMangaPictureAndSendMessage(manga, update, bot)
	}
	return nil
}
