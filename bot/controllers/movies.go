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

var movies []types.Movie
var searchClient clients.SearchMedia
var detailsClient clients.GetDetails
var providersClient clients.SearchProviders
var directorsClient clients.GetMovieCredits

func init() {
	var cfg *config.Config
	var err error
	cfg, err = config.Wire()
	if err != nil {
		panic(err)
	}
	mediaType := "movie"
	apiKey := cfg.MovieAcessKey.Key
	searchClient, err = clients.NewSearchMedia(mediaType, apiKey)
	if err != nil {
		panic(err)
	}
	detailsClient, err = clients.NewGetDetails(mediaType, apiKey)
	if err != nil {
		panic(err)
	}
	providersClient, err = clients.NewSearchProviders(mediaType, apiKey)
	if err != nil {
		panic(err)
	}
	directorsClient, err = clients.NewGetMovieCredits(apiKey)
	if err != nil {
		panic(err)
	}
}

// MoviesHandleUpdate send the movie message
func MoviesHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		return movieArrowButtonsAction(cfg, update, movies)
	}
	mediaType := update.Message.Command()
	movieName := strings.TrimSpace(update.Message.CommandArguments())
	if movieName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMovies)
		_, err := bot.Send(msg)
		return err
	}
	var err error
	movies, err = callMovieFunctions(update, mediaType, movieName)
	if err != nil {
		return err
	}
	movieMessage, err := getMoviesPictureAndSendMessage(update, movies[0])
	if err != nil {
		return err
	}
	kb := SendMoviesKeyboard(movies)
	if len(movies) > 1 {
		movieMessage.ReplyMarkup = kb[0]
	}
	_, err = bot.Send(movieMessage)
	return err
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
	movies, err = callMovieFunctions(update, mediaType, movies[i].Title)
	if err != nil {
		return err
	}
	movieMessage, err := getMoviesPictureAndSendMessage(update, movies[i])
	if err != nil {
		return err
	}
	kb := SendMoviesCallbackKeyboard(movies, i)
	err = msgs.EditMessage(
		cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+movies[i].PosterPath,
		movieMessage.Caption,
		tgbotapi.NewInlineKeyboardMarkup(
			kb,
		),
	)
	return err
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
		movMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID,
			"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+mov.PosterPath)
	}
	if update.CallbackQuery != nil {
		movMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID,
			"https://www.themoviedb.org/t/p/w300_and_h450_bestv2"+mov.PosterPath)
	}
	movMessage.Caption = strings.Join(moviesDetailsMessage, "") +
		strings.Join(moviesProvidersMessage, "") + "\nDiretor: " +
		strings.Join(moviesCreditsMessage, ",")
	return &movMessage, nil
}

func getMovieProviders(mov types.Movie) []string {
	country := mov.Providers.Results["BR"]
	if country == nil {
		return nil
	}
	movProvidersMessage := make([]string, 0, len(country.Buy)+len(country.Rent)+len(country.Flatrate))
	movProvidersMessage = append(movProvidersMessage, "\nPara Comprar: ")
	for i, providerBuy := range country.Buy {
		movProvidersMessage = append(movProvidersMessage, providerBuy.ProviderName)
		if i == len(country.Buy)-1 {
			movProvidersMessage = append(movProvidersMessage, ".")
		} else {
			movProvidersMessage = append(movProvidersMessage, ", ")
		}
	}
	movProvidersMessage = append(movProvidersMessage, "\nPara Alugar: ")
	for i, providerRent := range country.Rent {
		movProvidersMessage = append(movProvidersMessage, providerRent.ProviderName)
		if i == len(country.Rent)-1 {
			movProvidersMessage = append(movProvidersMessage, ".")
		} else {
			movProvidersMessage = append(movProvidersMessage, ", ")
		}
	}
	movProvidersMessage = append(movProvidersMessage, "\nServicos de streaming: ")
	for i, providerFlatrate := range country.Flatrate {
		movProvidersMessage = append(movProvidersMessage, providerFlatrate.ProviderName)
		if i == len(country.Flatrate)-1 {
			movProvidersMessage = append(movProvidersMessage, ".")
		} else {
			movProvidersMessage = append(movProvidersMessage, ", ")
		}
	}
	return movProvidersMessage
}

func getMovieDirector(mov types.Movie) []string {
	directors := make([]string, 0, len(mov.Credits.Crew))
	for _, crew := range mov.Credits.Crew {
		if crew.Job == "Director" && crew.Department == "Directing" {
			directors = append(directors, crew.Name)
		}
	}
	return directors
}

func callMovieFunctions(
	update *tgbotapi.Update,
	mediaType string,
	mediaTitle string,
) ([]types.Movie, error) {
	var arrayPos int
	var err error
	var res interface{}
	if update.CallbackQuery == nil {
		res, err = searchClient.SearchMedia(mediaType, mediaTitle)
		if err != nil {
			return nil, err
		}
		movies = res.([]types.Movie)
		if len(movies) == 0 {
			err = fmt.Errorf("No film results: %w", err)
			return nil, err
		}
	} else if arrayPos, err = strconv.Atoi(update.CallbackQuery.Data); err != nil {
		return nil, err
	}
	details, _, err := detailsClient.GetDetails(mediaType, strconv.Itoa(movies[arrayPos].ID))
	if err != nil {
		return nil, err
	}
	providers, err := providersClient.SearchProviders(mediaType, strconv.Itoa(movies[arrayPos].ID))
	if err != nil {
		return nil, err
	}
	credits, err := directorsClient.GetMovieCredits(strconv.Itoa(movies[arrayPos].ID))
	if err != nil {
		return nil, err
	}
	movies[arrayPos].Details = *details
	movies[arrayPos].Providers = *providers
	movies[arrayPos].Credits = *credits
	return movies, nil
}
