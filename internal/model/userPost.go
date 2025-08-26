package model

import "gorm.io/gorm"

type UserPosts struct {
	gorm.Model
	UserID    uint   `gorm:"not null"`
	PostToken string `gorm:"not null"`
}
