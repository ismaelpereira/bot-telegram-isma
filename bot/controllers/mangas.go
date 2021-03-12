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

var mangas []types.Manga

// MangaHandleUpdate is a function for manga work
func MangasHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		return mangasArrowButtonAction(cfg, update, mangas)
	}
	command := update.Message.Command()
	mangaName := strings.TrimSpace(update.Message.CommandArguments())
	if mangaName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMangas)
		_, err := bot.Send(msg)
		return err
	}
	jikanAPI, err := clients.NewJikanAPI(mangaName, command)
	if err != nil {
		return err
	}
	mangas, _, err = jikanAPI.SearchAnimeOrManga(mangaName, command)
	if err != nil {
		return err
	}
	if len(mangas) == 0 {
		return nil
	}
	mangaMessage, err := getMangasPictureAndSendMessage(update, mangas[0])
	if err != nil {
		return err
	}
	kb := SendMangasKeyboard(mangas)
	if len(mangas) > 1 {
		mangaMessage.ReplyMarkup = kb[0]
	}
	_, err = bot.Send(mangaMessage)
	return err
}

func getMangasPictureAndSendMessage(
	update *tgbotapi.Update,
	m types.Manga,
) (*tgbotapi.PhotoConfig, error) {
	var err error
	mangaName := strings.TrimSpace(m.Title)
	var mMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		mMessage = tgbotapi.NewPhotoShare(update.Message.Chat.ID, m.CoverPicture)
	}
	if update.CallbackQuery != nil {
		mMessage = tgbotapi.NewPhotoShare(update.CallbackQuery.Message.Chat.ID, m.CoverPicture)
	}
	volumesNumber := strconv.Itoa(m.Volumes)
	chaptersNumber := strconv.Itoa(m.Chapters)
	if volumesNumber == "0" {
		volumesNumber = "?"
	}
	if chaptersNumber == "0" {
		chaptersNumber = "?"
	}
	var mangaSearch clients.MangaDetails
	mangaSearch, err = clients.NewMangaAPI(strconv.Itoa(m.ID), mangaName)
	if err != nil {
		return nil, err
	}
	japaneseName, status, err := mangaSearch.GetMangaPageDetails(strconv.Itoa(m.ID), m.Title)
	if err != nil {
		return nil, err
	}
	mMessage.Caption = "Título: " + m.Title +
		"\nNome Japonês: " + string(japaneseName) +
		"\nNota: " + strconv.FormatFloat(m.Score, 'f', 2, 64) +
		"\nVolumes: " + volumesNumber +
		"\nCapítulos: " + chaptersNumber +
		"\nStatus: " + string(status)
	return &mMessage, nil
}

func mangasArrowButtonAction(
	cfg *config.Config,
	update *tgbotapi.Update,
	mangas []types.Manga,
) error {
	i, err := strconv.Atoi(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	mangaMessage, err := getMangasPictureAndSendMessage(update, mangas[i])
	if err != nil {
		return err
	}
	kb := SendMangasCallbackKeyboard(mangas, i)
	err = msgs.EditMessage(
		cfg,
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		mangas[i].CoverPicture,
		mangaMessage.Caption,
		tgbotapi.NewInlineKeyboardMarkup(
			kb,
		),
	)
	return err
}
