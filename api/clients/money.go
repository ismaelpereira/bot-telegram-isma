package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ismaelpereira/telegram-bot-isma/types"
)

type moneyAPI struct {
	apiKey string
}
type MoneyAPI interface {
	GetCurrencies() (*types.MoneySearchResult, error)
}

func NewMoneyAPI(apiKey string) (MoneyAPI, error) {
	return &moneyAPICached{
		api: &moneyAPI{
			apiKey: apiKey,
		},
	}, nil
}

func (t *moneyAPI) GetCurrencies() (*types.MoneySearchResult, error) {
	log.Println("money api")
	apiResponse, err := http.Get("http://data.fixer.io/api/latest?access_key=" + url.QueryEscape(t.apiKey))
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

type moneyAPICached struct {
	api   MoneyAPI
	cache interface{}
}

func (t *moneyAPICached) GetCurrencies() (*types.MoneySearchResult, error) {
	log.Println("money api cached")
	if t.cache != nil {
		return t.cache.(*types.MoneySearchResult), nil
	}
	res, err := t.api.GetCurrencies()
	if err != nil {
		return nil, err
	}
	t.cache = res
	return res, nil
}
