package redis

import (
	"context"
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "dankmuzikk:"

const (
	songLyricsTtlDays       = 7
	userSessionTokenTtlDays = 30
	otpTtlMinutes           = 30
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
	return fmt.Sprintf("%ssong:%d:lyrics", keyPrefix, songId)
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

func userTokenKey(sessionToken string) string {
	return fmt.Sprintf("%suser-session-token:%s", keyPrefix, sessionToken)
}

func (c *Cache) SetAuthenticatedUser(sessionToken string, profile models.Profile) error {
	profileJson, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return c.client.Set(context.Background(), userTokenKey(sessionToken), string(profileJson), userSessionTokenTtlDays*time.Hour*24).Err()
}

func (c *Cache) GetAuthenticatedUser(sessionToken string) (models.Profile, error) {
	res := c.client.Get(context.Background(), userTokenKey(sessionToken))
	if res == nil {
		return models.Profile{}, &app.ErrNotFound{
			ResourceName: "profile",
		}
	}
	value, err := res.Result()
	if err == redis.Nil {
		return models.Profile{}, &app.ErrNotFound{
			ResourceName: "profile",
		}
	} else if err != nil {
		return models.Profile{}, err
	}

	var profile models.Profile
	err = json.Unmarshal([]byte(value), &profile)
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (c *Cache) InvalidateAuthenticatedUser(sessionToken string) error {
	err := c.client.Del(context.Background(), userTokenKey(sessionToken)).Err()
	if err != nil {
		return err
	}

	return nil
}

func otpKey(userId uint) string {
	return fmt.Sprintf("%suser-otp:%d", keyPrefix, userId)
}

func (c *Cache) CreateOtp(accountId uint, otp string) error {
	err := c.client.Set(context.Background(), otpKey(accountId), otp, otpTtlMinutes*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetOtpForAccount(id uint) (string, error) {
	res := c.client.Get(context.Background(), otpKey(id))
	if res == nil {
		return "", &app.ErrNotFound{
			ResourceName: "otp",
		}
	}
	value, err := res.Result()
	if err == redis.Nil {
		return "", &app.ErrNotFound{
			ResourceName: "otp",
		}
	} else if err != nil {
		return "", err
	}

	return value, nil
}
