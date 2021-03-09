package types

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type EditMediaJSON struct {
	ChatID      int64                         `json:"chat_id"`
	MessageID   int                           `json:"message_id"`
	Media       InputMedia                    `json:"media"`
	ReplyMarkup tgbotapi.InlineKeyboardMarkup `json:"reply_markup"`
}

type InputMedia struct {
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

type JikanResponse struct {
	LastPage           int             `json:"last_page"`
	RequestCacheExpiry int             `json:"request_cache_expiry"`
	RequestCached      bool            `json:"request_cached"`
	RequestHash        string          `json:"request_hash"`
	Data               json.RawMessage `json:"results"`
}
type Anime struct {
	ID           int `json:"mal_id"`
	Title        string
	Airing       bool
	Episodes     int
	CoverPicture string `json:"image_url"`
	Score        float64
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
	Success   bool
	Timestamp int64
	Base      string
	Date      string
	Rates     map[string]float64
}

type MovieDBResponse struct {
	Page int
	Data json.RawMessage `json:"results"`
}

type Movie struct {
	ID            int
	Title         string
	OriginalTitle string  `json:"original_title"`
	ReleaseDate   string  `json:"release_date"`
	PosterPath    string  `json:"poster_path"`
	Popularity    float64 `json:"popularity"`
	Providers     WatchProviders
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

type WatchProviders struct {
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

type TVShow struct {
	ID            int
	Title         string `json:"name"`
	OriginalTitle string `json:"original_name"`
	Popularity    float64
	PosterPath    string `json:"poster_path"`
	ReleaseDate   string `json:"first_air_date"`
	TVShowDetails TVShowDetails
	Providers     WatchProviders
}

type TVShowDetails struct {
	Status       string
	SeasonNumber int `json:"number_of_seasons"`
	Seasons      []TVShowSeasonDetails
	Rating       float64 `json:"vote_average"`
	CreatedBy    []Crew  `json:"created_by"`
}

type Crew struct {
	Name       string
	Department string
	Job        string
}

type TVShowSeasonDetails struct {
	EpisodesCount int `json:"episode_count"`
	Name          string
	AirDate       string `json:"air_date"`
	PosterPath    string `json:"poster_path"`
}

type Checklist struct {
	Title string
	Itens []ChecklistItem
}
type ChecklistItem struct {
	Name      string
	IsChecked bool
}
