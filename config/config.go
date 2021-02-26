package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/gomodule/redigo/redis"
)

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

func Load() *Config {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return nil
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	return &configDecoded
}

func StartRedis() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
