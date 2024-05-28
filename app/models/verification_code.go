package models

import "time"

type EmailVerificationCode struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	AccountId uint
	Account   Account
	Code      string `gorm:"not null"`
	CreatedAt time.Time
}

func (e EmailVerificationCode) GetId() uint {
	return e.Id
}
