package msgs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/IsmaelPereira/telegram-bot-isma/config"
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
	IconPrevious      = "⇦"
	IconNext          = "⇨"

	MsgThumbsUp       = IconThumbsUp
	MsgCantUnderstand = IconX + " -- Desculpe, não entendi"
	MsgNotAuthorized  = IconDevil + " -- Desculpe, você não tem permissão para isso"
	MsgServerError    = IconSkull + " -- Desculpe, tem algo de errado comigo..."
	MsgNotFound       = IconWarning + " -- Desculpe, não consegui encontrar isso"
	MsgHelp           = IconThumbsUp + " -- Os comandos são:\n/admiral\n/anime\n/manga\n/money\n/movie"
	MsgAdmiral        = IconWarning + " -- The Admiral command is /admiral <admiral name> "
	MsgAnime          = IconWarning + " -- O comando é /anime <nome do anime>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgManga          = IconWarning + " -- O comando é /manga <nome do mangá>\nO resultado é baseado em uma pesquisa no MyanimeList"
	MsgMoney          = IconWarning + "-- O comando é /money <quantidade> <moeda principal> <moeda a ser convertida>"
	MsgMovie          = IconWarning + "-- O comando é /movie <nome do filme> O resultado é baseado em uma pesquisa do MovieDB"
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
	adMessage.Caption = "Nome real: " + ad.RealName + "\nNome de almirante: " +
		ad.AdmiralName + "\nIdade: " + strconv.Itoa(ad.Age) + "\nData de nascimento: " +
		ad.BirthDate + "\nSigno: " + ad.Sign + "\nAltura: " + strconv.FormatFloat(ad.Height, 'f', 2, 64) +
		"\nAkuma no Mi: " + ad.AkumaNoMi + "\nAnimal: " + ad.Animal + "\nPoder: " + ad.Power + "\nInspirado em: " +
		ad.ActorWhoInspire
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
	anMessage.Caption = "Título: " + an.Title + "\nNota: " + strconv.FormatFloat(an.Score, 'f', 2, 64) +
		"\nEpisódios: " + animeEpisodes + "\nPassando? " + airing
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
	mMessage.Caption = "Título: " + m.Title + "\nNome Japonês: " + string(m.JapaneseName) + "\nNota: " +
		strconv.FormatFloat(m.Score, 'f', 2, 64) + "\nVolumes: " + volumesNumber + "\nCapítulos: " + chaptersNumber +
		"\nStatus: " + m.Status
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
	startJp := bytes.Index(animeListCode, japaneseStartPosition)
	endJp := bytes.Index(animeListCode[startJp:], japaneseEndPosition)
	m.JapaneseName = bytes.TrimSpace(animeListCode[startJp+len(japaneseStartPosition) : startJp+endJp])

	statusStartPosition := string("Status:</span>")
	statusEndPosition := string("</div>")
	startSt := strings.Index(string(animeListCode), statusStartPosition)
	endSt := strings.Index(string(animeListCode)[startSt:], statusEndPosition)
	m.Status = strings.TrimSpace(string(animeListCode)[startSt+len(statusStartPosition) : startSt+endSt])
	return err
}

func GetMoviePictureAndSendMessage(mov types.MovieDbSearchResults, update *tgbotapi.Update, bot *tgbotapi.BotAPI) (*tgbotapi.PhotoConfig, error) {
	var movDetailsMessage []string
	releaseDate, err := time.Parse("2006-01-02", mov.ReleaseDate)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	movDetailsMessage = append(movDetailsMessage,
		"\nTítulo: "+mov.Title,
		"\nTítulo Original: "+mov.OriginalTitle,
		"\nPopularidade: "+strconv.FormatFloat(mov.Popularity, 'f', 2, 64),
		"\nData de lançamento: "+releaseDate.Format("02/01/2006"),
	)
	movPicture, err := http.Get("https://themoviedb.org/t/p/w300_and_h450_bestv2" + mov.PosterPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer movPicture.Body.Close()
	movPictureData, err := ioutil.ReadAll(movPicture.Body)
	movieProvidersMessage, err := GetMovieProviders(mov)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var movMessage tgbotapi.PhotoConfig
	if update.CallbackQuery == nil {
		movMessage = tgbotapi.NewPhotoUpload(update.Message.Chat.ID, tgbotapi.FileBytes{Bytes: movPictureData})
	}
	if update.CallbackQuery != nil {
		movMessage = tgbotapi.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileBytes{Bytes: movPictureData})
	}
	movMessage.Caption = strings.Join(movDetailsMessage, "") + strings.Join(movieProvidersMessage, "")
	return &movMessage, nil
}
func GetMovieProviders(mov types.MovieDbSearchResults) (movProvidersMessage []string, err error) {
	apiKey, err := config.GetMovieApiKey()
	if err != nil {
		log.Println(err)

	}
	watchProviders, err := http.Get("https://api.themoviedb.org/3/movie/" +
		url.QueryEscape(strconv.Itoa(mov.ID)) + "/watch/providers?api_key=" +
		url.QueryEscape(apiKey))
	defer watchProviders.Body.Close()
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var providers types.WatchProvidersResponse
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if country, ok := providers.Results["BR"]; ok && country != nil {
		if country.Buy != nil {
			movProvidersMessage = append(movProvidersMessage, "\nPara Comprar: ")
			for i, providerBuy := range country.Buy {
				movProvidersMessage = append(movProvidersMessage, providerBuy.ProviderName)
				if i == len(country.Buy)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}

		if country.Rent != nil {
			movProvidersMessage = append(movProvidersMessage, "\nPara Alugar: ")
			for i, providerRent := range country.Rent {
				movProvidersMessage = append(movProvidersMessage, providerRent.ProviderName)
				if i == len(country.Rent)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}

		if country.Flatrate != nil {
			movProvidersMessage = append(movProvidersMessage, "\nServicos de streaming: ")
			for i, providerFlatrate := range country.Flatrate {
				movProvidersMessage = append(movProvidersMessage, providerFlatrate.ProviderName)
				if i == len(country.Flatrate)-1 {
					movProvidersMessage = append(movProvidersMessage, ".")
				} else {
					movProvidersMessage = append(movProvidersMessage, ", ")
				}
			}
		}
	}
	return movProvidersMessage, err
}
