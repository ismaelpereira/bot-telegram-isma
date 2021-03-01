package controllers

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

//AnimeHandleUpdate is a function for anime work
func AnimesHandleUpdate(cfg *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	animeName := strings.TrimSpace(update.Message.CommandArguments())
	if animeName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgAnimes)
		_, err := bot.Send(msg)
		return err
	}
	var animesSearch *clients.AnimeApi
	searchResults, err := animesSearch.SearchAnime(animeName)
	if err != nil {
		return err
	}
	for _, anime := range searchResults.Results {
		getAnimesPictureAndSendMessage(anime, update, bot)
	}
	return nil
}

func getAnimesPictureAndSendMessage(an types.Anime, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	anMessage := tgbotapi.NewPhotoShare(update.Message.Chat.ID, an.CoverPicture)
	var airing string
	if an.Airing == true {
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
	_, err := bot.Send(anMessage)
	return err
}
