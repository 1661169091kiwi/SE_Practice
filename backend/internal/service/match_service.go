package service

import (
	"errors"

	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

type MatchService struct {
	matchRepo *repo.MatchRepo
	eventRepo *repo.EventRepo
	teamRepo  *repo.UserRepo // 复用用户仓库中的队伍查询方法
}

func NewMatchService() *MatchService {
	return &MatchService{
		matchRepo: repo.NewMatchRepo(),
		eventRepo: repo.NewEventRepo(),
		teamRepo:  repo.NewUserRepo(),
	}
}

// ListAllMatches 获取所有比赛
func (s *MatchService) ListAllMatches() ([]model.Match, error) {
	return s.matchRepo.ListAllMatches()
}

// ListMatchesByEvent 按赛事ID获取比赛列表
func (s *MatchService) ListMatchesByEvent(eventID int64) ([]model.Match, error) {
	return s.matchRepo.ListMatchesByEvent(eventID)
}

// CreateMatch 创建比赛
func (s *MatchService) CreateMatch(req *model.CreateMatchRequest) (*model.Match, error) {
	// 验证赛事是否存在
	event, err := s.eventRepo.GetEventByID(req.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}

	// 验证队伍是否存在
	teamA, err := s.teamRepo.GetTeamByID(req.TeamAID)
	if err != nil || teamA == nil {
		return nil, errors.New("team a not found")
	}
	teamB, err := s.teamRepo.GetTeamByID(req.TeamBID)
	if err != nil || teamB == nil {
		return nil, errors.New("team b not found")
	}

	// 构建比赛模型
	match := &model.Match{
		EventID: req.EventID,
		Name:    req.MatchName,
		Round:   req.Round,
		Time:    req.MatchTime,
		TeamAID: req.TeamAID,
		TeamBID: req.TeamBID,
		Status:  "not_started", // 默认未开始
	}

	// 保存到数据库
	id, err := s.matchRepo.CreateMatch(match)
	if err != nil {
		return nil, err
	}
	match.ID = id

	return match, nil
}

// GetMatchDetail 获取比赛详情
func (s *MatchService) GetMatchDetail(id int64) (*model.MatchDetailResponse, error) {
	match, err := s.matchRepo.GetMatchByID(id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, errors.New("match not found")
	}

	// 获取队伍信息
	teamA, _ := s.teamRepo.GetTeamByID(match.TeamAID)
	teamB, _ := s.teamRepo.GetTeamByID(match.TeamBID)

	return &model.MatchDetailResponse{
		MatchID:   match.ID,
		EventID:   match.EventID,
		MatchName: match.Name,
		Round:     match.Round,
		MatchTime: match.Time,
		TeamA: model.TeamBrief{
			ID:   teamA.ID,
			Name: teamA.TeamName,
		},
		TeamB: model.TeamBrief{
			ID:   teamB.ID,
			Name: teamB.TeamName,
		},
		ScoreA: match.ScoreA,
		ScoreB: match.ScoreB,
		Status: match.Status,
		Collectors: []model.UserBrief{
			// 实际项目中需要查询collectors表获取采集员信息
			{ID: match.Collector1ID},
			{ID: match.Collector2ID},
		},
	}, nil
}

// UpdateMatchScore 更新比赛分数
func (s *MatchService) UpdateMatchScore(id, scoreA, scoreB int64) error {
	if scoreA < 0 || scoreB < 0 {
		return errors.New("scores cannot be negative")
	}

	match, err := s.matchRepo.GetMatchByID(id)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found")
	}

	return s.matchRepo.UpdateMatchScore(id, int(scoreA), int(scoreB))
}
