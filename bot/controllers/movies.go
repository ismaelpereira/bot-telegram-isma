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

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

var movies []types.Movie

// MoviesHandleUpdate send the movie message
func MoviesHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery == nil {
		mediaType := update.Message.Command()
		movieName := strings.TrimSpace(update.Message.CommandArguments())
		if movieName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovies)
			_, err := bot.Send(msg)
			return err
		}
		apiKey := cfg.MovieAcessKey.Key
		searchClient, err := clients.NewSearchMedia(mediaType, movieName, apiKey)
		if err != nil {
			return err
		}
		movies, _, err = searchClient.SearchMedia(mediaType, movieName)
		if err != nil {
			return err
		}
		detailsClient, err := clients.NewGetDetails(mediaType, strconv.Itoa(movies[0].ID), apiKey)
		if err != nil {
			return err
		}
		details, _, err := detailsClient.GetDetails(mediaType, strconv.Itoa(movies[0].ID))
		if err != nil {
			return err
		}
		providersClient, err := clients.NewSearchProviders(mediaType, strconv.Itoa(movies[0].ID), apiKey)
		if err != nil {
			return err
		}
		providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(movies[0].ID))
		if err != nil {
			return err
		}
		directorsClient, err := clients.NewGetMovieCredits(strconv.Itoa(movies[0].ID), apiKey)
		if err != nil {
			return err
		}
		credits, err := directorsClient.GetMovieCredits(strconv.Itoa(movies[0].ID))
		if err != nil {
			return err
		}
		movies[0].Details = *details
		movies[0].Providers = *providers
		movies[0].Credits = *credits
		movieMessage, err := getMoviesPictureAndSendMessage(update, movies[0])
		if err != nil {
			return err
		}
		var kb []tgbotapi.InlineKeyboardMarkup
		if len(movies) > 1 {
			kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "movies:1"),
				),
			))
		}
		if len(movies) > 1 {
			movieMessage.ReplyMarkup = kb[0]
		}
		_, err = bot.Send(movieMessage)
		return err
	}
	return movieArrowButtonsAction(cfg, update, movies)
}

func movieArrowButtonsAction(
	cfg *config.Config,
	update *tgbotapi.Update,
	movies []types.Movie,
) error {
	mediaType := "movies"
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	apiKey := cfg.MovieAcessKey.Key
	detailsClient, err := clients.NewGetDetails(mediaType, strconv.Itoa(movies[i].ID), apiKey)
	if err != nil {
		return err
	}
	details, _, err := detailsClient.GetDetails(mediaType, strconv.Itoa(movies[i].ID))
	if err != nil {
		return err
	}
	providersClient, err := clients.NewSearchProviders(mediaType, strconv.Itoa(movies[i].ID), apiKey)
	if err != nil {
		return err
	}
	providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(movies[i].ID))
	if err != nil {
		return err
	}
	directorsClient, err := clients.NewGetMovieCredits(strconv.Itoa(movies[i].ID), apiKey)
	if err != nil {
		return err
	}
	credits, err := directorsClient.GetMovieCredits(strconv.Itoa(movies[i].ID))
	if err != nil {
		return err
	}
	movies[i].Details = *details
	movies[i].Providers = *providers
	movies[i].Credits = *credits
	movieMessage, err := getMoviesPictureAndSendMessage(update, movies[i])
	if err != nil {
		return err
	}
	var kb []tgbotapi.InlineKeyboardButton
	if i != 0 {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "movies:"+strconv.Itoa(i-1)),
		)
	}
	if i != (len(movies) - 1) {
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
	if movies[i].PosterPath == "" {
		msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
	} else {
		msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + movies[i].PosterPath
	}
	messageJSON, err := json.Marshal(msgEdit)
	if err != nil {
		return err
	}
	sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(cfg.Telegram.Key)+"/editMessageMedia",
		"application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer sendMessage.Body.Close()
	if sendMessage.StatusCode > 299 || sendMessage.StatusCode < 200 {
		err = fmt.Errorf("Error in post method %w", err)
		return err
	}
	return nil
}

func getMoviesPictureAndSendMessage(
	update *tgbotapi.Update,
	mov types.Movie,
) (*tgbotapi.PhotoConfig, error) {
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
	return directors
}
