package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type admiralJSON struct {
	path string
}

type AdmiralJSON interface {
	GetAdmiral(string) ([]types.Admiral, error)
}

func NewAdmiral(path string) (AdmiralJSON, error) {
	return &admiralCached{
		api: &admiralJSON{
			path: path,
		},
	}, nil
}

func (t *admiralJSON) GetAdmiral(archivePath string) ([]types.Admiral, error) {
	log.Println("admiral")
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	admiralsArchive, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var admiralsDecoded []types.Admiral
	err = json.Unmarshal(admiralsArchive, &admiralsDecoded)
	if err != nil {
		return nil, err
	}
	return admiralsDecoded, nil
}

type admiralCached struct {
	api   AdmiralJSON
	cache interface{}
}

func (t *admiralCached) GetAdmiral(archivePath string) ([]types.Admiral, error) {
	log.Println("admiral cached")
	if t.cache != nil {
		return t.cache.([]types.Admiral), nil
	}
	res, err := t.api.GetAdmiral(archivePath)
	if err != nil {
		return nil, err
	}
	t.cache = res
	return res, nil
}
