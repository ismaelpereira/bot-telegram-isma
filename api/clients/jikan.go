package clients

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/ismaelpereira/telegram-bot-isma/api/common"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type mediaDetails struct {
	ID    string
	title string
}

type jikanAPI struct {
	mediaTitle string
}

type JikanAPI interface {
	SearchAnimeOrManga(string, string) (interface{}, error)
}

func NewJikanAPI() (JikanAPI, error) {
	return &jikanAPICached{
		api:   &jikanAPI{},
		cache: nil,
		redis: nil,
	}, nil
}

type MangaDetails interface {
	GetMangaPageDetails(string, string) ([]byte, []byte, error)
}

func NewMangaAPI() (MangaDetails, error) {
	return &mangaAPICached{
		api:         &mediaDetails{},
		jpCache:     nil,
		statusCache: nil,
		redis:       nil,
	}, nil
}

func (t *jikanAPI) SearchAnimeOrManga(mediaType string, mediaTitle string) (interface{}, error) {
	log.Println("jikan api")
	if mediaType == "animes" {
		var animes []types.Anime
		url := "https://api.jikan.moe/v3/search/anime?q=" +
			url.QueryEscape(mediaTitle) + "&page=1"
		if err := t.httpGET(url, &animes); err != nil {
			return nil, err
		}

		return animes, nil
	}
	if mediaType == "mangas" {
		var mangas []types.Manga
		url := "https://api.jikan.moe/v3/search/manga?q=" +
			url.QueryEscape(mediaTitle) + "&page=1"
		if err := t.httpGET(url, &mangas); err != nil {
			return nil, err
		}
		return mangas, nil
	}
	return nil, nil
}

func (t *jikanAPI) httpGET(url string, v interface{}) error {
	var apiResult *http.Response
	var err error
	if apiResult, err = http.Get(url); err != nil {
		return err
	}
	defer apiResult.Body.Close()
	searchResult, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return err
	}
	var jikanResults types.JikanResponse
	err = json.Unmarshal(searchResult, &jikanResults)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(jikanResults.Data, &v); err != nil {
		return err
	}
	return nil
}

func (t *mediaDetails) GetMangaPageDetails(mangaID string, mangaName string) ([]byte, []byte, error) {
	log.Println("gettin manga details")
	animeListURL, err := http.Get("https://myanimelist.net/manga/" + url.QueryEscape(mangaID))
	if err != nil {
		return nil, nil, err
	}
	defer animeListURL.Body.Close()
	myAnimeListPageCode, err := ioutil.ReadAll(animeListURL.Body)
	if err != nil {
		return nil, nil, err
	}
	japaneseStartPosition := []byte("Japanese:</span>")
	japaneseEndPosition := []byte("</div>")
	startJp := bytes.Index(myAnimeListPageCode, japaneseStartPosition)
	endJp := bytes.Index(myAnimeListPageCode[startJp:], japaneseEndPosition)
	japaneseName := bytes.TrimSpace(myAnimeListPageCode[startJp+len(japaneseStartPosition) : startJp+endJp])
	statusStartPosition := []byte("Status:</span>")
	statusEndPosition := []byte("</div>")
	startSt := bytes.Index(myAnimeListPageCode, statusStartPosition)
	endSt := bytes.Index(myAnimeListPageCode[startSt:], statusEndPosition)
	status := bytes.TrimSpace(myAnimeListPageCode[startSt+len(statusStartPosition) : startSt+endSt])
	return japaneseName, status, nil
}

type mangaAPICached struct {
	api         MangaDetails
	jpCache     interface{}
	statusCache interface{}
	redis       *redis.Client
}

func (t *mangaAPICached) GetMangaPageDetails(mangaID string, mangaName string) ([]byte, []byte, error) {
	log.Println("manga details cached")
	err := t.searchRedisDetailsKeys(mangaName)
	if err != nil {
		return nil, nil, err
	}
	if t.jpCache != nil && t.statusCache != nil {
		return t.jpCache.([]byte), t.statusCache.([]byte), nil
	}
	jpRes, statusRes, err := t.api.GetMangaPageDetails(mangaID, mangaName)
	if err != nil {
		return nil, nil, err
	}
	t.jpCache = jpRes
	t.statusCache = statusRes
	jpKey := "telegram:manga:details:japaneseName:" + mangaName
	statusKey := "telegram:manga:details:status:" + mangaName
	if err = t.redis.Set(jpKey, jpRes, 72*time.Hour).Err(); err != nil {
		return nil, nil, err
	}
	if err = t.redis.Set(statusKey, statusRes, 72*time.Hour).Err(); err != nil {
		return nil, nil, err
	}
	return jpRes, statusRes, nil
}

func (t *mangaAPICached) searchRedisDetailsKeys(mangaName string) error {
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return err
	}
	keys, err := t.redis.Keys("telegram:manga:details:*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		details := strings.TrimPrefix(key, "telegram:manga:details:")
		if strings.TrimPrefix(details, "japaneseName:") == mangaName {
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return err
			}
			t.jpCache = data
			return nil
		}
		if strings.TrimPrefix(details, "status:") == mangaName {
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return err
			}
			t.statusCache = data
			return nil
		}
	}
	t.jpCache = nil
	t.statusCache = nil
	return nil
}

type jikanAPICached struct {
	api   JikanAPI
	cache interface{}
	redis *redis.Client
}

func (t *jikanAPICached) SearchAnimeOrManga(mediaType string, mediaTitle string) (interface{}, error) {
	log.Println("jikan api cached")
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return nil, err
	}
	cache, err := common.SearchItens(t.redis, mediaType, mediaTitle)
	if err != nil {
		return nil, err
	}
	t.cache = cache
	if mediaType == "animes" {
		if t.cache != nil {
			return t.cache.([]types.Anime), nil
		}
	}
	if mediaType == "mangas" {
		if t.cache != nil {
			return t.cache.([]types.Manga), nil
		}
	}
	var res interface{}
	if res, err = t.setInRedis(mediaType, mediaTitle); err != nil {
		return nil, err
	}
	t.cache = res
	return res, nil
}

func (t *jikanAPICached) setInRedis(mediaType string, mediaTitle string) (interface{}, error) {
	res, err := t.api.SearchAnimeOrManga(mediaType, mediaTitle)
	if err != nil {
		return nil, err
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	key := "telegram:" + mediaType + ":" + mediaTitle
	if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
		return nil, err
	}
	return res, nil
}
