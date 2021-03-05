package clients

import (
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

type movieDetails struct {
	ID string
}

type theMovieDBAPI struct {
	apiKey string
}

type SearchMedia interface {
	SearchMedia(string, string) ([]types.Movie, []types.TVShow, error)
}

func NewSearchMedia(mediaType string, mediaTitle string, apiKey string) (SearchMedia, error) {
	return &movieAPICached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
	}, nil
}

type SearchProviders interface {
	SearchProviders(string, string) (*types.WatchProviders, error)
}

func NewSearchProviders(mediaType string, mediaID string, apiKey string) (SearchProviders, error) {
	return &providersCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
	}, nil
}

type GetDetails interface {
	GetDetails(string, string) (*types.MovieDetails, *types.TVShowDetails, error)
}

func NewGetDetails(mediaType string, mediaID string, apiKey string) (GetDetails, error) {
	return &detailsCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
	}, nil
}

type GetMovieCredits interface {
	GetMovieCredits(string) (*types.MovieCredits, error)
}

func NewGetMovieCredits(movieID string, apiKey string) (GetMovieCredits, error) {
	return &moviesCreditCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
	}, nil
}

func (t *theMovieDBAPI) SearchMedia(
	mediaType string,
	mediaTitle string,
) ([]types.Movie, []types.TVShow, error) {
	log.Println("moviedb api")
	var movieAPI *http.Response
	var err error
	if mediaType == "movies" {
		movieAPI, err = http.Get("https://api.themoviedb.org/3/search/movie?api_key=" + url.QueryEscape(t.apiKey) +
			"&page=1&langague=pt-br&query=" + url.QueryEscape(mediaTitle))
		if err != nil {
			return nil, nil, err
		}
	}
	if mediaType == "tvshows" {
		movieAPI, err = http.Get("https://api.themoviedb.org/3/search/tv?api_key=" + url.QueryEscape(t.apiKey) +
			"&language=pt-BR&page=1&query=" + url.QueryEscape(mediaTitle))
		if err != nil {
			return nil, nil, err
		}
	}
	defer movieAPI.Body.Close()
	searchResults, err := ioutil.ReadAll(movieAPI.Body)
	if err != nil {
		return nil, nil, err
	}
	var theMovieResult types.MovieDBResponse
	err = json.Unmarshal(searchResults, &theMovieResult)
	if err != nil {
		return nil, nil, err
	}
	if mediaType == "movies" {
		var movies []types.Movie
		err = json.Unmarshal(theMovieResult.Data, &movies)
		if err != nil {
			return nil, nil, err
		}
		return movies, nil, nil
	}
	if mediaType == "tvshows" {
		var tvshows []types.TVShow
		err = json.Unmarshal(theMovieResult.Data, &tvshows)
		if err != nil {
			return nil, nil, err
		}
		return nil, tvshows, nil
	}
	return nil, nil, nil
}

func (t *theMovieDBAPI) SearchProviders(
	mediaType string,
	mediaID string,
) (*types.WatchProviders, error) {
	log.Println("providers api")
	var watchProviders *http.Response
	var err error
	if mediaType == "movies" {
		watchProviders, err = http.Get("https://api.themoviedb.org/3/movie/" +
			url.QueryEscape(mediaID) + "/watch/providers?api_key=" +
			url.QueryEscape(t.apiKey))
		if err != nil {
			return nil, err
		}
	}
	if mediaType == "tvshows" {
		watchProviders, err = http.Get("https://api.themoviedb.org/3/tv/" +
			url.QueryEscape(mediaID) + "/watch/providers?api_key=" +
			url.QueryEscape(t.apiKey))
		if err != nil {
			return nil, err
		}
	}
	defer watchProviders.Body.Close()
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		return nil, err
	}
	var providers types.WatchProviders
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		return nil, err
	}
	return &providers, nil
}

