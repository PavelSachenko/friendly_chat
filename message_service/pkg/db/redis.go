package db

import (
	"github.com/go-redis/redis/v7"
	"github.com/pavel/message_service/config"
)

func InitRedis(cfg *config.Config) (error, *redis.Client) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.DB.RedisDSN,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return err, nil
	}
	return nil, client
}
