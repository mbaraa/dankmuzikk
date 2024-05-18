package models

import (
	"time"

	"gorm.io/gorm"
)

type Playlist struct {
	Id uint `gorm:"primaryKey;autoIncrement"`

	PublicId   string `gorm:"unique;not null;index"`
	Title      string
	SongsCount int
	IsPublic   bool

	Songs     []*Song    `gorm:"many2many:playlist_songs;"`
	Owners    []*Profile `gorm:"many2many:playlist_owners;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Playlist) GetId() uint {
	return p.Id
}

func (p *Playlist) BeforeDelete(tx *gorm.DB) error {
	err := tx.
		Model(new(PlaylistOwner)).
		Delete(&PlaylistOwner{
			PlaylistId: p.Id,
		}, "playlist_id = ?", p.Id).
		Error
	if err != nil {
		return err
	}

	return tx.
		Model(new(PlaylistSong)).
		Delete(&PlaylistSong{
			PlaylistId: p.Id,
		}, "playlist_id = ?", p.Id).
		Error
}

type PlaylistSong struct {
	PlaylistId uint `gorm:"primaryKey"`
	SongId     uint `gorm:"primaryKey"`
	Votes      int
	PlayTimes  int

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

type PlaylistOwner struct {
	PlaylistId  uint `gorm:"primaryKey"`
	ProfileId   uint `gorm:"primaryKey"`
	Permissions PlaylistPermissions

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PlaylistOwner) GetId() uint {
	return p.ProfileId | p.PlaylistId
}

type PlaylistPermissions int8

const (
	ReadPermission PlaylistPermissions = 1 << iota
	WritePermission
	OwnerPermission
)
