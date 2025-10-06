package repository

import (
	"gorm.io/gorm"
	"multimodal-notes/model"
)

// NoteRepository 接口：修改方法名，支持关联标签
type NoteRepository interface {
	CreateNoteWithTags(note *model.Note) (uint, error) // 新增：保存笔记+关联标签
	FindNotesByUserID(userID uint) ([]model.Note, error)
	FindNoteByID(noteID uint) ([]model.Note, error) // 新增：按ID查笔记

}

type noteRepository struct {
	db *gorm.DB
}

func (r *noteRepository) FindNoteByID(noteID uint) ([]model.Note, error) {
	var notes []model.Note
	result := r.db.Where("id = ?", noteID).Find(&notes)
	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil
}
func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}

// CreateNoteWithTags 保存笔记并关联标签（GORM 会自动处理中间表）
func (r *noteRepository) CreateNoteWithTags(note *model.Note) (uint, error) {
	// 关键：使用 Create 保存笔记时，GORM 会自动将 Tags 关联到中间表 note_tags
	result := r.db.Create(note)
	if result.Error != nil {
		return 0, result.Error
	}
	return note.ID, nil
}

func (r *noteRepository) FindNotesByUserID(userID uint) ([]model.Note, error) {
	var notes []model.Note
	// Preload("Tags")：预加载标签，避免 N+1 查询问题
	result := r.db.Preload("Tags").Where("user_id = ?", userID).Order("created_at DESC").Find(&notes)
	if result.Error != nil {
		return nil, result.Error
	}
	return notes, nil
}
