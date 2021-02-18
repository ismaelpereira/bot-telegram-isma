package clients

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/IsmaelPereira/telegram-bot-isma/types"
)

type MovieDB struct {
	ApiKey string
}

type TVShowDB struct {
	ApiKey string
}

func (t *MovieDB) SearchMovie(movieName string) (*types.MovieResponse, error) {
	movieAPI, err := http.Get("https://api.themoviedb.org/3/search/movie?api_key=" + url.QueryEscape(t.ApiKey) +
		"&page=1&langague=pt-br&query=" + url.QueryEscape(movieName))
	if err != nil {
		return nil, err
	}
	defer movieAPI.Body.Close()
	searchResults, err := ioutil.ReadAll(movieAPI.Body)
	if err != nil {
		return nil, err
	}
	var moviesResults types.MovieResponse
	err = json.Unmarshal(searchResults, &moviesResults)
	if err != nil {
		return nil, err
	}
	return &moviesResults, nil
}

func (t *MovieDB) GetMovieProviders(movieID string) (*types.WatchProvidersResponse, error) {
	watchProviders, err := http.Get("https://api.themoviedb.org/3/movie/" +
		url.QueryEscape(movieID) + "/watch/providers?api_key=" +
		url.QueryEscape(t.ApiKey))
	if err != nil {
		return nil, err
	}
	defer watchProviders.Body.Close()
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		return nil, err
	}
	var providers types.WatchProvidersResponse
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		return nil, err
	}
	return &providers, nil
}

func (t *TVShowDB) SearchTVShow(TVShowTitle string) (*types.TVShowResponse, error) {
	TVShowAPI, err := http.Get("https://api.themoviedb.org/3/search/tv?api_key=" + url.QueryEscape(t.ApiKey) +
		"&language=pt-BR&page=1&include_adult=false&query=" + url.QueryEscape(TVShowTitle))
	if err != nil {
		return nil, err
	}
	defer TVShowAPI.Body.Close()
	searchResults, err := ioutil.ReadAll(TVShowAPI.Body)
	if err != nil {
		return nil, err
	}
	var TVShowsResults types.TVShowResponse
	err = json.Unmarshal(searchResults, &TVShowsResults)
	return &TVShowsResults, nil
}

func (t *TVShowDB) GetTVShowSeasonDetails(TVShowID string) (*types.TVShowDetails, error) {
	seasonDetails, err := http.Get("https://api.themoviedb.org/3/tv/" + url.QueryEscape(TVShowID) +
		"?api_key=" + url.QueryEscape(t.ApiKey) + "&language=pt-BR")
	if err != nil {
		return nil, err
	}
	defer seasonDetails.Body.Close()
	seasonDetailsResults, err := ioutil.ReadAll(seasonDetails.Body)
	if err != nil {
		return nil, err
	}
	var seasonsResults types.TVShowDetails
	err = json.Unmarshal(seasonDetailsResults, &seasonsResults)
	if err != nil {
		return nil, err
	}
	return &seasonsResults, nil
}

func (t *TVShowDB) GetTVShowProviders(TVShowID string) (*types.WatchProvidersResponse, error) {
	watchProviders, err := http.Get("https://api.themoviedb.org/3/tv/" +
		url.QueryEscape(TVShowID) + "/watch/providers?api_key=" +
		url.QueryEscape(t.ApiKey))
	if err != nil {
		return nil, err
	}
	defer watchProviders.Body.Close()
	providersValues, err := ioutil.ReadAll(watchProviders.Body)
	if err != nil {
		return nil, err
	}
	var providers types.WatchProvidersResponse
	err = json.Unmarshal(providersValues, &providers)
	if err != nil {
		return nil, err
	}
	return &providers, nil
}
