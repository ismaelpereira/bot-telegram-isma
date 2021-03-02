package config

type Config struct {
	Telegram      TelegramKey `json:"telegram"`
	AdmiralPath   AdmiralPath `json:"admiral"`
	MoneyAcessKey MoneyAPI    `json:"moneyAPI"`
	MovieAcessKey MovieAPI    `json:"movieAPI"`
	RedisAddress  Redis       `json:"redis"`
}

type TelegramKey struct {
	Key string `json:"key"`
}

type AdmiralPath struct {
	Path string `json:"path"`
}

type MoneyAPI struct {
	Key string `json:"accessKey"`
}

type MovieAPI struct {
	Key string `json:"acessKey"`
}

type Redis struct {
	Address string `json:"address"`
}
