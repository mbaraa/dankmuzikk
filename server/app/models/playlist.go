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
	Owners    []*Account `gorm:"many2many:playlist_owners;"`
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

	err = tx.
		Model(new(PlaylistSongVoter)).
		Delete(&PlaylistSongVoter{
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

type PlaylistOwner struct {
	PlaylistId  uint `gorm:"primaryKey"`
	AccountId   uint `gorm:"primaryKey"`
	Permissions PlaylistPermissions

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PlaylistOwner) GetId() uint {
	return p.AccountId | p.PlaylistId
}

type PlaylistPermissions int8

const (
	VisitorPermission PlaylistPermissions = 1 << iota
	JoinerPermission
	OwnerPermission
	NonePermission PlaylistPermissions = 0
)

// PlaylistSongVoter ensures that an account had voted only once.
type PlaylistSongVoter struct {
	PlaylistId uint `gorm:"primaryKey"`
	SongId     uint `gorm:"primaryKey"`
	AccountId  uint `gorm:"primaryKey"`
	VoteUp     bool

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p PlaylistSongVoter) GetId() uint {
	return p.SongId | p.PlaylistId | p.AccountId
}
