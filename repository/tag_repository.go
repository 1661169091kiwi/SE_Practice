package repository

import (
	"multimodal-notes/model"

	"gorm.io/gorm"
)

// TagRepository 标签仓库接口
type TagRepository interface {
	GetOrCreateTags(tagNames []string) ([]model.Tag, error) // 查重并创建标签
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// GetOrCreateTags 核心逻辑：对每个标签，存在则查询，不存在则创建
func (r *tagRepository) GetOrCreateTags(tagNames []string) ([]model.Tag, error) {
	var tagList []model.Tag
	for _, name := range tagNames {
		var tag model.Tag
		result := r.db.FirstOrCreate(&tag, model.Tag{Name: name})
		if result.Error != nil {
			return nil, result.Error
		}
		tagList = append(tagList, tag)
	}
	return tagList, nil
}
