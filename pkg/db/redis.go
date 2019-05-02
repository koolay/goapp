package db

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type Redis struct {
	Client *redis.Client
	Prefix string
}

func NewRedis(endpoint string, password string, prefix string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         endpoint,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Password:     password,
		PoolSize:     4,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, errors.Errorf("Failed to ping redis. with endpoint '%s'", endpoint)
	}

	return &Redis{
		Client: client,
		Prefix: prefix,
	}, nil
}
