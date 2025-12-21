package service

import (
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

type RatingService struct {
	repo *repo.RatingRepo
}

func NewRatingService() *RatingService {
	return &RatingService{repo: repo.NewRatingRepo()}
}

func (s *RatingService) Upsert(matchID int64, studentID string, raterID string, score float64, comment string) error {
	return s.repo.Upsert(matchID, studentID, raterID, score, comment)
}

func (s *RatingService) ListByMatch(matchID int64) ([]model.PlayerRating, error) {
	return s.repo.ListByMatch(matchID)
}
