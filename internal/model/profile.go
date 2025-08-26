package model

import (
	"gorm.io/gorm"
	"time"
)

type Profile struct {
	gorm.Model
	UserID            uint      `json:"-"`
	User              *User     `json:"-"`
	TelegramLink      string    `gorm:"type:varchar(50)" json:"telegram_link"`
	WhatsAppLink      string    `gorm:"type:varchar(50)" json:"whatsapp_link"`
	Instagram         string    `gorm:"type:varchar(50)" json:"instagram"`
	Eitaa             string    `gorm:"type:varchar(50)" json:"eitaa"`
	Robika            string    `gorm:"type:varchar(50)" json:"robika"`
	PhoneNumber       string    `gorm:"type:varchar(20)" json:"phone_number"`
	CustomPhoneNumber string    `gorm:"type:varchar(20)" json:"custom_phone_number"`
	Address           string    `gorm:"type:text" json:"address"`
	Lat               float64   `json:"lat"`
	Long              float64   `json:"long"`
	PostToken         string    `gorm:"type:varchar(100)" json:"post_token"`
	UserName          string    `gorm:"type:varchar(50)" json:"username"`
	IsConnected       bool      `gorm:"type:boolean" json:"is_connected"`
	DivarID           string    `gorm:"type:varchar(50)" json:"-"`
	ConnectedUntil    time.Time `json:"connected_until"`
}
