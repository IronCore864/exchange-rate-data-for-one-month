package redisclient

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/ironcore864/exchange-rate-data-for-one-month/config"
)

var client *redis.Client

func getClient() *redis.Client {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:        fmt.Sprintf("%s:%d", config.Conf.RedisHost, config.Conf.RedisPort),
			Password:    config.Conf.RedisPassword,
			DB:          0, // use default DB
			DialTimeout: 3 * time.Second,
			PoolSize:    10,
		})
	}
	return client
}

// Set a key
func Set(key string, data float32, expiration time.Duration) (string, error) {
	return getClient().Set(key, data, expiration).Result()
}
