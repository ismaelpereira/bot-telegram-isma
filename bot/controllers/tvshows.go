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

var tvShows []types.TVShow

func SeriesHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery == nil {
		mediaType := update.Message.Command()
		tvShowName := strings.TrimSpace(update.Message.CommandArguments())
		if tvShowName == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgTVShow)
			_, err := bot.Send(msg)
			return err
		}
		apiKey := cfg.MovieAcessKey.Key
		searchClient, err := clients.NewSearchMedia(mediaType, tvShowName, apiKey)
		if err != nil {
			return err
		}
		_, tvShows, err = searchClient.SearchMedia(mediaType, tvShowName)
		if err != nil {
			return err
		}
		detailsClient, err := clients.NewGetDetails(mediaType, strconv.Itoa(tvShows[0].ID), apiKey)
		if err != nil {
			return err
		}
		_, details, err := detailsClient.GetDetails(mediaType, strconv.Itoa(tvShows[0].ID))
		if err != nil {
			return err
		}
		providersClient, err := clients.NewSearchProviders(mediaType, strconv.Itoa(tvShows[0].ID), apiKey)
		if err != nil {
			return err
		}
		providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(tvShows[0].ID))
		if err != nil {
			return err
		}
		tvShows[0].TVShowDetails = *details
		tvShows[0].Providers = *providers
		tvShowMessage, err := getTVShowPictureAndSendMessage(cfg, bot, update, tvShows[0])
		if err != nil {
			return err
		}
		var kb []tgbotapi.InlineKeyboardMarkup
		if len(tvShows) > 1 {
			kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "tvshows:1"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Detalhes", "tvshows:seasons:0"),
				),
			))
		}
		if len(tvShows) > 1 {
			tvShowMessage.ReplyMarkup = kb[0]
		}
		if len(tvShows) == 1 {
			tvShowMessage.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Detalhes", "tvshows:seasons:0"),
				),
			)
		}
		_, err = bot.Send(tvShowMessage)
		return err
	}
	if strings.HasPrefix(update.CallbackQuery.Data, "seasons:") {
		return hasCallback(cfg, update, tvShows)
	}
	return tvShowArrowButtonsAction(cfg, bot, update, tvShows)
}

func hasCallback(cfg *config.Config, update *tgbotapi.Update, tvShows []types.TVShow) error {
	update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "seasons:")
	if strings.Contains(update.CallbackQuery.Data, ":") {
		lastBin := strings.Index(update.CallbackQuery.Data, ":")
		arrayPos, err := strconv.Atoi(update.CallbackQuery.Data[:lastBin])
		if err != nil {
			return err
		}
		season, err := strconv.Atoi(strings.TrimPrefix(update.CallbackQuery.Data[lastBin:], ":"))
		if err != nil {
			return err
		}
		var seasonDetails []string
		releaseDate, err := time.Parse("2006-01-02", tvShows[arrayPos].TVShowDetails.Seasons[season].AirDate)
		if err != nil {
			return err
		}
		seasonDetails = append(seasonDetails,
			"\nNúmero de Episódios: ", strconv.Itoa(tvShows[arrayPos].TVShowDetails.Seasons[season].EpisodesCount),
			"\nData de Lançamento: ", releaseDate.Format("02/06/2006"),
		)
		var msgEdit types.EditMediaJSON
		msgEdit.ChatID = update.CallbackQuery.Message.Chat.ID
		msgEdit.MessageID = update.CallbackQuery.Message.MessageID
		msgEdit.Media.Type = "photo"
		if tvShows[arrayPos].TVShowDetails.Seasons[season].PosterPath == "" {
			msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
		} else {
			msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2/" + tvShows[arrayPos].TVShowDetails.Seasons[season].PosterPath
		}
		msgEdit.Media.Caption = strings.Join(seasonDetails, "")
		msgEdit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Voltar", "tvshows:seasons:"+strconv.Itoa(arrayPos)),
			),
		)
		messageJSON, err := json.Marshal(msgEdit)
		if err != nil {
			return err
		}
		sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(cfg.Telegram.Key)+"/editmessagemedia",
			"application/json", bytes.NewBuffer(messageJSON))
		if err != nil {
			return err
		}
		defer sendMessage.Body.Close()
		if sendMessage.StatusCode < 200 && sendMessage.StatusCode > 299 {
			err = fmt.Errorf("Error in post method %w", err)
			return err
		}

		return nil
	}
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	if len(tvShows) != 0 {
		return sendSeasonKeyboard(cfg, update, tvShows[i])
	}
	return nil
}

