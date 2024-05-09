package models

import (
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	AccountId uint
	Account   Account

	PublicId   string `gorm:"unique;not null;index"`
	Title      string
	SongsCount int
	Songs      []*Song `gorm:"many2many:playlist_songs;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (p Playlist) GetId() uint {
	return p.Id
}

type PlaylistSong struct {
	PlaylistId uint `gorm:"primaryKey"`
	SongId     uint `gorm:"primaryKey"`
	Votes      int

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PlaylistSong) GetId() uint {
	return p.SongId | p.PlaylistId
}

func (p *PlaylistSong) BeforeCreate(tx *gorm.DB) error {
	var playlist Playlist
	err := tx.
		Select("songs_count").
		Where("id = ?", p.PlaylistId).
		First(&playlist).
		Error
	if err != nil {
		return err
	}

	return tx.
		Model(&playlist).
		Where("id = ?", p.PlaylistId).
		Update("songs_count", playlist.SongsCount+1).
		Error
}

func (p *PlaylistSong) BeforeDelete(tx *gorm.DB) error {
	var playlist Playlist
	err := tx.
		Select("songs_count").
		Where("id = ?", p.PlaylistId).
		First(&playlist).
		Error
	if err != nil {
		return err
	}

	return tx.
		Model(&playlist).
		Where("id = ?", p.PlaylistId).
		Update("songs_count", playlist.SongsCount-1).
		Error
}
