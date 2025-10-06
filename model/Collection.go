package model

import "time"

// Collection 收藏表：用户与笔记的多对多关联中间表
type Collection struct {
	UserID    uint      `gorm:"primaryKey;not null" json:"userId"` // 关联用户ID（复合主键之一）
	NoteID    uint      `gorm:"primaryKey;not null" json:"noteId"` // 关联笔记ID（复合主键之一）
	CreatedAt time.Time `json:"createdAt"`                         // 收藏时间（自动填充）

	// 关联用户和笔记（便于查询时预加载）
	User *User `gorm:"foreignKey:UserID" json:"-"`    // 隐藏前端不需要的关联字段
	Note *Note `gorm:"foreignKey:NoteID" json:"note"` // 收藏的笔记详情
}

func (Collection) TableName() string {
	return "collections"
}
