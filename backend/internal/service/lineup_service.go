package service

import (
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

type LineupService struct {
	repo *repo.LineupRepo
}

func NewLineupService() *LineupService {
	return &LineupService{repo: repo.NewLineupRepo()}
}

func (s *LineupService) Upsert(items []model.MatchLineup) error {
	for _, it := range items {
		if err := s.repo.Upsert(it); err != nil {
			return err
		}
	}
	return nil
}

func (s *LineupService) ListByMatch(matchID int64) ([]model.MatchLineup, error) {
	return s.repo.ListByMatch(matchID)
}
