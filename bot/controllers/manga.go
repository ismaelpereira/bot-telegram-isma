package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//MangaHandleUpdate is a function for manga work
func MangaHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	mangaName := strings.TrimSpace(update.Message.CommandArguments())
	if mangaName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgManga)
		_, err := bot.Send(msg)
		return err
	}
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/manga?q=" + url.QueryEscape(mangaName) + "&page=1&limit=3")
	if err != nil {
		return err
	}
	defer apiResult.Body.Close()
	readMangas, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return err
	}
	var searchResults types.MangaResponse
	err = json.Unmarshal(readMangas, &searchResults)
	if err != nil {
		return err
	}
	for _, manga := range searchResults.Results {
		getMangaPictureAndSendMessage(manga, update, bot)
	}
	return nil
}

func getMangaPictureAndSendMessage(m types.Manga, update *tgbotapi.Update, bot *tgbotapi.BotAPI) error {
	mPicture, err := http.Get(m.CoverPicture)
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgServerError)
	}
	defer mPicture.Body.Close()
	mPictureData, err := ioutil.ReadAll(mPicture.Body)
	if err != nil {
		tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgNotFound)
	}
	mMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: mPictureData})
	volumesNumber := strconv.Itoa(m.Volumes)
	chaptersNumber := strconv.Itoa(m.Chapters)
	if volumesNumber == "0" {
		volumesNumber = "?"
	}
	if chaptersNumber == "0" {
		chaptersNumber = "?"
	}
	getMangaStatus(&m)
	if err != nil {
		return err
	}
	mMessage.Caption = "Título: " + m.Title + "\nNome Japonês: " + string(m.JapaneseName) + "\nNota: " +
		strconv.FormatFloat(m.Score, 'f', 2, 64) + "\nVolumes: " + volumesNumber + "\nCapítulos: " + chaptersNumber +
		"\nStatus: " + m.Status
	_, err = bot.Send(mMessage)
	return err
}

func getMangaStatus(m *types.Manga) error {
	idManga := strconv.Itoa(m.ID)
	animeListURL, err := http.Get("https://myanimelist.net/manga/" + url.QueryEscape(idManga))
	if err != nil {
		return err
	}
	defer animeListURL.Body.Close()
	animeListCode, err := ioutil.ReadAll(animeListURL.Body)
	if err != nil {
		return err
	}
	japaneseStartPosition := []byte("Japanese:</span>")
	japaneseEndPosition := []byte("</div>")
	startJp := bytes.Index(animeListCode, japaneseStartPosition)
	endJp := bytes.Index(animeListCode[startJp:], japaneseEndPosition)
	m.JapaneseName = bytes.TrimSpace(animeListCode[startJp+len(japaneseStartPosition) : startJp+endJp])
	statusStartPosition := string("Status:</span>")
	statusEndPosition := string("</div>")
	startSt := strings.Index(string(animeListCode), statusStartPosition)
	endSt := strings.Index(string(animeListCode)[startSt:], statusEndPosition)
	m.Status = strings.TrimSpace(string(animeListCode)[startSt+len(statusStartPosition) : startSt+endSt])
	return nil
}
