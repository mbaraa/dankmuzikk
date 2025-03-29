package models

import "time"

type Account struct {
	Id       uint   `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"unique;not null"`
	IsOAuth  bool
	Playlist []*Playlist `gorm:"many2many:playlist_owners"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a Account) GetId() uint {
	return a.Id
}
