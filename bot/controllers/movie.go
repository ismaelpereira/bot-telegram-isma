package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MovieMenu is a map for mantain the results to make the button
var MovieMenu = make(map[int64][]types.MovieDbSearchResults)

//MovieHandleUpdate send the movie message
func MovieHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var movieResults types.MovieResponse
	if update.CallbackQuery == nil {
		movieName := strings.TrimSpace(update.Message.CommandArguments())
		if movieName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovie)
			_, err := bot.Send(msg)
			return err
		}
		apiKey := c.MovieAcessKey.Key
		movieAPI, err := http.Get("https://api.themoviedb.org/3/search/movie?api_key=" + url.QueryEscape(apiKey) +
			"&page=1&langague=pt-br&query=" + url.QueryEscape(movieName))
		if err != nil {
			return err
		}
		defer movieAPI.Body.Close()
		searchValues, err := ioutil.ReadAll(movieAPI.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(searchValues, &movieResults)
		if err != nil {
			return err
		}
		MovieMenu[update.Message.Chat.ID] = movieResults.Results
		if v, ok := MovieMenu[update.Message.Chat.ID]; ok && len(v) != 0 {
			movieMessage, err := getMoviePictureAndSendMessage(c, v[0], update, bot)
			if err != nil {
				return err
			}
			var kb []tgbotapi.InlineKeyboardMarkup
			if len(v) > 1 {
				kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "1"),
					),
				))
			}
			movieMessage.ReplyMarkup = kb[0]
			_, err = bot.Send(movieMessage)
			if err != nil {
				return err
			}
		}
		return nil
	}
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	if v, ok := MovieMenu[update.CallbackQuery.Message.Chat.ID]; ok && len(v) != 0 {
		movieMessage, err := getMoviePictureAndSendMessage(c, v[i], update, bot)
		if err != nil {
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
		if v[i].PosterPath == "" {
			msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
		}
		if v[i].PosterPath != "" {
			msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + v[i].PosterPath
		}

		messageJSON, err := json.Marshal(msgEdit)
		if err != nil {
			return err
		}
		telegramKey := c.Telegram.Key
		sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(telegramKey)+"/editMessageMedia", "application/json", bytes.NewBuffer(messageJSON))
		if err != nil {
			return err
		}
		if sendMessage.StatusCode > 299 && sendMessage.StatusCode < 200 {
			err = fmt.Errorf("Error in post method %w", err)
			return err
		}
	}
	return nil
}

func getMoviePictureAndSendMessage(c *config.Config, mov types.MovieDbSearchResults, update *tgbotapi.Update, bot *tgbotapi.BotAPI) (*tgbotapi.PhotoConfig, error) {
	var movDetailsMessage []string
	releaseDate, err := time.Parse("2006-01-02", mov.ReleaseDate)
	if err != nil {
		return nil, err
	}
	movDetailsMessage = append(movDetailsMessage,
		"\nTítulo: "+mov.Title,
		"\nTítulo Original: "+mov.OriginalTitle,
		"\nPopularidade: "+strconv.FormatFloat(mov.Popularity, 'f', 2, 64),
		"\nData de lançamento: "+releaseDate.Format("02/01/2006"),
	)
	var movPicture *http.Response
	movPicture, err = http.Get("https://themoviedb.org/t/p/w300_and_h450_bestv2" + mov.PosterPath)
	if err != nil {
		return nil, err
	}
	defer movPicture.Body.Close()
	movPictureData, err := ioutil.ReadAll(movPicture.Body)
	movieProvidersMessage, err := getMovieProviders(c, mov)
	if err != nil {
		return nil, err
	}
	var movMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		movMessage = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: movPictureData})
	}
	if update.CallbackQuery != nil {
		movMessage = tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{Bytes: movPictureData})
	}
	movMessage.Caption = strings.Join(movDetailsMessage, "") + strings.Join(movieProvidersMessage, "")
	return &movMessage, nil
}
func getMovieProviders(c *config.Config, mov types.MovieDbSearchResults) (movProvidersMessage []string, err error) {
	apiKey := c.MovieAcessKey.Key

	watchProviders, err := http.Get("https://api.themoviedb.org/3/movie/" +
		url.QueryEscape(strconv.Itoa(mov.ID)) + "/watch/providers?api_key=" +
		url.QueryEscape(apiKey))
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		return nil, err
	}
	defer watchProviders.Body.Close()
	var providers types.WatchProvidersResponse
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		return nil, err
	}
	if country, ok := providers.Results["BR"]; ok && country != nil {
		if country.Buy != nil {
			movProvidersMessage = append(movProvidersMessage, "\nPara Comprar: ")
			for i, providerBuy := range country.Buy {
				movProvidersMessage = append(movProvidersMessage, providerBuy.ProviderName)
				if i == len(country.Buy)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}

		if country.Rent != nil {
			movProvidersMessage = append(movProvidersMessage, "\nPara Alugar: ")
			for i, providerRent := range country.Rent {
				movProvidersMessage = append(movProvidersMessage, providerRent.ProviderName)
				if i == len(country.Rent)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}

		if country.Flatrate != nil {
			movProvidersMessage = append(movProvidersMessage, "\nServicos de streaming: ")
			for i, providerFlatrate := range country.Flatrate {
				movProvidersMessage = append(movProvidersMessage, providerFlatrate.ProviderName)
				if i == len(country.Flatrate)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}
	}
	return movProvidersMessage, err
}
