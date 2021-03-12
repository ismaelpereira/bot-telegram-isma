package common

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	r "github.com/ismaelpereira/telegram-bot-isma/redis"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

func SearchItens(redis *redis.Client, mediaType string, mediaTitle string) (interface{}, error) {
	keys, err := redis.Keys("telegram:" + mediaType + ":*").Result()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, nil
	}
	var cache interface{}
	for _, key := range keys {
		if strings.TrimPrefix(key, "telegram:"+mediaType+":") != mediaTitle {
			continue
		}
		data, err := redis.Get(key).Bytes()
		if err != nil {
			return nil, err
		}
		handler := map[string]types.RedisHandler{
			"animes":  cmdAnimes,
			"mangas":  cmdMangas,
			"movies":  cmdMovies,
			"tvshows": cmdTVShows,
		}
		if f, ok := handler[mediaType]; ok {
			cache, err = f(data)
			if err != nil {
				return nil, err
			}
		}
	}
	return cache, nil
}

func cmdAnimes(data []byte) (interface{}, error) {
	var animeMedia []types.Anime
	err := json.Unmarshal(data, &animeMedia)
	if err != nil {
		return nil, err
	}
	return animeMedia, nil
}

func cmdMangas(data []byte) (interface{}, error) {
	var mangaMedia []types.Manga
	err := json.Unmarshal(data, &mangaMedia)
	if err != nil {
		return nil, err
	}
	return mangaMedia, nil
}

func cmdMovies(data []byte) (interface{}, error) {
	var movieMedia []types.Movie
	err := json.Unmarshal(data, &movieMedia)
	if err != nil {
		return nil, err
	}
	return movieMedia, nil
}

func cmdTVShows(data []byte) (interface{}, error) {
	var tvShowMedia []types.TVShow
	err := json.Unmarshal(data, &tvShowMedia)
	if err != nil {
		return nil, err
	}
	return tvShowMedia, nil
}

func SetRedis() (*redis.Client, error) {
	cfg, err := config.Wire()
	if err != nil {
		return nil, err
	}
	redis, err := r.Wire(cfg)
	if err != nil {
		return nil, err
	}
	return redis, nil
}

func SetRedisKey(resJSON []byte, redis *redis.Client, mediaType string, mediaTitle string) error {
	key := "telegram:" + mediaType + ":" + mediaTitle
	err := redis.Set(key, resJSON, 72*time.Hour).Err()
	return err
}

func SetRedisKeyDetails(resJSON []byte, redis *redis.Client, mediaType string, mediaID string) error {
	key := "telegram:" + mediaType + ":details:" + mediaID
	err := redis.Set(key, resJSON, 72*time.Hour).Err()
	return err
}
