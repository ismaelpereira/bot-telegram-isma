package controllers

import (
	"log"

	"github.com/IsmaelPereira/telegram-bot-isma/bot/msgs"
	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func StorageHandleUpdate(bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
	userID, err := bot.GetMe()
	spew.Dump(userID)
	if userID.ID != 484887886 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgs.MsgNotAuthorized)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	fileID, err := fileHandler(bot, update)
	if err != nil {
		log.Println(err)
		return err
	}
	spew.Dump(fileID)
	url, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		log.Println(err)
		return err
	}
	spew.Dump(url)
	return nil
}

func fileHandler(bot *tgbotapi.BotAPI, update *tgbotapi.Update) (string, error) {
	var fileID string
	var err error
	switch {
	case update.Message.Animation != nil:
		{
			fileID = update.Message.Animation.FileID
		}
	case update.Message.Audio != nil:
		{
			fileID = update.Message.Audio.FileID
		}
	case update.Message.Document != nil:
		{
			fileID = update.Message.Document.FileID
		}
	case update.Message.Sticker != nil:
		{
			fileID = update.Message.Sticker.FileID
		}
	case update.Message.Video != nil:
		{
			fileID = update.Message.Document.FileID
		}
	case update.Message.VideoNote != nil:
		{
			fileID = update.Message.VideoNote.FileID
		}
	case update.Message.Voice != nil:
		{
			fileID = update.Message.Voice.FileID
		}
	}
	if fileID == "" {
		return "error", err
	}

	return fileID, err

}
