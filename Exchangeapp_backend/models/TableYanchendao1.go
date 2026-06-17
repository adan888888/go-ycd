package models

import (
	"time"

	"gorm.io/gorm"
)

type TableYanchendao1 struct {
	ID           int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Benjin       float64        `gorm:"type:decimal(18,4);not null;comment:本金" json:"benjin"`
	YongJin      float64        `gorm:"column:yongjin;type:decimal(10,4);not null;comment:佣金" json:"yongjin"`
	Mean         float64        `gorm:"type:decimal(10,4);not null;comment:数学期望" json:"mean"`
	RestartIdx   int            `gorm:"column:restart_index;type:int;not null;default:0;comment:重起位置" json:"restart_index"`
	LiushuiIdx   int            `gorm:"column:liushui_index;type:int;not null;default:0;comment:流水的位置" json:"liushui_index"`
	ZhuangZhanBi int            `gorm:"column:zhuang_zhan_bi;type:int;default:50;comment:庄占比(0-100)" json:"zhuang_zhan_bi"`
	TempIndex    string         `gorm:"column:temp_index;type:varchar(255);comment:局部平衡位置" json:"temp_index"`
	CreatedAt    time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	Uid          int64          `gorm:"column:uid" json:"uid"`
	DeletedAt    gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
}

func (TableYanchendao1) TableName() string {
	return "table_yanchendao1"
}
