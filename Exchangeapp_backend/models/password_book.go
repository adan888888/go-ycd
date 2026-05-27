package models

import (
	"time"
)

// PasswordItem 密码项模型
// 表名: password_items
// 用途: 存储密码信息
type PasswordItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:主键ID" json:"id"`
	Uid       int64     `gorm:"column:uid;index;not null;comment:所属用户uid，关联 users.uid" json:"uid"`
	Title     string    `gorm:"size:255;not null;index:idx_title;comment:密码项标题" json:"title"`     // 标题
	Username  string    `gorm:"size:255;not null;index:idx_username;comment:用户名" json:"username"` // 用户名
	Password  string    `gorm:"size:500;not null;comment:密码，建议加密存储" json:"password"`              // 密码
	Website   string    `gorm:"size:500;comment:网站地址" json:"website"`                             // 网站
	Notes     string    `gorm:"type:text;comment:备注信息" json:"notes"`                              // 备注
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// PasswordItemRequest 创建/更新密码项请求
// 用途: 用于接收客户端创建或更新密码项的请求数据
// 验证: 使用 gin 的 binding 标签进行参数验证
type PasswordItemRequest struct {
	Title    string `json:"title" binding:"required" example:"百度云盘"`                   // 标题（必填）
	Username string `json:"username" binding:"required" example:"shitlaoge@gmail.com"` // 用户名（必填）
	Password string `json:"password" binding:"required" example:"password123"`         // 密码（必填）
	Website  string `json:"website" example:"http://example.com"`                      // 网站（可选）
	Notes    string `json:"notes" example:"备注信息"`                                      // 备注（可选）
}

// PasswordItemResponse 密码项响应
// 用途: 用于返回给客户端的密码项数据
// 注意: 包含完整的密码信息，客户端需要妥善处理
type PasswordItemResponse struct {
	ID        uint      `json:"id"`         // 密码项ID
	Uid       int64     `json:"uid"`        // 所属用户
	Title     string    `json:"title"`      // 标题
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码
	Website   string    `json:"website"`    // 网站
	Notes     string    `json:"notes"`      // 备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// PasswordSearchRequest 搜索请求
// 用途: 用于接收密码搜索的请求参数
// 搜索: 支持按标题、用户名、网站、备注进行模糊搜索
type PasswordSearchRequest struct {
	Keyword string `json:"keyword" example:"小幺鸡"` // 搜索关键词
}
