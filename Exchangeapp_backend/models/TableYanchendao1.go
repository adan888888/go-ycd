package models

import (
	"time"
)

type TableYanchendao1 struct {
	ID               int       `gorm:"primaryKey;autoIncrement" json:"-"`
	ColumnBenjin     string    `gorm:"type:varchar(255);not null;comment:'本金'" json:"column_benjin"`
	ColumnYongJin    string    `gorm:"column:column_yongjin;type:varchar(255);not null;comment:'俑金'" json:"column_yongJin"`
	ColumnMean       string    `gorm:"type:varchar(255);not null;comment:'数学期望'" json:"column_mean"`
	ColumnRestartIdx string    `gorm:"column:column_restart_index;type:varchar(255);not null;comment:'重起位置'" json:"column_restart_index"`
	ColumnLiushuiIdx string    `gorm:"column:column_liushui_index;type:varchar(255);not null;comment:'流水的位置'" json:"column_liushui_index"`
	CreatedAt        time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"-"`
	Uid              int64     `gorm:"column:uid"`
}

// TableName 显式指定数据库表名
func (TableYanchendao1) TableName() string {
	return "table_yanchendao1"
}
