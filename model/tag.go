package model

import "time"

// Tag 标签模型：存储所有唯一的标签
type Tag struct {
	ID        uint      `gorm:"primaryKey"`                   // 标签自增ID
	Name      string    `gorm:"size:50;uniqueIndex;not null"` // 标签名（唯一，避免重复，如“数学”）
	CreatedAt time.Time // 标签创建时间（AI生成或用户自定义时自动记录）
	UpdatedAt time.Time // 标签更新时间

	// 多对多关联：一个标签对应多个笔记（反向关联，可选，便于通过标签查笔记）
	Notes []Note `gorm:"many2many:note_tags;foreignKey:ID;joinForeignKey:TagID;References:ID;joinReferences:NoteID"`
}
