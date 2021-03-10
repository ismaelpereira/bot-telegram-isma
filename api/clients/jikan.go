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

	"github.com/ismaelpereira/telegram-bot-isma/config"
	r "github.com/ismaelpereira/telegram-bot-isma/redis"
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
	SearchAnimeOrManga(string, string) ([]types.Manga, []types.Anime, error)
}

func NewJikanAPI(mediaTitle string, mediaType string) (JikanAPI, error) {
	return &jikanAPICached{
		api: &jikanAPI{
			mediaTitle: mediaTitle,
		},
		cache: nil,
		redis: nil,
	}, nil
}

type MangaDetails interface {
	GetMangaPageDetails(string, string) ([]byte, []byte, error)
}

func NewMangaAPI(mangaID string, mangaName string) (MangaDetails, error) {
	return &mangaAPICached{
		api: &mediaDetails{
			ID:    mangaID,
			title: mangaName,
		},
		jpCache:     nil,
		statusCache: nil,
		redis:       nil,
	}, nil
}

func (t *jikanAPI) SearchAnimeOrManga(mediaTitle string, mediaType string) ([]types.Manga, []types.Anime, error) {
	log.Println("jikan api")
	var apiResult *http.Response
	var err error
	if mediaType == "animes" {
		apiResult, err = http.Get("https://api.jikan.moe/v3/search/anime?q=" +
			url.QueryEscape(mediaTitle) + "&page=1")
		if err != nil {
			return nil, nil, err
		}
		defer apiResult.Body.Close()
	}
	if mediaType == "mangas" {
		apiResult, err = http.Get("https://api.jikan.moe/v3/search/manga?q=" +
			url.QueryEscape(mediaTitle) + "&page=1")
		if err != nil {
			return nil, nil, err
		}
		defer apiResult.Body.Close()
	}
	searchResult, err := ioutil.ReadAll(apiResult.Body)
	if err != nil {
		return nil, nil, err
	}
	var jikanResults types.JikanResponse
	err = json.Unmarshal(searchResult, &jikanResults)
	if err != nil {
		return nil, nil, err
	}
	if mediaType == "animes" {
		var animes []types.Anime
		err = json.Unmarshal(jikanResults.Data, &animes)
		if err != nil {
			return nil, nil, err
		}
		return nil, animes, nil
	}
	if mediaType == "mangas" {
		var mangas []types.Manga
		err = json.Unmarshal(jikanResults.Data, &mangas)
		if err != nil {
			return nil, nil, err
		}
		return mangas, nil, nil
	}
	return nil, nil, nil
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
	cfg, err := config.Wire()
	if err != nil {
		return nil, nil, err
	}
	t.redis, err = r.Wire(cfg)
	if err != nil {
		return nil, nil, err
	}
	keys, err := t.redis.Keys("telegram:manga:details:*").Result()
	if err != nil {
		return nil, nil, err
	}
	for _, key := range keys {
		details := strings.TrimPrefix(key, "telegram:manga:details:")
		if strings.TrimPrefix(details, "japaneseName:") == mangaName {
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return nil, nil, err
			}
			t.jpCache = data
		}
		if strings.TrimPrefix(details, "status:") == mangaName {
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return nil, nil, err
			}
			t.statusCache = data
		}
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
	if err = t.redis.Set(jpKey, jpRes, 30*time.Second).Err(); err != nil {
		return nil, nil, err
	}
	if err = t.redis.Set(statusKey, statusRes, 30*time.Second).Err(); err != nil {
		return nil, nil, err
	}

	return jpRes, statusRes, nil
}

type jikanAPICached struct {
	api   JikanAPI
	cache interface{}
	redis *redis.Client
}

func (t *jikanAPICached) SearchAnimeOrManga(mediaTitle string, mediaType string) ([]types.Manga, []types.Anime, error) {
	log.Println("jikan api cached")
	cfg, err := config.Wire()
	if err != nil {
		return nil, nil, err
	}
	t.redis, err = r.Wire(cfg)
	if err != nil {
		return nil, nil, err
	}
	keys, err := t.redis.Keys("telegram:" + mediaType + ":*").Result()
	if err != nil {
		return nil, nil, err
	}
	if len(keys) != 0 {
		for _, key := range keys {
			if strings.TrimPrefix(key, "telegram:"+mediaType+":") == mediaTitle {
				var animeMedia []types.Anime
				var mangaMedia []types.Manga
				data, err := t.redis.Get(key).Bytes()
				if err != nil {
					return nil, nil, err
				}
				if mediaType == "animes" {
					err = json.Unmarshal(data, &animeMedia)
					if err != nil {
						return nil, nil, err
					}
					t.cache = animeMedia
				}
				if mediaType == "mangas" {
					err = json.Unmarshal(data, &mangaMedia)
					if err != nil {
						return nil, nil, err
					}
					t.cache = mangaMedia
				}
			}
		}
	}
	var resJSON []byte
	if mediaType == "animes" {
		if t.cache != nil {
			return nil, t.cache.([]types.Anime), nil
		}
		_, resAnime, err := t.api.SearchAnimeOrManga(mediaTitle, mediaType)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resAnime)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":" + mediaTitle
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resAnime
		return nil, resAnime, nil
	}
	if mediaType == "mangas" {
		if t.cache != nil {
			return t.cache.([]types.Manga), nil, nil
		}
		resManga, _, err := t.api.SearchAnimeOrManga(mediaTitle, mediaType)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resManga)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":" + mediaTitle
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resManga
		return resManga, nil, nil
	}
	return nil, nil, nil
}
