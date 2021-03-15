package controllers

import (
	"strconv"
	"strings"

	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

var jikanAPI clients.JikanAPI
var animes []types.Anime

func init() {
	var err error
	jikanAPI, err = clients.NewJikanAPI()
	if err != nil {
		panic(err)
	}
}

// AnimeHandleUpdate is a function for anime work
func AnimesHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		return animesArrowButtonAction(cfg, update, animes)
	}
	mediaType := "animes"
	animeName := strings.TrimSpace(update.Message.CommandArguments())
	if animeName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAnimes)
		_, err := bot.Send(msg)
		return err
	}
	res, err := jikanAPI.SearchAnimeOrManga(mediaType, animeName)
	if err != nil {
		return err
	}
	animes = res.([]types.Anime)
	if len(animes) == 0 {
		return nil
	}
	animeMessage := getAnimesPictureAndSendMessage(update, animes[0])
	kb := SendAnimesKeyboard(animes)
	if len(animes) > 1 {
		animeMessage.ReplyMarkup = kb[0]
	}
	_, err = bot.Send(animeMessage)
	return err
}

func getAnimesPictureAndSendMessage(
	update *tgbotapi.Update,
	an types.Anime,
) *tgbotapi.PhotoConfig {
	var anMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		anMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID, an.CoverPicture)
	}
	if update.CallbackQuery != nil {
		anMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID, an.CoverPicture)
	}
	var airing string
	if an.Airing {
		airing = "Passando"
	} else {
		airing = "Finalizado"
	}
	animeEpisodes := strconv.Itoa(an.Episodes)
	if animeEpisodes == "0" {
		animeEpisodes = "?"
	}
	anMessage.Caption = "Título: " + an.Title +
		"\nNota: " + strconv.FormatFloat(an.Score, 'f', 2, 64) +
		"\nEpisódios: " + animeEpisodes +
		"\nStatus: " + airing
	return &anMessage
}

func animesArrowButtonAction(
	cfg *config.Config,
	update *tgbotapi.Update,
	animes []types.Anime,
) error {
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	animeMessage := getAnimesPictureAndSendMessage(update, animes[i])
	kb := SendAnimesCallbackKeyboard(animes, i)
	err = msgs.EditMessage(
		cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		animes[i].CoverPicture,
		animeMessage.Caption,
		tgbotapi.NewInlineKeyboardMarkup(
			kb,
		),
	)
	return err
}
