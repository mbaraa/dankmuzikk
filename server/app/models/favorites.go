package models

import "time"

type FavoriteSong struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	SongId    uint
	AccountId uint
	CreatedAt time.Time
}
