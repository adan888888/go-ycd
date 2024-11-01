package models

import (
	"time"
)

type Article struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	Title     string    `binding:"required" gorm:"size:10"` //类型是varchar  gorm:"size:255" 长度是255个字符 gorm:"size:10" 类型是varchar 长度是10
	Content   string    `binding:"required"`                //没有指定的默认是longtext类型
	Preview   string    `binding:"required"`
}
