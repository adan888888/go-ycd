package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title   string `binding:"required" gorm:"size:10"` //类型是varchar  gorm:"size:255" 长度是255个字符 gorm:"size:10" 类型是varchar 长度是10
	Content string `binding:"required"`                //没有指定的默认是longtext类型
	Preview string `binding:"required"`
}
