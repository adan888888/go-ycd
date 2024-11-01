package models

import (
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	Username  string    `gorm:"unique"`
	Password  string
}
