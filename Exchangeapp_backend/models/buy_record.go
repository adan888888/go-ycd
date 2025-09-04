package models

import (
	"time"
)

// BuyRecord 买币记录表
type BuyRecord struct {
	ID        uint      `gorm:"primaryKey;unique;autoIncrement;comment:主键ID" json:"id"`
	Currency  string    `gorm:"column:currency;size:20;not null;index;comment:币种标识，如 BTC、ETH、USDT等" json:"currency"`
	BuyPrice  float64   `gorm:"column:buy_price;type:decimal(20,8);not null;comment:买入价格" json:"buy_price"`
	BuyAmount float64   `gorm:"column:buy_amount;type:decimal(20,8);not null;comment:买入了多少钱" json:"buy_amount"`
	BuyTime   time.Time `gorm:"column:buy_time;not null;comment:买入时间" json:"buy_time"`

	CreatedAt time.Time `gorm:"comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:更新时间" json:"updated_at"`
	DeletedAt time.Time `gorm:"default:NULL;comment:删除时间" json:"deleted_at"`
}

// TableName 指定表名
func (BuyRecord) TableName() string {
	return "buy_records"
}



// BuyRecordResponse 买币记录的响应结构
type BuyRecordResponse struct {
	ID        uint      `json:"id"`
	Currency  string    `json:"currency"`
	BuyPrice  float64   `json:"buy_price"`
	BuyAmount float64   `json:"buy_amount"`
	BuyTime   time.Time `json:"buy_time"`
	CreatedAt time.Time `json:"created_at"`
}
