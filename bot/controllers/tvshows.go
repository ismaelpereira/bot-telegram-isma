package controllers

import (
	"fmt"
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

func TVShowHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, "seasons:") {
			update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "seasons:")
			return hasCallback(cfg, update, tvShows)
		}
		return tvShowArrowButtonsAction(cfg, update, tvShows)
	}
	mediaType := update.Message.Command()
	tvShowName := strings.TrimSpace(update.Message.CommandArguments())
	if tvShowName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgTVShow)
		_, err := bot.Send(msg)
		return err
	}
	var err error
	tvShows, err = callTVShowsfunc(cfg, update, mediaType, tvShowName)
	if err != nil {
		return err
	}
	tvShowMessage, err := getTVShowPictureAndSendMessage(update, tvShows[0])
	if err != nil {
		return err
	}
	kb := SendTVShowsKeyboard(tvShows)
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

func hasCallback(
	cfg *config.Config,
	update *tgbotapi.Update,
	tvShows []types.TVShow,
) error {
	if !strings.Contains(update.CallbackQuery.Data, ":") {
		i, err := strconv.Atoi(update.CallbackQuery.Data)
		if err != nil {
			return err
		}
		if len(tvShows) != 0 {
			return sendSeasonKeyboard(cfg, update, tvShows[i])
		}
	}
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
		"\nData de Lançamento: ", releaseDate.Format("02/01/2006"),
	)
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Voltar", "tvshows:seasons:"+strconv.Itoa(arrayPos)),
		),
	)
	err = msgs.EditMessage(cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShows[arrayPos].TVShowDetails.Seasons[season].PosterPath,
		strings.Join(seasonDetails, ""),
		kb)
	return err
}

func tvShowArrowButtonsAction(
	cfg *config.Config,
	update *tgbotapi.Update,
	tvShows []types.TVShow,
) error {
	mediaType := "tvshows"
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	tvShows, err = callTVShowsfunc(cfg, update, mediaType, tvShows[i].Title)
	if err != nil {
		return err
	}
	tvShowMessage, err := getTVShowPictureAndSendMessage(update, tvShows[i])
	if err != nil {
		return err
	}
	kb := SendTVShowsCallbackKeyboard(tvShows, i)

	kbDetails := tgbotapi.NewInlineKeyboardMarkup(
		kb,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Detalhes", "tvshows:seasons:"+strconv.Itoa(i)),
		),
	)
	err = msgs.EditMessage(
		cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShows[i].PosterPath,
		tvShowMessage.Caption,
		kbDetails,
	)
	return err
}

func sendSeasonKeyboard(
	cfg *config.Config,
	update *tgbotapi.Update,
	tvShow types.TVShow,
) error {
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	kb := make([][]tgbotapi.InlineKeyboardButton, 0, len(tvShow.TVShowDetails.Seasons))
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
	err = msgs.EditMessage(
		cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShow.PosterPath,
		"",
		kbComplete,
	)
	return err
}

func getTVShowPictureAndSendMessage(
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
		"\nData de lançamento: "+releaseDate.Format("02/01/2006"),
		"\nNota: "+strconv.FormatFloat(tvShow.TVShowDetails.Rating, 'f', 2, 64),
	)
	tvShowSeasonsDetails := getTVShowSeasonDetails(tvShow)
	tvShowProvidersMessage := getTVShowProviders(tvShow)
	tvShowDirector := getTvShowDirector(tvShow)
	var tvShowMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		tvShowMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID,
			"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShow.PosterPath)
	}
	if update.CallbackQuery != nil {
		tvShowMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID,
			"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+tvShow.PosterPath)
	}
	tvShowMessage.Caption = strings.Join(tvShowDetailsMessage, "") +
		strings.Join(tvShowSeasonsDetails, "") + strings.Join(tvShowProvidersMessage, "") +
		"\nDiretores: " + strings.Join(tvShowDirector, ",")
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
		tvShowProvidersMessage = append(tvShowProvidersMessage, "\nPara Comprar: ")
		for i, providerBuy := range country.Buy {
			tvShowProvidersMessage = append(tvShowProvidersMessage, providerBuy.ProviderName)
			if i == len(country.Buy)-1 {
				tvShowProvidersMessage = append(tvShowProvidersMessage, ".")
			} else {
				tvShowProvidersMessage = append(tvShowProvidersMessage, ", ")
			}
		}
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
	return tvShowProvidersMessage
}

func getTvShowDirector(tvShow types.TVShow) []string {
	directors := make([]string, 0, len(tvShow.TVShowDetails.CreatedBy))
	for _, director := range tvShow.TVShowDetails.CreatedBy {
		directors = append(directors, director.Name)
	}
	if len(tvShow.TVShowDetails.CreatedBy) == 0 {
		directors = append(directors, "-")
	}
	return directors
}

func callTVShowsfunc(
	cfg *config.Config,
	update *tgbotapi.Update,
	mediaType string,
	mediaTitle string,
) ([]types.TVShow, error) {
	var arrayPos int
	var err error
	var res interface{}
	apiKey := cfg.MovieAcessKey.Key
	if update.CallbackQuery == nil {
		searchClient, err := clients.NewSearchMedia(mediaType, mediaTitle, apiKey)
		if err != nil {
			return nil, err
		}
		res, err = searchClient.SearchMedia(mediaType, mediaTitle)
		if err != nil {
			return nil, err
		}
		tvShows = res.([]types.TVShow)
		if len(tvShows) == 0 {
			err = fmt.Errorf("No TVShows results: %w", err)
			return nil, err
		}
	} else if arrayPos, err = strconv.Atoi(update.CallbackQuery.Data); err != nil {
		return nil, err
	}
	detailsClient, err := clients.NewGetDetails(mediaType, strconv.Itoa(tvShows[arrayPos].ID), apiKey)
	if err != nil {
		return nil, err
	}
	res, err = detailsClient.GetDetails(mediaType, strconv.Itoa(tvShows[arrayPos].ID))
	if err != nil {
		return nil, err
	}
	details := res.(*types.TVShowDetails)
	providersClient, err := clients.NewSearchProviders(mediaType, strconv.Itoa(tvShows[arrayPos].ID), apiKey)
	if err != nil {
		return nil, err
	}
	providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(tvShows[arrayPos].ID))
	if err != nil {
		return nil, err
	}
	tvShows[arrayPos].TVShowDetails = *details
	tvShows[arrayPos].Providers = *providers
	return tvShows, nil
}
