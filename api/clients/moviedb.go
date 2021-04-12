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
	"github.com/ismaelpereira/telegram-bot-isma/api/common"
	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type theMovieDBAPI struct {
	apiKey string
}

type SearchMedia interface {
	SearchMedia(string, string) (interface{}, error)
}

func NewSearchMedia(mediaType string, apiKey string) (SearchMedia, error) {
	return &movieAPICached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
		cache: nil,
		redis: nil,
	}, nil
}

type SearchProviders interface {
	SearchProviders(string, string) (*types.WatchProviders, error)
}

func NewSearchProviders(mediaType string, apiKey string) (SearchProviders, error) {
	return &providersCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
		cache: nil,
		redis: nil,
	}, nil
}

type GetDetails interface {
	GetDetails(string, string) (*types.MovieDetails, *types.TVShowDetails, error)
}

func NewGetDetails(mediaType string, apiKey string) (GetDetails, error) {
	return &detailsCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
		cache: nil,
		redis: nil,
	}, nil
}

type GetMovieCredits interface {
	GetMovieCredits(string) (*types.MovieCredits, error)
}

func NewGetMovieCredits(apiKey string) (GetMovieCredits, error) {
	return &moviesCreditCached{
		api: &theMovieDBAPI{
			apiKey: apiKey,
		},
		cache: nil,
		redis: nil,
	}, nil
}

func (t *theMovieDBAPI) SearchMedia(
	mediaType string,
	mediaTitle string,
) (interface{}, error) {
	log.Println("moviedb api")
	if mediaType == "movies" {
		var movies []types.Movie
		url := "https://api.themoviedb.org/3/search/movie?api_key=" +
			url.QueryEscape(t.apiKey) +
			"&page=1&langague=pt-br&query=" + url.QueryEscape(mediaTitle)
		if err := t.httpGET(url, &movies); err != nil {
			return nil, err
		}
		return movies, nil
	}
	if mediaType == "tvshows" {
		var tvShows []types.TVShow
		url := "https://api.themoviedb.org/3/search/tv?api_key=" +
			url.QueryEscape(t.apiKey) +
			"&page=1&langague=pt-br&query=" + url.QueryEscape(mediaTitle)
		if err := t.httpGET(url, &tvShows); err != nil {
			return nil, err
		}
		return tvShows, nil
	}
	return nil, nil
}

func (t *theMovieDBAPI) httpGET(url string, v interface{}) error {
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
	var theMovieResult types.MovieDBResponse
	if err = json.Unmarshal(searchResult, &theMovieResult); err != nil {
		return err
	}
	if err = json.Unmarshal(theMovieResult.Data, &v); err != nil {
		return err
	}
	return nil
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
		defer watchProviders.Body.Close()
	}
	if mediaType == "tvshows" {
		watchProviders, err = http.Get("https://api.themoviedb.org/3/tv/" +
			url.QueryEscape(mediaID) + "/watch/providers?api_key=" +
			url.QueryEscape(t.apiKey))
		if err != nil {
			return nil, err
		}
		defer watchProviders.Body.Close()
	}
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
	if mediaType == "movies" {
		var details *types.MovieDetails
		url := "https://api.themoviedb.org/3/movie/" + url.QueryEscape(mediaID) +
			"?api_key=" + url.QueryEscape(t.apiKey) + "&language=pt_BR"
		if err := t.httpGETDetails(url, &details); err != nil {
			return nil, nil, err
		}
		return details, nil, nil
	}
	if mediaType == "tvshows" {
		var details *types.TVShowDetails
		url := "https://api.themoviedb.org/3/tv/" + url.QueryEscape(mediaID) +
			"?api_key=" + url.QueryEscape(t.apiKey) + "&language=pt_BR"
		if err := t.httpGETDetails(url, &details); err != nil {
			return nil, nil, err
		}
		return nil, details, nil
	}
	return nil, nil, nil
}

func (t *theMovieDBAPI) httpGETDetails(url string, v interface{}) error {
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
	err = json.Unmarshal(searchResult, &v)
	if err != nil {
		return err
	}
	return nil
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
	if err != nil {
		return nil, err
	}
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

func (t *movieAPICached) SearchMedia(mediaType string, mediaTitle string) (interface{}, error) {
	log.Println("moviedb api cached")
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return nil, err
	}
	t.cache, err = common.SearchItens(t.redis, mediaType, mediaTitle)
	if err != nil {
		return nil, err
	}
	if mediaType == "movies" {
		if t.cache != nil {
			return t.cache.([]types.Movie), nil
		}
	}
	if mediaType == "tvshows" {
		if t.cache != nil {
			return t.cache.([]types.TVShow), nil
		}
	}
	var res interface{}
	if res, err = t.setMediaInRedis(mediaType, mediaTitle); err != nil {
		return nil, err
	}
	t.cache = res
	return res, nil
}

