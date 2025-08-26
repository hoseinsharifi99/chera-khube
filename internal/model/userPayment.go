package model

import "gorm.io/gorm"

type UserPayment struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"user_id"`
	RefID  string `gorm:"not null" json:"ref_id"`
	Amount int    `gorm:"not null" json:"amount"`
	Status bool   `gorm:"default:false; not null" json:"status"`
}
