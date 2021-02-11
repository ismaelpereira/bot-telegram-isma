package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var MovieMenu = make(map[int64][]types.MovieDbSearchResults)

func MovieHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var movieResults types.MovieResponse

	if update.CallbackQuery == nil {
		movieName := strings.TrimSpace(update.Message.CommandArguments())
		if movieName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovie)
			_, err := bot.Send(msg)
			if err != nil {
				return err
			}
			return err
		}
		apiKey, err := config.GetMovieApiKey()
		if err != nil {
			return err
		}
		movieApi, err := http.Get("https://api.themoviedb.org/3/search/movie?api_key=" + url.QueryEscape(apiKey) +
			"&page=1&langague=pt-br&query=" + url.QueryEscape(movieName))
		if err != nil {
			return err
		}
		defer movieApi.Body.Close()
		searchValues, err := ioutil.ReadAll(movieApi.Body)
		if err != nil {
			log.Println(err)
			return err
		}
		err = json.Unmarshal(searchValues, &movieResults)
		if err != nil {
			log.Println(err)
			return err
		}

		MovieMenu[update.Message.Chat.ID] = movieResults.Results

		var i int
		if v, ok := MovieMenu[update.Message.Chat.ID]; ok && len(v) != 0 {
			movieMessage, err := msgs.GetMoviePictureAndSendMessage(v[i], update, bot)
			if err != nil {
				log.Println(err)
				return err
			}
			var kb []tgbotapi.InlineKeyboardMarkup
			if i != 0 {
				kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, strconv.Itoa(i-1)),
					),
				))
			}
			if i != len(v)-1 {
				kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, strconv.Itoa(i+1)),
					),
				))
			}
			movieMessage.ReplyMarkup = kb[0]
			_, err = bot.Send(movieMessage)
			if err != nil {
				log.Println(err)
				return err
			}

		}
		return nil
	}
	if update.CallbackQuery != nil {
		i, _ := strconv.Atoi(update.CallbackQuery.Data)
		if v, ok := MovieMenu[update.CallbackQuery.Message.Chat.ID]; ok && len(v) != 0 {
			movieMessage, err := msgs.GetMoviePictureAndSendMessage(v[i], update, bot)
			if err != nil {
				log.Println(err)
				return err
			}
			var kb []tgbotapi.InlineKeyboardButton
			if i != 0 {
				kb = append(kb,
					tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, strconv.Itoa(i-1)),
				)
			}
			if i != (len(v) - 1) {
				kb = append(kb,
					tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, strconv.Itoa(i+1)),
				)
			}
			var msgEdit types.EditMediaJSON
			msgEdit.ChatID = update.CallbackQuery.Message.Chat.ID
			msgEdit.MessageID = update.CallbackQuery.Message.MessageID
			msgEdit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				kb,
			)
			msgEdit.Media.Caption = movieMessage.Caption
			msgEdit.Media.Type = "photo"
			msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + v[i].PosterPath
			messageJSON, err := json.Marshal(msgEdit)
			if err != nil {
				log.Println(err)
				return err
			}
			telegramKey, err := config.GetTelegramKey()
			if err != nil {
				log.Println(err)
				return err
			}
			sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(telegramKey)+"/editMessageMedia", "application/json", bytes.NewBuffer(messageJSON))
			if err != nil {
				log.Println(err)
				return err
			}

			if sendMessage.StatusCode > 299 && sendMessage.StatusCode < 200 {
				err = fmt.Errorf("Error in post method %w", err)
				return err
			}
			return nil
		}
	}
	return nil
}