func (t *movieAPICached) setMediaInRedis(mediaType string, mediaTitle string) (interface{}, error) {
	res, err := t.api.SearchMedia(mediaType, mediaTitle)
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

type providersCached struct {
	api   SearchProviders
	cache interface{}
	redis *redis.Client
}

func (t *providersCached) SearchProviders(mediaType string, mediaID string) (*types.WatchProviders, error) {
	log.Println("providers cached")
	err := t.searchRedisProvidersKeys(mediaType, mediaID)
	if err != nil {
		return nil, err
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

func (t *providersCached) searchRedisProvidersKeys(mediaType string, mediaID string) error {
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return err
	}
	keys, err := t.redis.Keys("telegram:" + mediaType + ":providers:*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		if strings.TrimPrefix(key, "telegram:"+mediaType+":providers:") == mediaID {
			var providers *types.WatchProviders
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &providers)
			if err != nil {
				return err
			}
			t.cache = providers
			return nil
		}
	}
	t.cache = nil
	return nil
}

type detailsCached struct {
	api   GetDetails
	cache interface{}
	redis *redis.Client
}

func (t *detailsCached) GetDetails(
	mediaType string,
	mediaID string,
) (*types.MovieDetails, *types.TVShowDetails, error) {
	log.Println("details cached")
	err := t.searchDetailsKeys(mediaType, mediaID)
	if err != nil {
		return nil, nil, err
	}
	if mediaType == "movies" {
		if t.cache != nil {
			return t.cache.(*types.MovieDetails), nil, nil
		}
		res, _, err := t.api.GetDetails(mediaType, mediaID)
		if err != nil {
			return nil, nil, err
		}
		t.cache = res
		return res, nil, err
	}
	if mediaType == "tvshows" {
		if t.cache != nil {
			return nil, t.cache.(*types.TVShowDetails), nil
		}
		_, res, err := t.api.GetDetails(mediaType, mediaID)
		if err != nil {
			return nil, nil, err
		}
		t.cache = res
		return nil, res, err
	}
	return nil, nil, nil
}

func (t *detailsCached) searchDetailsKeys(mediaType string, mediaID string) error {
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return err
	}
	keys, err := t.redis.Keys("telegram:" + mediaType + ":details:*").Result()
	if err != nil {
		return err
	}
	for _, key := range keys {
		if strings.TrimPrefix(key, "telegram:"+mediaType+":details:") != mediaID {
			continue
		}
		var movieDetails types.MovieDetails
		var tvShowDetails types.TVShowDetails
		data, err := t.redis.Get(key).Bytes()
		if err != nil {
			return err
		}
		if mediaType == "movies" {
			if err = json.Unmarshal(data, &movieDetails); err != nil {
				return err
			}
			t.cache = movieDetails
			return nil
		}
		if mediaType == "tvshows" {
			if err = json.Unmarshal(data, &tvShowDetails); err != nil {
				return err
			}
			t.cache = tvShowDetails
			return nil
		}
	}
	t.cache = nil
	return nil
}

type moviesCreditCached struct {
	api   GetMovieCredits
	cache interface{}
	redis *redis.Client
}

func (t *moviesCreditCached) GetMovieCredits(movieID string) (*types.MovieCredits, error) {
	log.Println("credits cached")
	err := t.getCreditKeys(movieID)
	if err != nil {
		return nil, err
	}
	if t.cache != nil {
		return t.cache.(*types.MovieCredits), nil
	}
	res, err := t.api.GetMovieCredits(movieID)
	if err != nil {
		return nil, err
	}
	var resJSON []byte
	resJSON, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = common.SetRedisKeyMovieCredits(resJSON, t.redis, movieID)
	if err != nil {
		return nil, err
	}
	t.cache = res
	return res, nil
}

func (t *moviesCreditCached) getCreditKeys(movieID string) error {
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return err
	}
	keys, err := t.redis.Keys("telegram:movies:credits:*").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		if strings.TrimPrefix(key, "telegram:movies:credits:") == movieID {
			var credits *types.MovieCredits
			data, err := t.redis.Get(key).Bytes()
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, &credits)
			if err != nil {
				return err
			}
			t.cache = credits
			return nil
		}
	}
	t.cache = nil
	return nil
}
