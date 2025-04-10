package redis

import (
	"context"
	"dankmuzikk-web/config"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "dankmuzikk-web:"

const (
	redirectPathTtlMinutes = 30
)

type Cache struct {
	client *redis.Client
}

func New() *Cache {
	return &Cache{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Env().Cache.Host,
			Password: config.Env().Cache.Password,
			DB:       0,
		}),
	}
}

func redirectPathKey(clientHash string) string {
	return fmt.Sprintf("%sredirect-path:%s", keyPrefix, clientHash)
}

func (c *Cache) SetRedirectPath(clientHash, path string) error {
	return c.client.Set(context.Background(), redirectPathKey(clientHash), path, redirectPathTtlMinutes*time.Minute).Err()
}

func (c *Cache) GetRedirectPath(clientHash string) (string, error) {
	value, err := c.client.Get(context.Background(), redirectPathKey(clientHash)).Result()
	if err == redis.Nil {
		return "", errors.New("oopsie")
	} else if err != nil {
		return "", err
	}

	return value, nil
}
