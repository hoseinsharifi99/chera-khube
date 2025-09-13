package model

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	Code        string `gorm:"type:varchar(255)" json:"token"`
	Description string `gorm:"type:varchar(255)" json:"description"`
	Category    string `gorm:"type:varchar(255)" json:"category"`
	ServiceName string `gorm:"type:varchar(255)" json:"-"`
}
