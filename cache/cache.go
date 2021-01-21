package cache

import (
	"time"
)

type Cache struct {
	data     interface{}
	expireAt time.Time
}

func (t *Cache) Get(now time.Time) interface{} {
	if t.expireAt.Before(now) {
		return nil
	}
	return t.data
}

func (t *Cache) Set(expireAt time.Time, data interface{}) {
	t.expireAt = expireAt
	t.data = data
}
