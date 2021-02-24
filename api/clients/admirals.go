package clients

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type AdmiralJSON struct {
	Path string
}

func (t *AdmiralJSON) GetAdmiral(archivePath string) ([]types.Admiral, error) {
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
