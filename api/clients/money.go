package clients

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type MoneyApi struct {
	ApiKey string
}

func (t *MoneyApi) GetCurrencies() (*types.MoneySearchResult, error) {
	apiResponse, err := http.Get("http://data.fixer.io/api/latest?access_key=" + url.QueryEscape(t.ApiKey))
	if err != nil {
		return nil, err
	}
	defer apiResponse.Body.Close()
	moneyValues, err := ioutil.ReadAll(apiResponse.Body)
	if err != nil {
		return nil, err
	}
	var moneyCurrencies types.MoneySearchResult
	err = json.Unmarshal(moneyValues, &moneyCurrencies)
	if err != nil {
		return nil, err
	}
	return &moneyCurrencies, nil
}
