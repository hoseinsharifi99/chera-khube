package model

import "gorm.io/gorm"

type PricingLogic struct {
	gorm.Model
	ServiceID            uint   `json:"-"`
	Title                string `json:"title"`
	Price                uint   `gorm:"default:0" json:"price"`
	PriceWithOutDiscount uint   `json:"price_with_out_discount"`
	Color                string `gorm:"type:varchar(255)" json:"color"`
	Popular              bool   `json:"popular"`
	SpecialDiscount      bool   `json:"-"`
	Token                uint   `json:"token"`
	Tag                  string `json:"tag"`
	TagColor             string `json:"tag_color"`
}
