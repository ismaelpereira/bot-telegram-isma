package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EditMediaJSON struct {
	ChatID      int64                         `json:"chat_id"`
	MessageID   int                           `json:"message_id"`
	Media       Media                         `json:"media"`
	ReplyMarkup tgbotapi.InlineKeyboardMarkup `json:"reply_markup"`
}

type Media struct {
	Type    string `json:"type"`
	URL     string `json:"media"`
	Caption string `json:"caption"`
}

type Admiral struct {
	RealName        string
	AdmiralName     string
	AkumaNoMi       string
	Animal          string
	Power           string
	Sign            string
	ActorWhoInspire string
	BirthDate       string
	Height          float64
	Age             int
	ProfilePicture  string
}

type AnimeResponse struct {
	LastPage           int
	RequestCacheExpiry int
	RequestCached      bool
	RequestHash        string
	Results            []Anime
}
type Anime struct {
	ID           int `json:"mal_id"`
	Title        string
	Airing       bool
	Episodes     int
	CoverPicture string `json:"image_url"`
	Score        float64
}
type MangaResponse struct {
	LastPage           int
	RequestCacheExpiry int
	RequestCached      bool
	RequestHash        string
	Results            []Manga
}

type Manga struct {
	ID           int `json:"mal_id"`
	Title        string
	Publishing   bool
	Chapters     int
	Volumes      int
	Score        float64
	Status       string
	CoverPicture string `json:"image_url"`
	JapaneseName []byte
}

type MoneySearchResult struct {
	Success   bool   `json:"sucess"`
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     map[string]float64
}

type MovieResponse struct {
	Page    int
	Results []Movie
}

type Movie struct {
	ID            int
	Title         string
	OriginalTitle string  `json:"original_title"`
	ReleaseDate   string  `json:"release_date"`
	PosterPath    string  `json:"poster_path"`
	Popularity    float64 `json:"popularity"`
	Providers     WatchProvidersResponse
	Details       MovieDetails
	Credits       MovieCredits
}

type MovieDetails struct {
	Duration int     `json:"runtime"`
	Rating   float64 `json:"vote_average"`
}

type MovieCredits struct {
	ID   int
	Crew []Crew
}

type Crew struct {
	Name       string
	Department string
	Job        string
}

type WatchProvidersResponse struct {
	ID      int
	Results map[string]*CountryOptions
}

type CountryOptions struct {
	Link     string
	Rent     []*Provider
	Buy      []*Provider
	Flatrate []*Provider
}

type Provider struct {
	DisplayPriority int    `json:"display_priority"`
	ProviderID      int    `json:"provider_id"`
	ProviderName    string `json:"provider_name"`
}

type TVShowResponse struct {
	Page    int
	Results []TVShow
}

type TVShow struct {
	ID            int
	Title         string `json:"name"`
	OriginalTitle string `json:"original_name"`
	Popularity    float64
	PosterPath    string `json:"poster_path"`
	ReleaseDate   string `json:"first_air_date"`
	TVShowDetails TVShowDetails
	Providers     WatchProvidersResponse
}

type TVShowDetails struct {
	Status       string
	SeasonNumber int `json:"number_of_seasons"`
	Seasons      []TVShowSeasonDetails
	Rating       float64 `json:"vote_average"`
	CreatedBy    []Crew  `json:"created_by"`
}

type TVShowSeasonDetails struct {
	EpisodesCount int `json:"episode_count"`
	Name          string
	AirDate       string `json:"air_date"`
	PosterPath    string `json:"poster_path"`
}
