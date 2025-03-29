package models

import "time"

type History struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	SongId    uint
	AccountId uint
	CreatedAt time.Time
}

func (h History) GetId() uint {
	return h.Id
}
