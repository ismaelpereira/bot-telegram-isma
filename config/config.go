package config

import (
	"encoding/json"
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

func GetAdmiralPath() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	pathAdmiral := configDecoded.AdmiralPath.Path
	return pathAdmiral, err
}

func GetTelegramKey() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
	}

	telegramKey := configDecoded.Telegram.Key
	return telegramKey, err
}

func GetMoneyApiKey() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		return "", err
	}
	accessKey := configDecoded.MoneyAcessKey.Key
	return accessKey, err
}

func GetMovieApiKey() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		return "", err
	}
	acessKey := configDecoded.MovieAcessKey.Key
	return acessKey, err
}
