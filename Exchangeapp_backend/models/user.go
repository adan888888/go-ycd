package models

import (
	"time"
)

/*
*
references 通常用于指定外键引用的是源表的主键字段，是一种标准的外键关联方式。
AssociationForeignKey 允许你指定源表中用于关联的非主键字段，提供了更灵活的关联方式。
*/
type User struct {
	ID                uint  `gorm:"primaryKey;unique;autoIncrement"`
	Uid               int64 `gorm:"column:uid;unique;default:NULL"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Username          string `gorm:"unique"`
	Password          string
	Token             string
	TableYanchendao1s []TableYanchendao1 `gorm:"foreignKey:Uid;references:Uid"`    // 一个用户可以有多个表1数据
	TableYanchendao2s []TableYanchendao2 `gorm:"foreignKey:UserID;references:Uid"` // foreignKey:UserID 是外键
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
