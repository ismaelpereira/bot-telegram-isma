package controllers

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

func SendAnimesKeyboard(animes []types.Anime) []tgbotapi.InlineKeyboardMarkup {
	var kb []tgbotapi.InlineKeyboardMarkup
	if len(animes) > 1 {
		kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "animes:1"),
			),
		))
	}
	return kb
}

func SendAnimesCallbackKeyboard(animes []types.Anime, i int) []tgbotapi.InlineKeyboardButton {
	var kb []tgbotapi.InlineKeyboardButton
	if i != 0 {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "animes:"+strconv.Itoa(i-1)),
		)
	}
	if i != (len(animes) - 1) {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "animes:"+strconv.Itoa(i+1)),
		)
	}
	return kb
}

func SendMangasKeyboard(mangas []types.Manga) []tgbotapi.InlineKeyboardMarkup {
	var kb []tgbotapi.InlineKeyboardMarkup
	if len(mangas) > 1 {
		kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "mangas:1"),
			),
		))
	}
	return kb
}

func SendMangasCallbackKeyboard(mangas []types.Manga, i int) []tgbotapi.InlineKeyboardButton {
	var kb []tgbotapi.InlineKeyboardButton
	if i != 0 {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconPrevious, "mangas:"+strconv.Itoa(i-1)),
		)
	}
	if i != (len(mangas) - 1) {
		kb = append(kb,
			tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "mangas:"+strconv.Itoa(i+1)),
		)
	}
	return kb
}

func SendMoviesKeyboard(movies []types.Movie) []tgbotapi.InlineKeyboardMarkup {
	var kb []tgbotapi.InlineKeyboardMarkup
	if len(movies) > 1 {
		kb = append(kb, tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(msgs.IconNext, "movies:1"),
			),
		))
	}
	return kb
}

func SendMoviesCallbackKeyboard(movies []types.Movie, i int) []tgbotapi.InlineKeyboardButton {
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
	return kb
}

func SendTVShowsKeyboard(tvShows []types.TVShow) []tgbotapi.InlineKeyboardMarkup {
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
	return kb
}

func SendTVShowsCallbackKeyboard(tvShows []types.TVShow, i int) []tgbotapi.InlineKeyboardButton {
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
	return kb
}
