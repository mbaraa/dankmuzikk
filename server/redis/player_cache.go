package redis

import (
	"context"
	"dankmuzikk/app"
	"dankmuzikk/app/models"
	"dankmuzikk/config"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	playerSettingsTtlHours = 72
	playerQueueTtlHours    = 2
)

var _ app.PlayerCache = &playerCache{}

type playerCache struct {
	client *redis.Client
}

func NewPlayerCache() *playerCache {
	return &playerCache{
		client: redis.NewClient(&redis.Options{
			Addr:     config.Env().Cache.Host,
			Password: config.Env().Cache.Password,
			DB:       0,
		}),
	}
}

type playerSetting string

const (
	playerSettingShuffle playerSetting = "shuffle"
	playerSettingLoop    playerSetting = "loop"
)

func playerQueueKey(accountId uint) string {
	return fmt.Sprintf("%splayer-queue:%d", keyPrefix, accountId)
}

func playerShuffledQueueKey(accountId uint) string {
	return fmt.Sprintf("%splayer-shuffled-queue:%d", keyPrefix, accountId)
}

// func playerQueueSongKey(songId uint) string {
// return fmt.Sprintf("%splayer-song:%d", keyPrefix, songId)
// }

func (c *playerCache) CreateSongsQueue(accountId uint, initialSongIds ...uint) error {
	for _, songId := range initialSongIds {
		err := c.AddSongToQueue(accountId, songId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *playerCache) CreateSongsShuffledQueue(accountId uint, initialSongIds ...uint) error {
	for _, songId := range initialSongIds {
		err := c.AddSongToShuffledQueue(accountId, songId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *playerCache) AddSongToQueue(accountId uint, songId uint) error {
	status := c.client.RPush(context.Background(),
		playerQueueKey(accountId),
		songId,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	res := c.client.Expire(context.Background(), playerQueueKey(accountId), playerQueueTtlHours*time.Hour)
	if res == nil {
		return nil
	}

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (c *playerCache) AddSongToQueueAfterIndex(accountId uint, songId uint, index int) error {
	res := c.client.LLen(context.Background(), playerQueueKey(accountId))
	if res == nil {
		return nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return err
	}

	if index > int(queueLen) || index < 0 {
		return errors.New("invalid song index")
	}

	res2 := c.client.LIndex(context.Background(), playerQueueKey(accountId), int64(index))
	if res2 == nil {
		return nil
	}

	err = res2.Err()
	if err != nil {
		return err
	}

	currentSong := res2.Val()

	res3 := c.client.LSet(context.Background(), playerQueueKey(accountId), int64(index), "replace-next")
	if res3 == nil {
		return nil
	}

	err = res3.Err()
	if err != nil {
		return err
	}

	res4 := c.client.LInsert(context.Background(), playerQueueKey(accountId), "BEFORE", "replace-next", currentSong)
	if res4 == nil {
		return nil
	}

	err = res4.Err()
	if err != nil {
		return err
	}

	res5 := c.client.LInsert(context.Background(), playerQueueKey(accountId), "AFTER", "replace-next", songId)
	if res5 == nil {
		return nil
	}

	err = res5.Err()
	if err != nil {
		return err
	}

	res6 := c.client.LRem(context.Background(), playerQueueKey(accountId), 1, "replace-next")
	if res6 == nil {
		return nil
	}

	err = res6.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) AddSongToShuffledQueue(accountId uint, songId uint) error {
	status := c.client.RPush(context.Background(),
		playerShuffledQueueKey(accountId),
		songId,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	res := c.client.Expire(context.Background(), playerQueueKey(accountId), playerQueueTtlHours*time.Hour)
	if res == nil {
		return nil
	}

	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (c *playerCache) AddSongToShuffledQueueAfterIndex(accountId uint, songId uint, index int) error {
	res := c.client.LLen(context.Background(), playerShuffledQueueKey(accountId))
	if res == nil {
		return nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return err
	}

	if index > int(queueLen) || index < 0 {
		return errors.New("invalid song index")
	}

	res2 := c.client.LIndex(context.Background(), playerShuffledQueueKey(accountId), int64(index))
	if res2 == nil {
		return nil
	}

	err = res2.Err()
	if err != nil {
		return err
	}

	currentSong := res2.Val()

	res3 := c.client.LSet(context.Background(), playerShuffledQueueKey(accountId), int64(index), "replace-next")
	if res3 == nil {
		return nil
	}

	err = res3.Err()
	if err != nil {
		return err
	}

	res4 := c.client.LInsert(context.Background(), playerShuffledQueueKey(accountId), "BEFORE", "replace-next", currentSong)
	if res4 == nil {
		return nil
	}

	err = res4.Err()
	if err != nil {
		return err
	}

	res5 := c.client.LInsert(context.Background(), playerShuffledQueueKey(accountId), "AFTER", "replace-next", songId)
	if res5 == nil {
		return nil
	}

	err = res5.Err()
	if err != nil {
		return err
	}

	res6 := c.client.LRem(context.Background(), playerShuffledQueueKey(accountId), 1, "replace-next")
	if res6 == nil {
		return nil
	}

	err = res6.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) RemoveSongFromQueue(songIndex int, accountId uint) error {
	res := c.client.LLen(context.Background(), playerQueueKey(accountId))
	if res == nil {
		return nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return err
	}

	if songIndex > int(queueLen) || songIndex < 0 {
		return errors.New("invalid song index")
	}

	res2 := c.client.LSet(context.Background(), playerQueueKey(accountId), int64(songIndex), "deleted")
	if res2 == nil {
		return nil
	}

	err = res2.Err()
	if err != nil {
		return err
	}

	res3 := c.client.LRem(context.Background(), playerQueueKey(accountId), 1, "deleted")
	if res3 == nil {
		return nil
	}

	err = res3.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) RemoveSongFromShuffledQueue(songIndex int, accountId uint) error {
	res := c.client.LLen(context.Background(), playerShuffledQueueKey(accountId))
	if res == nil {
		return nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return err
	}

	if songIndex > int(queueLen) || songIndex < 0 {
		return errors.New("invalid song index")
	}

	res2 := c.client.LSet(context.Background(), playerShuffledQueueKey(accountId), int64(songIndex), "deleted")
	if res2 == nil {
		return nil
	}

	err = res2.Err()
	if err != nil {
		return err
	}

	res3 := c.client.LRem(context.Background(), playerShuffledQueueKey(accountId), 1, "deleted")
	if res3 == nil {
		return nil
	}

	err = res3.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) ClearQueue(accountId uint) error {
	res := c.client.Del(context.Background(), playerQueueKey(accountId))
	if res == nil {
		return nil
	}
	_, err := res.Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) ClearShuffledQueue(accountId uint) error {
	res := c.client.Del(context.Background(), playerShuffledQueueKey(accountId))
	if res == nil {
		return nil
	}
	_, err := res.Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) GetSongsQueue(accountId uint) ([]uint, error) {
	res := c.client.LRange(context.Background(), playerQueueKey(accountId), 0, -1)
	if res == nil {
		return nil, &app.ErrNotFound{
			ResourceName: "player_queue",
		}
	}
	queueRaw, err := res.Result()
	if err != nil {
		return nil, err
	}

	queue := make([]uint, 0, len(queueRaw))
	for _, songIdStr := range queueRaw {
		songId, err := strconv.Atoi(songIdStr)
		if err != nil {
			continue
		}
		queue = append(queue, uint(songId))
	}

	return queue, nil
}

func (c *playerCache) GetSongsShuffledQueue(accountId uint) ([]uint, error) {
	res := c.client.LRange(context.Background(), playerShuffledQueueKey(accountId), 0, -1)
	if res == nil {
		return nil, &app.ErrNotFound{
			ResourceName: "player_queue",
		}
	}
	queueRaw, err := res.Result()
	if err != nil {
		return nil, err
	}

	queue := make([]uint, 0, len(queueRaw))
	for _, songIdStr := range queueRaw {
		songId, err := strconv.Atoi(songIdStr)
		if err != nil {
			continue
		}
		queue = append(queue, uint(songId))
	}

	return queue, nil
}

func playerCurrentPlayingSongKey(accountId uint, clientHash string) string {
	return fmt.Sprintf("%splayer-playing-index:%d:%s", keyPrefix, accountId, clientHash)
}

func (c *playerCache) SetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string, songIndex int) error {
	status := c.client.Set(context.Background(),
		playerCurrentPlayingSongKey(accountId, clientHash),
		songIndex,
		playerSettingsTtlHours*time.Hour,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	return nil

}

func (c *playerCache) GetCurrentPlayingSongIndexInQueue(accountId uint, clientHash string) (int, error) {
	res := c.client.Get(context.Background(), playerCurrentPlayingSongKey(accountId, clientHash))
	if res == nil {
		return 0, redis.Nil
	}
	value, err := res.Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func playerCurrentPlayingSongKeyInShuffledQueue(accountId uint, clientHash string) string {
	return fmt.Sprintf("%splayer-playing-index-shuffled:%d:%s", keyPrefix, accountId, clientHash)
}

func (c *playerCache) SetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string, songIndex int) error {
	status := c.client.Set(context.Background(),
		playerCurrentPlayingSongKeyInShuffledQueue(accountId, clientHash),
		songIndex,
		playerSettingsTtlHours*time.Hour,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	return nil

}

func (c *playerCache) GetCurrentPlayingSongIndexInShuffledQueue(accountId uint, clientHash string) (int, error) {
	res := c.client.Get(context.Background(), playerCurrentPlayingSongKeyInShuffledQueue(accountId, clientHash))
	if res == nil {
		return 0, redis.Nil
	}
	value, err := res.Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(value)
}

func playerSettingKey(accountId uint, settingName playerSetting) string {
	return fmt.Sprintf("%splayer-setting:%d:%s", keyPrefix, accountId, settingName)
}

func (c *playerCache) SetShuffled(accountId uint, shuffled bool) error {
	status := c.client.Set(context.Background(),
		playerSettingKey(accountId, playerSettingShuffle),
		shuffled,
		playerSettingsTtlHours*time.Hour,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *playerCache) GetShuffled(accountId uint) (bool, error) {
	res := c.client.Get(context.Background(), playerSettingKey(accountId, playerSettingShuffle))
	if res == nil {
		return false, nil
	}
	value, err := res.Result()
	if err != nil {
		return false, err
	}

	return value == "1", nil
}

func (c *playerCache) SetLoopMode(accountId uint, mode models.PlayerLoopMode) error {
	status := c.client.Set(context.Background(),
		playerSettingKey(accountId, playerSettingLoop),
		string(mode),
		playerSettingsTtlHours*time.Hour,
	)
	if status == nil {
		return redis.Nil
	}

	err := status.Err()
	if err != nil {
		return err
	}

	return nil

}

func (c *playerCache) GetLoopMode(accountId uint) (models.PlayerLoopMode, error) {
	res := c.client.Get(context.Background(), playerSettingKey(accountId, playerSettingLoop))
	if res == nil {
		return models.LoopOffMode, nil
	}
	value, err := res.Result()
	if err != nil {
		return models.LoopOffMode, err
	}

	return models.PlayerLoopMode(value), nil
}

func playerPlayingPlaylistKey(accountId uint) string {
	return fmt.Sprintf("%splayer-playlist:%d", keyPrefix, accountId)
}

func (c *playerCache) SetCurrentPlayingPlaylistInQueue(accountId uint, playlistId uint) error {
	if playlistId == 0 {
		status := c.client.Del(context.Background(), playerPlayingPlaylistKey(accountId))
		if status == nil {
			return redis.Nil
		}
		if err := status.Err(); err != nil {
			return err
		}
	}

	status := c.client.Set(context.Background(), playerPlayingPlaylistKey(accountId), playlistId, playerQueueTtlHours*time.Hour)
	if status == nil {
		return redis.Nil
	}
	if err := status.Err(); err != nil {
		return err
	}

	return nil
}

func (c *playerCache) GetCurrentPlayingPlaylistInQueue(accountId uint) (uint, error) {
	res := c.client.Get(context.Background(), playerPlayingPlaylistKey(accountId))
	if res == nil {
		return 0, nil
	}
	value, err := res.Result()
	if err != nil {
		return 0, err
	}

	valueInt, err := strconv.Atoi(value)
	return uint(valueInt), err
}

func (c *playerCache) GetQueueLength(accountId uint) (uint, error) {
	res := c.client.LLen(context.Background(), playerQueueKey(accountId))
	if res == nil {
		return 0, redis.Nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return 0, err
	}

	return uint(queueLen), nil
}

func (c *playerCache) GetShuffledQueueLength(accountId uint) (uint, error) {
	res := c.client.LLen(context.Background(), playerShuffledQueueKey(accountId))
	if res == nil {
		return 0, redis.Nil
	}
	queueLen, err := res.Result()
	if err != nil {
		return 0, err
	}

	return uint(queueLen), nil
}

func (c *playerCache) GetSongIdAtIndexFromQueue(accountId uint, index int) (uint, error) {
	res := c.client.LIndex(context.Background(), playerQueueKey(accountId), int64(index))
	if res == nil {
		return 0, redis.Nil
	}

	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	songId, err := res.Uint64()

	return uint(songId), err
}

func (c *playerCache) GetSongIdAtIndexFromShuffledQueue(accountId uint, index int) (uint, error) {
	res := c.client.LIndex(context.Background(), playerShuffledQueueKey(accountId), int64(index))
	if res == nil {
		return 0, redis.Nil
	}

	err := res.Err()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	songId, err := res.Uint64()

	return uint(songId), err
}
