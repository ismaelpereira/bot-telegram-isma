package clients

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type AnimeApi struct {
	title string
}

type MangaApi struct {
	title string
	ID    int
}

func (t *AnimeApi) SearchAnime(animeName string) (*types.AnimeResponse, error) {
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/anime?q=" + url.QueryEscape(animeName) + "&page=1&limit=3")
	if err != nil {
		return nil, err
	}
	defer apiResult.Body.Close()
	readAnimes, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return nil, err
	}
	var animeSearchResults types.AnimeResponse
	err = json.Unmarshal(readAnimes, &animeSearchResults)
	if err != nil {
		return nil, err
	}
	return &animeSearchResults, nil
}

func (t *MangaApi) SearchManga(mangaName string) (*types.MangaResponse, error) {
	apiResult, err := http.Get("https://api.jikan.moe/v3/search/manga?q=" + url.QueryEscape(mangaName) + "&page=1&limit=3")
	if err != nil {
		return nil, err
	}
	defer apiResult.Body.Close()
	searchResult, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return nil, err
	}
	var mangasSearchResults types.MangaResponse
	err = json.Unmarshal(searchResult, &mangasSearchResults)
	if err != nil {
		return nil, err
	}
	return &mangasSearchResults, nil
}

func (t *MangaApi) GetMangaPageDetails(mangaID string) ([]byte, string, error) {
	animeListURL, err := http.Get("https://myanimelist.net/manga/" + url.QueryEscape(mangaID))
	if err != nil {
		return nil, "", err
	}
	defer animeListURL.Body.Close()
	myAnimeListPageCode, err := ioutil.ReadAll(animeListURL.Body)
	if err != nil {
		return nil, "", err
	}
	japaneseStartPosition := []byte("Japanese:</span>")
	japaneseEndPosition := []byte("</div>")
	startJp := bytes.Index(myAnimeListPageCode, japaneseStartPosition)
	endJp := bytes.Index(myAnimeListPageCode[startJp:], japaneseEndPosition)
	japaneseName := bytes.TrimSpace(myAnimeListPageCode[startJp+len(japaneseStartPosition) : startJp+endJp])
	statusStartPosition := string("Status:</span>")
	statusEndPosition := string("</div>")
	startSt := strings.Index(string(myAnimeListPageCode), statusStartPosition)
	endSt := strings.Index(string(myAnimeListPageCode)[startSt:], statusEndPosition)
	status := strings.TrimSpace(string(myAnimeListPageCode)[startSt+len(statusStartPosition) : startSt+endSt])
	return japaneseName, status, nil
}
