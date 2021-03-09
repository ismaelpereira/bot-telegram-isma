package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func Wire() (*Config, error) {
	file, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer file.Close()
	configEncoded, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var configDecoded Config
	err = json.Unmarshal(configEncoded, &configDecoded)
	if err != nil {
		return nil, err
	}
	return &configDecoded, nil
}
