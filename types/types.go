package types

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
