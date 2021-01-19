package msgs

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/IsmaelPereira/telegram-bot-isma/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	IconThumbsUp      = "👍"
	IconX             = "❌"
	IconDevil         = "😈"
	IconPointingRight = "👉"
	IconPointingDown  = "👇"
	IconSkull         = "💀"
	IconWarning       = "⚠"
	IconAlarmClock    = "⏰"

	MsgThumbsUp       = IconThumbsUp
	MsgCantUnderstand = IconX + " -- Desculpe, não entendi"
	MsgNotAuthorized  = IconDevil + " -- Desculpe, você não tem permissão para isso"
	MsgServerError    = IconSkull + " -- Desculpe, tem algo de errado comigo..."
	MsgNotFound       = IconWarning + " -- Desculpe, não consegui encontrar isso"
	MsgHelp           = IconThumbsUp + " -- Os comandos são:\n/admiral\n/anime\n/manga\n"
	MsgAdmiral        = IconWarning + " -- The Admiral command is /admiral <admiral name> "
	MsgAnime          = IconWarning + " -- O comando é /anime <nome do anime>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgManga          = IconWarning + " -- O comando é /manga <nome do mangá>\nO resultado é baseado em uma pesquisa no MyanimeList"
)

//GetAdmiralPictureAndSendMessage is a function for admiral controller
func GetAdmiralPictureAndSendMessage(ad types.Admiral, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	adPicture, err := http.Get(ad.ProfilePicture)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, MsgServerError)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
	defer adPicture.Body.Close()
	adPictureData, err := ioutil.ReadAll(adPicture.Body)
	if err != nil {
		log.Println(err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, MsgNotFound)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}

	adMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: adPictureData})
	adMessage.Caption = "Nome real: " + ad.RealName + "\nNome de almirante: " + ad.AdmiralName + "\nIdade: " + strconv.Itoa(ad.Age) + "\nData de nascimento: " + ad.BirthDate + "\nSigno: " + ad.Sign + "\nAltura: " + strconv.FormatFloat(ad.Height, 'f', 2, 64) + "\nAkuma no Mi: " + ad.AkumaNoMi + "\nAnimal: " + ad.Animal + "\nPoder: " + ad.Power + "\nInspirado em: " + ad.ActorWhoInspire
	_, err = bot.Send(adMessage)
	if err != nil {
		log.Println(err)
	}
}

//GetAnimePictureAndSendMessage is a function for anime controller
func GetAnimePictureAndSendMessage(an types.Anime, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	anPicture, err := http.Get(an.CoverPicture)
	if err != nil {
		log.Println(err)
		tgbotapi.NewMessage(update.Message.Chat.ID, MsgServerError)
	}
	defer anPicture.Body.Close()
	anPictureData, err := ioutil.ReadAll(anPicture.Body)
	if err != nil {
		log.Println(err)
		tgbotapi.NewMessage(update.Message.Chat.ID, MsgNotFound)
	}
	anMessage := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: anPictureData})
	var airing string
	if an.Airing == true {
		airing = "Sim"
	} else {
		airing = "Não"
	}
	animeEpisodes := strconv.Itoa(an.Episodes)
	if animeEpisodes == "0" {
		animeEpisodes = "?"
	}
	anMessage.Caption = "Título: " + an.Title + "\nNota: " + strconv.FormatFloat(an.Score, 'f', 2, 64) + "\nEpisódios: " + animeEpisodes + "\nPassando? " + airing
	_, err = bot.Send(anMessage)
	if err != nil {
		log.Println(err)
	}
}

//GetMangaPictureAndSendMessage is a function for manga controller
func GetMangaPictureAndSendMessage(m types.Manga, update *tgbotapi.Update, bot *tgbotapi.BotAPI) {
	mPicture, err := http.Get(m.CoverPicture)
	if err != nil {
		log.Println(err)
		tgbotapi.NewMessage(update.Message.Chat.ID, MsgServerError)
	}
	defer mPicture.Body.Close()
	mPictureData, err := ioutil.ReadAll(mPicture.Body)
	if err != nil {
		log.Println(err)
		tgbotapi.NewMessage(update.Message.Chat.ID, MsgNotFound)
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
	GetMangaStatus(&m)
	if err != nil {
		log.Println(err)
	}
	mMessage.Caption = "Título: " + m.Title + "\nNome Japonês: " + string(m.JapaneseName) + "\nNota: " + strconv.FormatFloat(m.Score, 'f', 2, 64) + "\nVolumes: " + volumesNumber + "\nCapítulos: " + chaptersNumber + "\nStatus: " + m.Status
	_, err = bot.Send(mMessage)
	if err != nil {
		log.Println(err)
	}
}

//GetMangaStatus is a function for get the required manga in MAL site
func GetMangaStatus(m *types.Manga) error {
	idManga := strconv.Itoa(m.ID)
	animeListURL, err := http.Get("https://myanimelist.net/manga/" + url.QueryEscape(idManga))
	if err != nil {
		log.Println(err)
		return err
	}
	defer animeListURL.Body.Close()
	animeListCode, err := ioutil.ReadAll(animeListURL.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	japaneseStartPosition := []byte("Japanese:</span>")
	japaneseEndPosition := []byte("</div>")
	// statusStartPosition := string("Status:</span>")
	// statusEndPosition := string("</div>")
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

// chromeDriver := webdriver.NewChromeDriver("./chromedriver")
// err := chromeDriver.Start()
// if err != nil {
// 	log.Println(err)
// }
// defer chromeDriver.Stop()
// desired := webdriver.Capabilities{"Plataform": "Linux"}
// required := webdriver.Capabilities{}
// session, err := chromeDriver.NewSession(desired, required)
// if err != nil {
// 	log.Println(err)
// }
// defer session.Delete()
// idManga := strconv.Itoa(m.ID)
// err = session.Url("https://myanimelist.net/manga/" + url.QueryEscape(idManga))
// if err != nil {
// 	log.Println(err)
// }
// mangaDetailsBytes, err := session.ExecuteScript(`return Array.from(document.querySelectorAll(".dark_text")).map(el=>el.parentNode.innerText)`, []interface{}{})
// if err != nil {
// 	log.Println(err)
// }
// var mangaDetails []string
// err = json.Unmarshal(mangaDetailsBytes, &mangaDetails)
// if err != nil {
// 	log.Println(err)
// }
// for _, ssData := range mangaDetails {
// 	if strings.HasPrefix(ssData, "Status: ") == true {
// 		m.Status = ssData
// 	}
// 	if strings.HasPrefix(ssData, "Japanese: ") == true {
// 		m.JapaneseName = ssData
// 	}
// }
