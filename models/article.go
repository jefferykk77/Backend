package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `gorm:"required"`
	Content string `gorm:"required"`
	Preview string `gorm:"required"`
	Likes   int    `gorm:"default:0"` // 点赞数，默认为0
}