func tvShowArrowButtonsAction(
	cfg *config.Config,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
	tvShows []types.TVShow,
) error {
	mediaType := "tvshows"
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	apiKey := cfg.MovieAcessKey.Key
	detailsClient, err := clients.NewGetDetails(mediaType, strconv.Itoa(tvShows[i].ID), apiKey)
	if err != nil {
		return err
	}
	_, details, err := detailsClient.GetDetails(mediaType, strconv.Itoa(tvShows[i].ID))
	if err != nil {
		return err
	}
	providersClient, err := clients.NewSearchProviders(mediaType, strconv.Itoa(tvShows[i].ID), apiKey)
	if err != nil {
		return err
	}
	providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(tvShows[i].ID))
	if err != nil {
		return err
	}
	tvShows[i].TVShowDetails = *details
	tvShows[i].Providers = *providers
	tvShowMessage, err := getTVShowPictureAndSendMessage(cfg, bot, update, tvShows[i])
	if err != nil {
		return err
	}
	var kb []tgbotapi.InlineKeyboardButton
	if i != 0 {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "tvshows:"+strconv.Itoa(i-1)),
		)
	}
	if i != (len(tvShows) - 1) {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "tvshows:"+strconv.Itoa(i+1)),
		)
	}
	var msgEdit types.EditMediaJSON
	msgEdit.ChatID = update.CallbackQuery.Message.Chat.ID
	msgEdit.MessageID = update.CallbackQuery.Message.MessageID
	msgEdit.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		kb,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Detalhes", "tvshows:seasons:"+strconv.Itoa(i)),
		),
	)
	msgEdit.Media.Caption = tvShowMessage.Caption
	msgEdit.Media.Type = "photo"
	if tvShows[i].PosterPath == "" {
		msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
	} else {
		msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + tvShows[i].PosterPath
	}
	messageJSON, err := json.Marshal(msgEdit)
	if err != nil {
		return err
	}
	sendMessage, err := http.Post("https://api.telegram.org/bot"+url.QueryEscape(cfg.Telegram.Key)+"/editmessagemedia",
		"application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return err
	}
	defer sendMessage.Body.Close()
	if sendMessage.StatusCode > 299 && sendMessage.StatusCode < 200 {
		err = fmt.Errorf("Error in post method %w", err)
		return err
	}

	return nil
}

func sendSeasonKeyboard(cfg *config.Config,
	update *tgbotapi.Update,
	tvShow types.TVShow,
) error {
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	var kb [][]tgbotapi.InlineKeyboardButton
	for s, season := range tvShow.TVShowDetails.Seasons {
		kb = append(kb, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(season.Name, "tvshows:seasons:"+strconv.Itoa(i)+":"+strconv.Itoa(s)),
		))
	}
	kb = append(kb, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Voltar", "tvshows:"+strconv.Itoa(i)),
	))
	var kbComplete tgbotapi.InlineKeyboardMarkup
	kbComplete.InlineKeyboard = kb
	var msgEdit types.EditMediaJSON
	msgEdit.ChatID = update.CallbackQuery.Message.Chat.ID
	msgEdit.MessageID = update.CallbackQuery.Message.MessageID
	msgEdit.ReplyMarkup = kbComplete
	msgEdit.Media.Type = "photo"
	if tvShow.PosterPath == "" {
		msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
	} else {
		msgEdit.Media.URL = "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" + tvShow.PosterPath
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
	if sendMessage.StatusCode > 299 && sendMessage.StatusCode < 200 {
		err = fmt.Errorf("Error in post method %w", err)
		return err
	}
	return nil
}

