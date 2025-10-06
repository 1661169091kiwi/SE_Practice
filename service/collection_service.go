package service

import (
	"errors"
	"multimodal-notes/model"
	"multimodal-notes/repository"
)

// CollectionService 收藏服务接口
type CollectionService interface {
	AddCollection(userID, noteID uint) error
	RemoveCollection(userID, noteID uint) error
	// 收藏列表分页的
	GetUserCollections(userID uint, page, pageSize int) ([]model.Collection, int64, error)
	IsCollected(userID, noteID uint) (bool, error)
}

type collectionService struct {
	colRepo  repository.CollectionRepository
	noteRepo repository.NoteRepository // 依赖笔记仓库：检查笔记是否存在
}

func NewCollectionService(colRepo repository.CollectionRepository, noteRepo repository.NoteRepository) CollectionService {
	return &collectionService{
		colRepo:  colRepo,
		noteRepo: noteRepo,
	}
}

func (s *collectionService) AddCollection(userID, noteID uint) error {

	notes, err := s.noteRepo.FindNoteByID(noteID)
	if err != nil || len(notes) == 0 {
		return errors.New("笔记不存在")
	}
	return s.colRepo.AddCollection(userID, noteID)
}

func (s *collectionService) RemoveCollection(userID, noteID uint) error {
	return s.colRepo.RemoveCollection(userID, noteID)
}

func (s *collectionService) GetUserCollections(userID uint, page, pageSize int) ([]model.Collection, int64, error) {
	return s.colRepo.GetUserCollections(userID, page, pageSize)
}

func (s *collectionService) IsCollected(userID, noteID uint) (bool, error) {
	return s.colRepo.IsCollected(userID, noteID)
}
