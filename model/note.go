package model

import "time"

type Note struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text"`
	ImageURL  string `gorm:"type:varchar(255)"`
	Summary   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
