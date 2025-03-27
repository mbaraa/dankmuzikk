package redis

import (
	"context"
	"dankmuzikk/app"
	"dankmuzikk/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	songLyricsTtlDays = 7
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

func songLyricsKey(songId uint) string {
	return fmt.Sprintf("%d:lyrics", songId)
}

func (c *Cache) StoreLyrics(songId uint, lyrics []string) error {
	lyricsJson, err := json.Marshal(lyrics)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), songLyricsKey(songId), string(lyricsJson), songLyricsTtlDays*time.Hour*24).Err()
}

func (c *Cache) GetLyrics(songId uint) ([]string, error) {
	res := c.client.Get(context.Background(), songLyricsKey(songId))
	if res == nil {
		return nil, &app.ErrNotFound{
			ResourceName: "lyrics",
		}
	}
	value, err := res.Result()
	if err == redis.Nil {
		return nil, &app.ErrNotFound{
			ResourceName: "lyrics",
		}
	} else if err != nil {
		return nil, err
	}

	var lyrics []string
	err = json.Unmarshal([]byte(value), &lyrics)
	if err != nil {
		return nil, err
	}

	return lyrics, nil
}
