package models

import (
	"time"
)

type Song struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	YtId            string
	PublicId        string `gorm:"unique;not null;index"`
	Title           string
	Artist          string
	ThumbnailUrl    string
	Duration        string
	RealDuration    time.Duration
	Playlists       []*Playlist `gorm:"many2many:playlist_songs;"`
	FullyDownloaded bool
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Playlist only fields

	AddedAt   string `gorm:"-"`
	Votes     int    `gorm:"-"`
	PlayTimes int    `gorm:"-"`

	// View only fields

	Favorite bool `gorm:"-"`
}

func (s Song) GetId() uint {
	return s.Id
}
