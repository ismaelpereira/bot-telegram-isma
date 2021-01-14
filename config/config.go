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
}

type TelegramKey struct {
	Key string `json:"key"`
}

type AdmiralPath struct {
	Path string `json:"path"`
}

func GetAdmiralPath() string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	pathAdmiral := configDecoded.AdmiralPath.Path
	return pathAdmiral
}

func GetTelegramKey() string {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		log.Println(err)
	}

	telegramKey := configDecoded.Telegram.Key
	return telegramKey
}
