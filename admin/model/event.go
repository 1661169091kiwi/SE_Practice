package model

import (
	"time"
)

type Event struct {
	EventId   int       `gorm:"primaryKey;autoIncrement"`
	EventName string    `gorm:"not null"`
	SportId   int       `gorm:"not null"`
	Location  string    `gorm:"column:location"` // 注意：你的DB设计里events表看似没有location，建议加上，或者存入description
	StartDate time.Time `gorm:"type:date"`
	Status    string    `gorm:"default:'upcoming'"`
	// 关联
	Sport Sport `gorm:"foreignKey:SportId"`
}

type Match struct {
	MatchId    int       `gorm:"primaryKey;autoIncrement"`
	EventId    int       `gorm:"not null"`
	MatchName  string    // 对应前端 "小组赛第一轮"
	MatchTime  time.Time `gorm:"not null"`
	TeamAId    int       `gorm:"not null"`
	TeamBId    int       `gorm:"not null"`
	ScoreTeamA int       `gorm:"default:0"`
	ScoreTeamB int       `gorm:"default:0"`
	Status     string
	Location   string `gorm:"-"` // 数据库matches表没设计Location字段，建议加上或存JSON
}
