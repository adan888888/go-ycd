package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint  `gorm:"primarykey"`
	Uid       int64 `gorm:"colmun:uid;default:NULL"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime `gorm:"index"`
	Username  string       `gorm:"unique"`
	Password  string
}
