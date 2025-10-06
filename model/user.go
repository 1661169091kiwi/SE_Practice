package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `gorm:"not null" json:"password"`
	Email     string `gorm:"unique;not null" json:"email"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