func getTVShowPictureAndSendMessage(
	cfg *config.Config,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
	tvShow types.TVShow,
) (*tgbotapi.PhotoConfig, error) {
	var tvShowDetailsMessage []string
	releaseDate, err := time.Parse("2006-01-02", tvShow.ReleaseDate)
	if err != nil {
		return nil, err
	}
	tvShowDetailsMessage = append(tvShowDetailsMessage,
		"\nTítulo: "+tvShow.Title,
		"\nTítulo Original: "+tvShow.OriginalTitle,
		"\nPopularidade: "+strconv.FormatFloat(tvShow.Popularity, 'f', 2, 64),
		"\nData de lançamento: "+releaseDate.Format("02/06/2006"),
		"\nNota: "+strconv.FormatFloat(tvShow.TVShowDetails.Rating, 'f', 2, 64),
	)
	tvShowSeasonsDetails := getTVShowSeasonDetails(tvShow)
	tvShowProvidersMessage := getTVShowProviders(tvShow)
	tvShowDirector := getTvShowDirector(tvShow)
	var tvShowMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		tvShowMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID, "https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShow.PosterPath)
	}
	if update.CallbackQuery != nil {
		tvShowMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID, "https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShow.PosterPath)
	}
	tvShowMessage.Caption = strings.Join(tvShowDetailsMessage, "") + strings.Join(tvShowSeasonsDetails, "") + strings.Join(tvShowProvidersMessage, "") + "\nDiretores: " + strings.Join(tvShowDirector, ",")
	return &tvShowMessage, nil
}

func getTVShowSeasonDetails(tvShow types.TVShow) []string {
	var seriesSeasonDetails []string
	seriesSeasonDetails = append(seriesSeasonDetails,
		"\nNúmero de temporadas: "+strconv.Itoa(tvShow.TVShowDetails.SeasonNumber),
		"\nStatus: "+tvShow.TVShowDetails.Status)
	return seriesSeasonDetails
}

func getTVShowProviders(tvShow types.TVShow) []string {
	var tvShowProvidersMessage []string
	if country, ok := tvShow.Providers.Results["BR"]; ok && country != nil {
		if country.Buy != nil {
			tvShowProvidersMessage = append(tvShowProvidersMessage, "\nPara Comprar: ")
			for i, providerBuy := range country.Buy {
				tvShowProvidersMessage = append(tvShowProvidersMessage, providerBuy.ProviderName)
				if i == len(country.Buy)-1 {
					tvShowProvidersMessage = append(tvShowProvidersMessage, ".")
				} else {
					tvShowProvidersMessage = append(tvShowProvidersMessage, ", ")
				}
			}
		}
		if country.Flatrate != nil {
			tvShowProvidersMessage = append(tvShowProvidersMessage, "\nServicos de streaming: ")
			for i, providerFlatrate := range country.Flatrate {
				tvShowProvidersMessage = append(tvShowProvidersMessage, providerFlatrate.ProviderName)
				if i == len(country.Flatrate)-1 {
					tvShowProvidersMessage = append(tvShowProvidersMessage, ".")
				} else {
					tvShowProvidersMessage = append(tvShowProvidersMessage, ", ")
				}
			}
		}
	}
	return tvShowProvidersMessage
}

func getTvShowDirector(tvShow types.TVShow) []string {
	var directors []string
	for _, director := range tvShow.TVShowDetails.CreatedBy {
		directors = append(directors, director.Name)
	}
	if len(tvShow.TVShowDetails.CreatedBy) == 0 {
		directors = append(directors, "-")
	}
	return directors
}
