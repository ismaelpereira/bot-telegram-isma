package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

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

func NewSearchMedia(mediaType string, mediaID string, apiKey string) (SearchMedia, error) {
	return &theMovieDBAPI{
		apiKey: apiKey,
	}, nil
}

type SearchProviders interface {
	SearchProviders(string, string) (*types.WatchProviders, error)
}

func NewSearchProviders(mediaType string, mediaID string, apiKey string) (SearchProviders, error) {
	return &theMovieDBAPI{
		apiKey: apiKey,
	}, nil
}

type GetDetails interface {
	GetDetails(string, string) (*types.MovieDetails, *types.TVShowDetails, error)
}

func NewGetDetails(mediaType string, mediaID string, apiKey string) (GetDetails, error) {
	return &theMovieDBAPI{
		apiKey: apiKey,
	}, nil
}

type GetMovieCredits interface {
	GetMovieCredits(string) (*types.MovieCredits, error)
}

func NewGetMovieCredits(movieID string, apiKey string) (GetMovieCredits, error) {
	return &theMovieDBAPI{
		apiKey: apiKey,
	}, nil
}

func (t *theMovieDBAPI) SearchMedia(
	mediaTitle string,
	mediaType string,
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
