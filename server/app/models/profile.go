package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	Id        uint `gorm:"primaryKey;autoIncrement"`
	AccountId uint
	Account   Account
	Name      string `gorm:"not null"`
	PfpLink   string
	Username  string `gorm:"unique;not null;index"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p Profile) GetId() uint {
	return p.Id
}

func (p *Profile) AfterDelete(tx *gorm.DB) error {
	return tx.
		Model(new(Account)).
		Delete(&p.Account, "id = ?", p.AccountId).
		Error
}
