package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func MovieHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	movieName := strings.TrimSpace(update.Message.CommandArguments())
	if movieName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovie)
		_, err := bot.Send(msg)
		if err != nil {
			return err
		}
		return nil
	}
	apiKey, err := config.GetMovieApiKey()
	if err != nil {
		return err
	}
	movieApi, err := http.Get("https://api.themoviedb.org/3/search/movie?api_key=" + url.QueryEscape(apiKey) + "&page=1&langague=pt-br&query=" + url.QueryEscape(movieName))
	if err != nil {
		return err
	}
	defer movieApi.Body.Close()
	searchValues, err := ioutil.ReadAll(movieApi.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	var movieResults types.MovieResponse
	err = json.Unmarshal(searchValues, &movieResults)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, movie := range movieResults.Results {
		msgs.GetMoviePictureAndSendMessage(movie, update, bot)
	}
	return nil
}
