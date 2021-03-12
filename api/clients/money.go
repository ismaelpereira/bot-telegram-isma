package clients

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/ismaelpereira/telegram-bot-isma/api/common"
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
		cache: nil,
		redis: nil,
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
	redis *redis.Client
}

func (t *moneyAPICached) GetCurrencies() (*types.MoneySearchResult, error) {
	log.Println("money api cached")
	err := t.getMoneyKeys()
	if err != nil {
		return nil, err
	}
	if t.cache != nil {
		return t.cache.(*types.MoneySearchResult), nil
	}
	res, err := t.api.GetCurrencies()
	if err != nil {
		return nil, err
	}
	t.cache = res
	key := "telegram:rates"
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	if err = t.redis.Set(key, resJSON, 1*time.Hour).Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (t *moneyAPICached) getMoneyKeys() error {
	var err error
	t.redis, err = common.SetRedis()
	if err != nil {
		return err
	}
	keys, err := t.redis.Keys("telegram:rates").Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		var rates *types.MoneySearchResult
		data, err := t.redis.Get(key).Bytes()
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &rates)
		if err != nil {
			return err
		}
		t.cache = rates
	}
	return err
}
