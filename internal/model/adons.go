package model

import "gorm.io/gorm"

type Adons struct {
	gorm.Model
	PostID      uint    `json:"-"`
	Post        *Post   `json:"-"`
	IsConnected bool    `gorm:"default: false" json:"is_connected"`
	Description string  `gorm:"type:text" json:"new_desc"`
	Codes       *string `gorm:"type:varchar(255)" json:"new_codes"`
	Service     string  `gorm:"type:varchar(50)" json:"-"`
}
