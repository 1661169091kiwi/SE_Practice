package model

type Sport struct {
	SportId   int    `gorm:"primaryKey;autoIncrement"`
	SportName string `gorm:"not null"`
}

type User struct {
	StudentId string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"default:'123456'"` // 默认密码
}
