package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	PhoneNumber       string    `gorm:"type:varchar(20);unique;not null" json:"-"`
	Balance           int       `gorm:"default:0" json:"balance"`
	PostToken         string    `gorm:"type:varchar(255)" json:"post_token"`
	AccessToken       string    `gorm:"unique" json:"-"`
	RefreshToken      string    `gorm:"type:varchar(255)" json:"-"`
	ExpiresAt         int64     `gorm:"default:0" json:""`
	Password          string    `gorm:"type:varchar(255);not null" json:"-"`
	JwtToken          string    `gorm:"-" json:"token"`
	IsUse             bool      `gorm:"default:false" json:"is_use"`
	IsFirstAgahi      bool      `gorm:"default:false" json:"is_first_agahi"`
	IsFirstProfileBuy bool      `gorm:"default:false" json:"is_first_profile_buy"`
	JwtExpireAt       time.Time `gorm:"-" json:"expire_at"`
	ExpirationTime    time.Time `json:"expiration_time"`
	ServiceId         int       `json:"-"`
}
