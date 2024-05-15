package models

import (
	"time"
)

type Song struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	YtId         string `gorm:"unique;not null;index"`
	Title        string
	Artist       string
	ThumbnailUrl string
	Duration     string
	Playlists    []*Playlist `gorm:"many2many:playlist_songs;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (s Song) GetId() uint {
	return s.Id
}
