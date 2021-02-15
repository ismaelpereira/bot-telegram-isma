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

var SeriesMenu = make(map[int64][]types.SeriesDbSearchResults)

func SeriesHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	var seriesResult types.SeriesResponse
	if update.CallbackQuery == nil {
		serieName := strings.TrimSpace(update.Message.CommandArguments())
		if serieName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgSerie)
			_, err := bot.Send(msg)
			return err
		}
		apiKey := c.MovieAcessKey.Key
		seriesAPI, err := http.Get("https://api.themoviedb.org/3/search/tv?api_key=" + url.QueryEscape(apiKey) +
			"&langague=pt-BR&page=1&query=" + url.QueryEscape(serieName))
		if err != nil {
			return err
		}
		defer seriesAPI.Body.Close()
		searchValues, err := ioutil.ReadAll(seriesAPI.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(searchValues, &seriesResult)
		if err != nil {
			return err
		}
		SeriesMenu[update.Message.Chat.ID] = seriesResult.Results
		if v, ok := SeriesMenu[update.Message.Chat.ID]; ok && len(v) != 0 {
			serieMessage, err := getSeriesPictureAndSendMessage(c, v[0], update, bot)
			if err != nil {
				return err
			}
			var kb []tgbotapi.InlineKeyboardMarkup
			if len(v) > 1 {
				kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "serie:1"),
					),
				))
			}
			if len(seriesResult.Results) > 1 {
				serieMessage.ReplyMarkup = kb[0]
			}
			_, err = bot.Send(serieMessage)
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
	if v, ok := SeriesMenu[update.CallbackQuery.Message.Chat.ID]; ok && len(v) != 0 {
		serieMessage, err := getSeriesPictureAndSendMessage(c, v[i], update, bot)
		if err != nil {
			return err
		}
		var kb []tgbotapi.InlineKeyboardButton
		if i != 0 {
			kb = append(kb,
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "serie:"+strconv.Itoa(i-1)),
			)
		}
		if i != (len(v) - 1) {
			kb = append(kb,
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "serie:"+strconv.Itoa(i+1)),
			)
		}
		var msgEdit types.EditMediaJSON
		msgEdit.ChatID = update.CallbackQuery.Message.Chat.ID
		msgEdit.MessageID = update.CallbackQuery.Message.MessageID
		msgEdit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			kb,
		)
		msgEdit.Media.Caption = serieMessage.Caption
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

func getSeriesPictureAndSendMessage(c *config.Config, ser types.SeriesDbSearchResults, update *tgbotapi.Update, bot *tgbotapi.BotAPI) (*tgbotapi.PhotoConfig, error) {
	var seriesDetailMessage []string
	releaseDate, err := time.Parse("2006-01-02", ser.ReleaseDate)
	if err != nil {
		return nil, err
	}
	seriesDetailMessage = append(seriesDetailMessage,
		"\nTítulo: "+ser.Title,
		"\nTítulo Original: "+ser.OriginalTitle,
		"\nPopularidade: "+strconv.FormatFloat(ser.Popularity, 'f', 2, 64),
		"\nData de lançamento: "+releaseDate.Format("02/01/2006"))
	seriePicture, err := http.Get("https://themoviedb.org/t/p/w300_and_h450_bestv2" + ser.PosterPath)
	if err != nil {
		return nil, err
	}
	defer seriePicture.Body.Close()
	seriePictureData, err := ioutil.ReadAll(seriePicture.Body)
	seriesProviderMessage, err := getSeriesProviders(c, ser)
	if err != nil {
		return nil, err
	}
	var serMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		serMessage = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: seriePictureData})
	}
	if update.CallbackQuery != nil {
		serMessage = tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{Bytes: seriePictureData})
	}
	serMessage.Caption = strings.Join(seriesDetailMessage, "") + strings.Join(seriesProviderMessage, "")
	return &serMessage, err
}

func getSeriesProviders(c *config.Config, ser types.SeriesDbSearchResults) (seriesProvidersMessage []string, err error) {
	apiKey := c.MovieAcessKey.Key
	watchProviders, err := http.Get("https://api.themoviedb.org/3/tv/" + url.QueryEscape(strconv.Itoa(ser.ID)) +
		"/watch/providers?api_key=" + url.QueryEscape(apiKey))
	if err != nil {
		return nil, err
	}
	defer watchProviders.Body.Close()
	var providers types.WatchProvidersResponse
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		return nil, err
	}
	if country, ok := providers.Results["BR"]; ok && country != nil {
		if country.Buy != nil {
			seriesProvidersMessage = append(seriesProvidersMessage, "\nPara Comprar: ")
			for i, providersBuy := range country.Buy {
				seriesProvidersMessage = append(seriesProvidersMessage, providersBuy.ProviderName)
				if i == len(country.Buy)-1 {
					seriesProvidersMessage = append(seriesProvidersMessage, ".")
				} else {
					seriesProvidersMessage = append(seriesProvidersMessage, ",")
				}
			}
		}
		if country.Flatrate != nil {
			seriesProvidersMessage = append(seriesProvidersMessage, "\nServiços de streaming: ")
			for i, providersFlatrate := range country.Flatrate {
				seriesProvidersMessage = append(seriesProvidersMessage, providersFlatrate.ProviderName)
				if i == len(country.Buy)-1 {
					seriesProvidersMessage = append(seriesProvidersMessage, ".")
				} else {
					seriesProvidersMessage = append(seriesProvidersMessage, ",")
				}
			}
		}
	}
	return seriesProvidersMessage, err
}
