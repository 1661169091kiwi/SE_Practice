package model

import "time"

type Note struct {
	ID        uint      `gorm:"primaryKey"`        // 笔记自增ID
	Content   string    `gorm:"type:text"`         // 原始笔记内容
	ImageURL  string    `gorm:"type:varchar(255)"` // 图片Base64
	Summary   string    `gorm:"type:text"`         // 总结内容
	UserID    uint      `gorm:"not null"`          // 关联用户ID（归属）
	CreatedAt time.Time // 笔记创建时间
	UpdatedAt time.Time // 笔记更新时间

	// 新增：多对多关联——一个笔记对应多个标签
	// gorm:"many2many:note_tags..." 表示自动生成关联表 note_tags
	Tags []Tag `gorm:"many2many:note_tags;foreignKey:ID;joinForeignKey:NoteID;References:ID;joinReferences:TagID"`

	User        *User        `gorm:"foreignKey:UserID" json:"author"` // 笔记作者（前端显示）
	Collections []Collection `gorm:"foreignKey:NoteID" json:"-"`      // 收藏该笔记的用户（新增，隐藏）
}
