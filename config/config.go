package config

import "github.com/gomodule/redigo/redis"

type Config struct {
	Telegram      TelegramKey `json:"telegram"`
	AdmiralPath   AdmiralPath `json:"admiral"`
	MoneyAcessKey MoneyAPI    `json:"moneyAPI"`
	MovieAcessKey MovieAPI    `json:"movieAPI"`
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

func StartRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