func (t *theMovieDBAPI) GetDetails(
	mediaType string,
	mediaID string,
) (*types.MovieDetails, *types.TVShowDetails, error) {
	log.Println("gettin details")
	var resDetails *http.Response
	var err error
	if mediaType == "movies" {
		resDetails, err = http.Get("https://api.themoviedb.org/3/movie/" + url.QueryEscape(mediaID) +
			"?api_key=" + url.QueryEscape(t.apiKey) + "&language=pt_BR")
		if err != nil {
			return nil, nil, err
		}
	}
	if mediaType == "tvshows" {
		resDetails, err = http.Get("https://api.themoviedb.org/3/tv/" + url.QueryEscape(mediaID) +
			"?api_key=" + url.QueryEscape(t.apiKey) + "&language=pt-BR")
		if err != nil {
			return nil, nil, err
		}
	}
	defer resDetails.Body.Close()
	details, err := ioutil.ReadAll(resDetails.Body)
	if err != nil {
		return nil, nil, err
	}
	if mediaType == "movies" {
		var movDetails types.MovieDetails
		err = json.Unmarshal(details, &movDetails)
		if err != nil {
			return nil, nil, err
		}
		return &movDetails, nil, err

	}
	if mediaType == "tvshows" {
		var tvShowDetails types.TVShowDetails
		err = json.Unmarshal(details, &tvShowDetails)
		if err != nil {
			return nil, nil, err
		}
		return nil, &tvShowDetails, err

	}
	return nil, nil, nil
}

func (t *theMovieDBAPI) GetMovieCredits(movieID string) (*types.MovieCredits, error) {
	log.Println("getting credits")
	movieCredits, err := http.Get("https://api.themoviedb.org/3/movie/" + url.QueryEscape(movieID) +
		"/credits?api_key=" + url.QueryEscape(t.apiKey))
	if err != nil {
		return nil, err
	}
	defer movieCredits.Body.Close()
	movCredits, err := ioutil.ReadAll(movieCredits.Body)
	var credits types.MovieCredits
	err = json.Unmarshal(movCredits, &credits)
	if err != nil {
		return nil, err
	}
	return &credits, nil
}

type movieAPICached struct {
	api   SearchMedia
	cache interface{}
	redis *redis.Client
}

func (t *movieAPICached) SearchMedia(mediaType string, mediaTitle string) ([]types.Movie, []types.TVShow, error) {
	log.Println("moviedb api cached")
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
				var movieMedia []types.Movie
				var tvShowMedia []types.TVShow
				data, err := t.redis.Get(key).Bytes()
				if err != nil {
					return nil, nil, err
				}
				if mediaType == "movies" {
					err = json.Unmarshal(data, &movieMedia)
					if err != nil {
						return nil, nil, err
					}
					t.cache = movieMedia
				}
				if mediaType == "tvshows" {
					err = json.Unmarshal(data, &tvShowMedia)
					if err != nil {
						return nil, nil, err
					}
					t.cache = tvShowMedia
				}
			}
		}
	}
	var resJSON []byte
	if mediaType == "movies" {
		if t.cache != nil {
			return t.cache.([]types.Movie), nil, nil
		}
		resMovies, _, err := t.api.SearchMedia(mediaType, mediaTitle)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resMovies)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":" + mediaTitle
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resMovies
		return resMovies, nil, nil
	}
	if mediaType == "tvshows" {
		if t.cache != nil {
			return nil, t.cache.([]types.TVShow), nil
		}
		_, resTVShow, err := t.api.SearchMedia(mediaType, mediaTitle)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resTVShow)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":" + mediaTitle
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resTVShow
		return nil, resTVShow, nil
	}
	return nil, nil, nil
}

type providersCached struct {
	api   SearchProviders
	cache interface{}
	redis *redis.Client
}

func (t *providersCached) SearchProviders(mediaType string, mediaID string) (*types.WatchProviders, error) {
	log.Println("providers cached")
	cfg, err := config.Wire()
	if err != nil {
		return nil, err
	}
	t.redis, err = r.Wire(cfg)
	if err != nil {
		return nil, err
	}
	keys, err := t.redis.Keys("telegram:" + mediaType + ":providers:*").Result()
	if err != nil {
		return nil, err
	}
	if len(keys) != 0 {
		for _, key := range keys {
			if strings.TrimPrefix(key, "telegram:"+mediaType+":providers:") == mediaID {
				var providers *types.WatchProviders
				data, err := t.redis.Get(key).Bytes()
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(data, &providers)
				if err != nil {
					return nil, err
				}
				t.cache = providers
			}
		}
	}
	if t.cache != nil {
		return t.cache.(*types.WatchProviders), nil
	}
	res, err := t.api.SearchProviders(mediaType, mediaID)
	if err != nil {
		return nil, err
	}
	t.cache = res
	var resJSON []byte
	resJSON, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}
	key := "telegram:" + mediaType + ":providers:" + mediaID
	if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
		return nil, err
	}
	return res, nil
}

