package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/ismaelpereira/telegram-bot-isma/config"
)

func Wire(cfg *config.Config) (*redis.Client, error) {
	address := cfg.RedisAddress.Address
	opt, err := redis.ParseURL(address)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
