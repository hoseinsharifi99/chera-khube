package model

import "gorm.io/gorm"

type AppLog struct {
	gorm.Model
	PostToken   string
	ServiceName string
}
