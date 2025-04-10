package models

import (
	"time"
)

type TableYanchendao2 struct {
	ID                int       `gorm:"primaryKey;autoIncrement" json:"id"` //binding:"required"
	ColumnXiazhujine  string    `gorm:"type:varchar(255);not null;comment:'下注的金额'" json:"column_xiazhujine"`
	ColmunShuyingzhi  string    `gorm:"type:varchar(255);not null;comment:'输赢值'" json:"colmun_shuyingzhi"`
	ColmunShuyingzhiD string    `gorm:"type:varchar(255);not null;comment:'消数后的输赢值'" json:"colmun_shuyingzhi_d"`
	ColmunShengfulu   string    `gorm:"type:varchar(10);not null;comment:'胜负路（输赢标记）'" json:"colmun_shengfulu"`
	ColmunZX          string    `gorm:"type:varchar(10);not null;comment:'开出的是庄还是闲'" json:"colmun_zx"`
	ColmunRemark      string    `gorm:"type:varchar(255);comment:'输赢标记备注'" json:"colmun_remark"`
	ColumnCurrentJin  string    `gorm:"type:varchar(255);not null;comment:'当前的钱'" json:"column_current_jin"`
	ColumnRefresh     bool      `gorm:"default:false;comment:'用来刷新用'" json:"-"`
	CreatedAt         time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"-"`
	UserID            int64     `gorm:"column:user_id"`
}

// TableName 显式指定数据库表名
func (TableYanchendao2) TableName() string {
	return "table_yanchendao2"
}
