package controllers

import (
	"strconv"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/api/clients"
	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/config"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MangaHandleUpdate is a function for manga work
func MangasHandleUpdate(c *config.Config, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	mangaName := strings.TrimSpace(update.Message.CommandArguments())
	if mangaName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgMangas)
		_, err := bot.Send(msg)
		return err
	}
	var mangaSearch *clients.MangaApi
	searchResults, err := mangaSearch.SearchManga(mangaName)
	if err != nil {
		return err
	}
	for _, manga := range searchResults.Results {
		getMangasPictureAndSendMessage(manga, update, bot)
	}
	return nil
}

func getMangasPictureAndSendMessage(m types.Manga, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	mMessage := tgbotapi.NewPhotoShare(update.Message.Chat.ID, m.CoverPicture)
	volumesNumber := strconv.Itoa(m.Volumes)
	chaptersNumber := strconv.Itoa(m.Chapters)
	if volumesNumber == "0" {
		volumesNumber = "?"
	}
	if chaptersNumber == "0" {
		chaptersNumber = "?"
	}
	var mangaSearch *clients.MangaApi
	japaneseName, status, err := mangaSearch.GetMangaPageDetails(strconv.Itoa(m.ID))
	if err != nil {
		return err
	}
	mMessage.Caption = "Título: " + m.Title +
		"\nNome Japonês: " + string(japaneseName) +
		"\nNota: " + strconv.FormatFloat(m.Score, 'f', 2, 64) +
		"\nVolumes: " + volumesNumber +
		"\nCapítulos: " + chaptersNumber +
		"\nStatus: " + status
	_, err = bot.Send(mMessage)
	return err
}
