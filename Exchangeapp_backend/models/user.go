package models

import (
	"time"
)

type User struct {
	ID        uint  `gorm:"primarykey"`
	Uid       int64 `gorm:"colmun:uid;default:NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time /*`gorm:"index"`*/
	Username  string    `gorm:"unique"`
	Password  string
	Token     string
}

type UserBody struct {
	Username string `json:"username" binding:"required" example:"admin"` // 用户名（必填）
	Password string `json:"password" binding:"required" example:"123"`   // 用户年龄（必填）
	Age      int    `json:"age" example:"25"`                            // 用户年龄（可选）
}

type Account struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"account name"`
}
type JSONResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type Order struct {
	OrderId  string `json:"order_id" example:"12313421"`
	OderName string `json:"order_name" example:"order name"`
}
