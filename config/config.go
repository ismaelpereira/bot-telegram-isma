package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

func GetAdmiralPath(c *Config) (string, error) {
	pathAdmiral := c.AdmiralPath.Path
	err := fmt.Errorf("Cannot load admiral path")
	return pathAdmiral, err
}

func GetTelegramKey(c *Config) (string, error) {
	telegramKey := c.Telegram.Key
	err := fmt.Errorf("Cannot load Telegram KEY")
	return telegramKey, err
}

func GetMoneyApiKey(c *Config) (string, error) {
	accessKey := c.MoneyAcessKey.Key
	err := fmt.Errorf("Cannot load Money API KEY")
	return accessKey, err
}

func GetMovieApiKey(c *Config) (string, error) {
	acessKey := c.MovieAcessKey.Key
	err := fmt.Errorf("CCannot load Movie API KEY")
	return acessKey, err
}
