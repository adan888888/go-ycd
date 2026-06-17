package models

import (
	"time"

	"gorm.io/gorm"
)

type TableYanchendao2 struct {
	ID                  int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Xiazhujine          float64        `gorm:"type:decimal(18,4);not null;comment:下注金额" json:"xiazhujine"`
	Shuyingzhi          float64        `gorm:"type:decimal(18,4);not null;comment:输赢值" json:"shuyingzhi"`
	ShuyingzhiXiaoshu   *float64       `gorm:"column:shuyingzhi_xiaoshu;type:decimal(18,4);comment:消数后输赢值" json:"shuyingzhi_xiaoshu"`
	Shengfulu           string         `gorm:"type:varchar(10);not null;comment:胜负路" json:"shengfulu"`
	ZX                  string         `gorm:"column:zx;type:varchar(10);not null;comment:庄闲" json:"zx"`
	Remark              string         `gorm:"type:varchar(255);comment:输赢标记备注" json:"remark"`
	CurrentJin          float64        `gorm:"column:current_jin;type:decimal(18,4);not null;comment:当前金额" json:"current_jin"`
	RestartStatSnapshot string         `gorm:"column:restart_stat_snapshot;type:varchar(255);default:'';comment:重启统计快照" json:"restartStatSnapshot"`
	ColumnRefresh       bool           `gorm:"default:false;comment:刷新标记" json:"-"`
	CreatedAt           time.Time      `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`
	UserID              int64          `gorm:"column:user_id" json:"user_id"`
	DeletedAt           gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
}

func (TableYanchendao2) TableName() string {
	return "table_yanchendao2"
}