type detailsCached struct {
	api   GetDetails
	cache interface{}
	redis *redis.Client
}

func (t *detailsCached) GetDetails(mediaType string, mediaID string) (*types.MovieDetails, *types.TVShowDetails, error) {
	log.Println("details cached")
	cfg, err := config.Wire()
	if err != nil {
		return nil, nil, err
	}
	t.redis, err = r.Wire(cfg)
	if err != nil {
		return nil, nil, err
	}
	keys, err := t.redis.Keys("telegram:" + mediaType + ":details:*").Result()
	if err != nil {
		return nil, nil, err
	}
	if len(keys) != 0 {
		for _, key := range keys {
			if strings.TrimPrefix(key, "telegram:"+mediaType+":details:") == mediaID {
				var movieDetails *types.MovieDetails
				var tvShowDetails *types.TVShowDetails
				data, err := t.redis.Get(key).Bytes()
				if err != nil {
					return nil, nil, err
				}
				if mediaType == "movies" {
					err = json.Unmarshal(data, &movieDetails)
					if err != nil {
						return nil, nil, err
					}
					t.cache = movieDetails
				}
				if mediaType == "tvshows" {
					err = json.Unmarshal(data, &tvShowDetails)
					if err != nil {
						return nil, nil, err
					}
					t.cache = tvShowDetails
				}
			}
		}
	}
	var resJSON []byte
	if mediaType == "movies" {
		if t.cache != nil {
			return t.cache.(*types.MovieDetails), nil, nil
		}
		resMovies, _, err := t.api.GetDetails(mediaType, mediaID)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resMovies)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":details:" + mediaID
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resMovies
		return resMovies, nil, nil
	}
	if mediaType == "tvshows" {
		if t.cache != nil {
			return nil, t.cache.(*types.TVShowDetails), nil
		}
		_, resTVShow, err := t.api.GetDetails(mediaType, mediaID)
		if err != nil {
			return nil, nil, err
		}
		resJSON, err = json.Marshal(resTVShow)
		if err != nil {
			return nil, nil, err
		}
		key := "telegram:" + mediaType + ":details:" + mediaID
		if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
			return nil, nil, err
		}
		t.cache = resTVShow
		return nil, resTVShow, nil
	}
	return nil, nil, nil
}

type moviesCreditCached struct {
	api   GetMovieCredits
	cache interface{}
	redis *redis.Client
}

func (t *moviesCreditCached) GetMovieCredits(movieID string) (*types.MovieCredits, error) {
	log.Println("credits cached")
	cfg, err := config.Wire()
	if err != nil {
		return nil, err
	}
	t.redis, err = r.Wire(cfg)
	if err != nil {
		return nil, err
	}
	keys, err := t.redis.Keys("telegram:movies:credits:*").Result()
	if err != nil {
		return nil, err
	}
	if len(keys) != 0 {
		for _, key := range keys {
			if strings.TrimPrefix(key, "telegram:movies:credits:") == movieID {
				var credits types.MovieCredits
				data, err := t.redis.Get(key).Bytes()
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(data, &credits)
				if err != nil {
					return nil, err
				}
				t.cache = credits
			}
		}
	}
	if t.cache != nil {
		return t.cache.(*types.MovieCredits), nil
	}
	res, err := t.api.GetMovieCredits(movieID)
	if err != nil {
		return nil, err
	}
	t.cache = res
	var resJSON []byte
	resJSON, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}
	key := "telegram:movies:credits" + movieID
	if err = t.redis.Set(key, resJSON, 72*time.Hour).Err(); err != nil {
		return nil, err
	}
	return res, nil
}
