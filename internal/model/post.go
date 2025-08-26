package model

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID      uint   `json:"-"`
	User        *User  `json:"-"`
	Token       string `gorm:"type:varchar(255); uniqueIndex" json:"token"`
	Data        string `gorm:"type:text" json:"-"`
	IsConnected bool   `gorm:"default: false" json:"is_connected"`
	Title       string `gorm:"type:varchar(255)" json:"title"`
	Category    string `gorm:"type:varchar(255)" json:"-"`
	City        string `gorm:"type:varchar(255)" json:"-"`
	District    string `gorm:"type:varchar(255)" json:"-"`
}
