package common

import (
	"encoding/json"
	"strings"

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
		switch mediaType {
		case "animes":
			{
				var animeMedia []types.Anime
				err = json.Unmarshal(data, &animeMedia)
				if err != nil {
					return nil, err
				}
				cache = animeMedia
			}
		case "mangas":
			{
				var mangaMedia []types.Manga
				err = json.Unmarshal(data, &mangaMedia)
				if err != nil {
					return nil, err
				}
				cache = mangaMedia
			}
		case "movies":
			{
				var movieMedia []types.Movie
				err = json.Unmarshal(data, &movieMedia)
				if err != nil {
					return nil, err
				}
				cache = movieMedia
			}
		case "tvshows":
			{
				var tvShowMedia []types.TVShow
				err = json.Unmarshal(data, &tvShowMedia)
				if err != nil {
					return nil, err
				}
				cache = tvShowMedia
			}
		}
	}
	return cache, nil
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
