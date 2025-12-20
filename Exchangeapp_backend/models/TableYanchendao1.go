package models

import (
	"time"

	"gorm.io/gorm"
)

type TableYanchendao1 struct {
	ID               int            `gorm:"primaryKey;autoIncrement" json:"id"` //json:"-"
	ColumnBenjin     string         `gorm:"type:varchar(255);not null;comment:'本金'" json:"column_benjin"`
	ColumnYongJin    string         `gorm:"column:column_yongjin;type:varchar(255);not null;comment:'俑金'" json:"column_yongJin"`
	ColumnMean       string         `gorm:"type:varchar(255);not null;comment:'数学期望'" json:"column_mean"`
	ColumnRestartIdx string         `gorm:"column:column_restart_index;type:varchar(255);not null;comment:'重起位置'" json:"column_restart_index"`
	ColumnLiushuiIdx string         `gorm:"column:column_liushui_index;type:varchar(255);not null;comment:'流水的位置'" json:"column_liushui_index"`
	ColumnZhuangZhanBi int          `gorm:"column:column_zhuang_zhan_bi;type:int;default:50;comment:'庄占比(0-100)'" json:"column_zhuang_zhan_bi"`
	TempIndex        string         `gorm:"column:temp_index;type:varchar(255)" json:"temp_index"`
	CreatedAt        time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:'创建时间'" json:"created_at"`
	Uid              int64          `gorm:"column:uid"`
	DeletedAt        gorm.DeletedAt `gorm:"index;comment:'删除时间'" json:"deleted_at"`
}

// TableName 显式指定数据库表名
func (TableYanchendao1) TableName() string {
	return "table_yanchendao1"
}
