package cache_test

import (
	"testing"
	"time"

	"github.com/ismaelpereira/telegram-bot-isma/cache"
	"github.com/ismaelpereira/telegram-bot-isma/types"
	"github.com/stretchr/testify/assert"
)

type CacheInterface interface {
	Get(now time.Time) interface{}
	Set(expireAt time.Time, data interface{})
}

func TestMoneySearchResultCache(t *testing.T) {
	var c CacheInterface = &cache.Cache{}
	data := &types.MoneySearchResult{}
	now := time.Now().Add(-24 * time.Hour)
	assert.Nil(t, c.Get(now))
	c.Set(now.Add(time.Hour), data)
	assert.Same(t, data, c.Get(now))
	assert.Same(t, data, c.Get(now.Add(-time.Hour)))
	assert.Nil(t, c.Get(now.Add(time.Hour+1)))
}
