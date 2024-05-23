package models

import "time"

type History struct {
	SongId    uint `gorm:"primaryKey"`
	ProfileId uint `gorm:"primaryKey"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (h History) GetId() uint {
	return h.ProfileId | h.SongId
}
