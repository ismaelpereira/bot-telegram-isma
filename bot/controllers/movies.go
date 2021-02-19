package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IsmaelPereira/telegram-bot-isma/api/clients"
	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MovieMenu is a map for mantain the results to make the button
var MoviesMenu = make(map[int64][]types.Movie)
var MoviesSearch clients.MovieDB

//MovieHandleUpdate send the movie message
func MoviesHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	if update.CallbackQuery == nil {
		movieName := strings.TrimSpace(update.Message.CommandArguments())
		if movieName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovies)
			_, err := bot.Send(msg)
			return err
		}
		MoviesSearch.ApiKey = c.MovieAcessKey.Key
		moviesResults, err := MoviesSearch.SearchMovie(movieName)
		if err != nil {
			return err
		}
		MoviesMenu[update.Message.Chat.ID] = moviesResults.Results
		if v, ok := MoviesMenu[update.Message.Chat.ID]; ok && len(v) != 0 {
			moviesProviders, err := MoviesSearch.GetMovieProviders(strconv.Itoa(v[0].ID))
			if err != nil {
				return err
			}
			moviesDetails, err := MoviesSearch.GetMovieDetails(strconv.Itoa(v[0].ID))
			if err != nil {
				return err
			}
			moviesCredits, err := MoviesSearch.GetMovieCredits(strconv.Itoa(v[0].ID))
			if err != nil {
				return err
			}
			v[0].Providers = *moviesProviders
			v[0].Details = *moviesDetails
			v[0].Credits = *moviesCredits
			movieMessage, err := getMoviesPictureAndSendMessage(c, bot, update, v[0])
			if err != nil {
				return err
			}
			var kb []tgbotapi.InlineKeyboardMarkup
			if len(v) > 1 {
				kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "movies:1"),
					),
				))
			}
			if len(moviesResults.Results) > 1 {
				movieMessage.ReplyMarkup = kb[0]
			}
			_, err = bot.Send(movieMessage)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return movieArrowButtonsAction(c, bot, update)
}

func movieArrowButtonsAction(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	if v, ok := MoviesMenu[update.CallbackQuery.Message.Chat.ID]; ok && len(v) != 0 {
		moviesProviders, err := MoviesSearch.GetMovieProviders(strconv.Itoa(v[i].ID))
		if err != nil {
			return err
		}
		v[i].Providers = *moviesProviders
		moviesDetails, err := MoviesSearch.GetMovieDetails(strconv.Itoa(v[i].ID))
		if err != nil {
			return err
		}
		v[i].Details = *moviesDetails
		moviesCredits, err := MoviesSearch.GetMovieCredits(strconv.Itoa(v[i].ID))
		v[i].Credits = *moviesCredits
		movieMessage, err := getMoviesPictureAndSendMessage(c, bot, update, v[i])
		if err != nil {
			return err
		}
		var kb []tgbotapi.InlineKeyboardButton
		if i != 0 {
			kb = append(kb,
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "movies:"+strconv.Itoa(i-1)),
			)
		}
		if i != (len(v) - 1) {
			kb = append(kb,
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "movies:"+strconv.Itoa(i+1)),
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
		} else {
			msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + v[i].PosterPath
		}
		messageJSON, err := json.Marshal(msgEdit)
		if err != nil {
			return err
		}

		sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(c.Telegram.Key)+"/editMessageMedia",
			"application/json", bytes.NewBuffer(messageJSON))
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

func getMoviesPictureAndSendMessage(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update, mov types.Movie) (*tgbotapi.PhotoConfig, error) {
	var moviesDetailsMessage []string
	releaseDate, err := time.Parse("2006-01-02", mov.ReleaseDate)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(strconv.Itoa(mov.Details.Duration) + "m")
	if err != nil {
		return nil, err
	}
	moviesDetailsMessage = append(moviesDetailsMessage,
		"\nTítulo: "+mov.Title,
		"\nTítulo Original: "+mov.OriginalTitle,
		"\nPopularidade: "+strconv.FormatFloat(mov.Popularity, 'f', 2, 64),
		"\nData de lançamento: "+releaseDate.Format("02/01/2006"),
		"\nDuração: "+duration.String(),
		"\nNota: "+strconv.FormatFloat(mov.Details.Rating, 'f', 2, 64),
	)
	moviesProvidersMessage := getMovieProviders(mov)
	moviesCreditsMessage := getMovieDirector(mov)
	var movMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		movMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID, "https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+mov.PosterPath)
	}
	if update.CallbackQuery != nil {
		movMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID, "https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+mov.PosterPath)
	}
	movMessage.Caption = strings.Join(moviesDetailsMessage, "") + strings.Join(moviesProvidersMessage, "") + "\nDiretor: " + strings.Join(moviesCreditsMessage, ",")
	return &movMessage, nil
}

func getMovieProviders(mov types.Movie) []string {
	var movProvidersMessage []string
	if country, ok := mov.Providers.Results["BR"]; ok && country != nil {
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
	return movProvidersMessage
}

func getMovieDirector(mov types.Movie) []string {
	directors := make([]string, 0, 2)
	for _, crew := range mov.Credits.Crew {
		if crew.Job == "Director" && crew.Department == "Directing" {
			directors = append(directors, crew.Name)
		}
	}
	// separators := strings.Count(strings.Join(directors, ""), ".")
	// directorString := strings.Replace(strings.Join(directors, ""), ".", ",", separators-1)
	return directors
}
