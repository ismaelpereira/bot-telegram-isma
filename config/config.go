package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Telegram    TelegramKey `json:"telegram"`
	AdmiralPath AdmiralPath `json:"admiral"`
	StoragePath Storage     `json:"storage"`
}

type TelegramKey struct {
	Key string `json:"key"`
}

type AdmiralPath struct {
	Path string `json:"path"`
}

type Storage struct {
	Path string `json:"path"`
}

func GetAdmiralPath() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		log.Println(err)
	}

	telegramKey := configDecoded.Telegram.Key
	return telegramKey, err
}

func GetStoragePath() (string, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		log.Println(err)
		return "", err
	}
	storagePath := configDecoded.StoragePath.Path
	return storagePath, err
}
