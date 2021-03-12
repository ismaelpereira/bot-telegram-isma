package msgs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

const (
	IconThumbsUp      = "👍"
	IconX             = "❌"
	IconOk            = "✅"
	IconDevil         = "😈"
	IconPointingRight = "👉"
	IconPointingDown  = "👇"
	IconSkull         = "💀"
	IconWarning       = "⚠"
	IconAlarmClock    = "⏰"
	IconPrevious      = "❰"
	IconNext          = "❯"

	MsgThumbsUp       = IconThumbsUp
	MsgCantUnderstand = IconX + " -- Desculpe, não entendi"
	MsgNotAuthorized  = IconDevil + " -- Desculpe, você não tem permissão para isso"
	MsgServerError    = IconSkull + " -- Desculpe, tem algo de errado comigo..."
	MsgNotFound       = IconWarning + " -- Desculpe, não consegui encontrar isso"
	MsgHelp           = IconThumbsUp + " -- Os comandos são:\n/admirals\n/animes\n/mangas\n/money\n/movies\n" +
		"/tvshows\n/now\n/reminder\n/checklist"
	MsgAdmirals = IconWarning + " -- The Admiral command is /admirals <admiral name> "
	MsgAnimes   = IconWarning + " -- O comando é /animes <nome do anime>\n" +
		"O resultado é baseado em uma pesquisa no MyanimeList"
	MsgMangas = IconWarning + " -- O comando é /mangas <nome do mangá>\n" +
		"O resultado é baseado em uma pesquisa no MyanimeList"
	MsgMoney    = IconWarning + "-- O comando é /money <quantidade> <moeda principal> <moeda a ser convertida>"
	MsgMovies   = IconWarning + "-- O comando é /movies <nome do filme> O resultado é baseado em uma pesquisa do MovieDB"
	MsgTVShow   = IconWarning + "-- O comando é /tvshows <nome da serie> O resultado é baseado em uma pesquisa do MovieDB"
	MsgReminder = IconWarning + "-- O comando é /reminder <tempo> <medida de tempo> <mensagem>"
	MsgNow      = IconWarning + "-- O comando é /now <operação> <tempo> <medida de tempo>"
)

func EditMessage(
	cfg *config.Config,
	chatID int64,
	messageID int,
	posterPath string,
	caption string,
	replyMarkup tgbotapi.InlineKeyboardMarkup,
) error {
	var msgEdit types.EditMediaJSON
	msgEdit.ChatID = chatID
	msgEdit.MessageID = messageID
	msgEdit.Media.Type = "photo"
	if posterPath == "" || posterPath == "https://www.themoviedb.org/t/p/w300_and_h450_bestv2" {
		msgEdit.Media.URL = "https://badybassitt.sp.gov.br/lib/img/no-image.jpg"
	} else {
		msgEdit.Media.URL = posterPath
	}
	msgEdit.Media.Caption = caption
	msgEdit.ReplyMarkup = replyMarkup
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
	if sendMessage.StatusCode < 200 || sendMessage.StatusCode > 299 {
		err = fmt.Errorf("Error in post method %w", err)
		return err
	}
	return nil
}
