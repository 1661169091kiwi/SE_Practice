package repository

import (
	"multimodal-notes/model"

	"gorm.io/gorm"
)

// CollectionRepository 收藏仓库接口
type CollectionRepository interface {
	//添加收藏userID-用户ID，noteID-笔记ID
	AddCollection(userID, noteID uint) error
	RemoveCollection(userID, noteID uint) error
	// 查询用户收藏列表：按收藏时间倒序，预加载笔记标签
	GetUserCollections(userID uint, page, pageSize int) ([]model.Collection, int64, error)
	//检查是否已收藏：返回 true 表示已收藏
	IsCollected(userID, noteID uint) (bool, error)
}

type collectionRepository struct {
	db *gorm.DB
}

func NewCollectionRepository(db *gorm.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

func (r *collectionRepository) AddCollection(userID, noteID uint) error {
	collection := &model.Collection{
		UserID: userID,
		NoteID: noteID,
	}
	// FirstOrCreate避免重复收藏
	result := r.db.FirstOrCreate(collection, model.Collection{UserID: userID, NoteID: noteID})
	return result.Error
}

func (r *collectionRepository) RemoveCollection(userID, noteID uint) error {
	result := r.db.Delete(&model.Collection{}, "user_id = ? AND note_id = ?", userID, noteID)
	return result.Error
}

func (r *collectionRepository) GetUserCollections(userID uint, page, pageSize int) ([]model.Collection, int64, error) {
	var (
		collections []model.Collection
		total       int64
	)
	// 先查总数（用于分页）
	err := r.db.Model(&model.Collection{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 分页查询：预加载笔记详情、笔记的标签、笔记的作者（只查用户名，隐藏密码）
	offset := (page - 1) * pageSize
	err = r.db.Preload("Note.Tags"). // 预加载笔记的标签
						Preload("Note.User", "id, username"). // 预加载笔记作者（只查id和用户名）
						Where("user_id = ?", userID).
						Order("created_at DESC"). // 按收藏时间倒序（最新收藏在前）
						Limit(pageSize).
						Offset(offset).
						Find(&collections).Error
	if err != nil {
		return nil, 0, err
	}
	return collections, total, err
}

// 4. 检查是否已收藏：用于前端判断“显示收藏按钮”还是“取消收藏按钮”
func (r *collectionRepository) IsCollected(userID, noteID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Collection{}).Where("user_id = ? AND note_id = ?", userID, noteID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil //  count>0 表示已收藏
}
