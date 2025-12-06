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

// GetMatchData 获取赛事数据（积分榜/赛程等）
func (s *MatchService) GetMatchData(sportType int, matchTime, dataType string) (*model.MatchDataResponse, error) {
	var list []model.ResultItem
	var err error

	switch dataType {
	case "积分榜":
		list, err = s.matchRepo.GetStandings(sportType, matchTime)
	case "赛程":
		list, err = s.matchRepo.GetSchedule(sportType, matchTime)
	case "历史":
		list, err = s.matchRepo.GetHistory(sportType, matchTime)
	case "球员榜", "球队榜":
		// TODO: 暂时返回空数据，后续实现
		list = []model.ResultItem{}
	default:
		return nil, errors.New("invalid dataType")
	}

	if err != nil {
		return nil, err
	}

	// 如果没有数据，返回空数组而不是 nil (符合前端预期)
	if list == nil {
		list = []model.ResultItem{}
	}

	return &model.MatchDataResponse{
		DataList: list,
	}, nil
}

// GetUserSubscribedMatches 获取用户订阅赛事的比赛列表
func (s *MatchService) GetUserSubscribedMatches(studentID string) (*model.SubscribedMatchResponse, error) {
	if studentID == "" {
		return nil, errors.New("studentID is required")
	}

	matches, err := s.matchRepo.GetMatchesBySubscription(studentID)
	if err != nil {
		return nil, err
	}

	if matches == nil {
		matches = []model.SubscribedMatchItem{}
	}

	return &model.SubscribedMatchResponse{
		SubscribedMatches: matches,
	}, nil
}

// SubscribeMatch 处理订阅/取消订阅逻辑
func (s *MatchService) SubscribeMatch(req *model.SubscribeRequest) (string, error) {
	// 1. 根据 MatchID 找到 EventID
	eventID, err := s.matchRepo.GetEventIDByMatchID(req.MatchID)
	if err != nil {
		return "", err
	}
	if eventID == 0 {
		return "", errors.New("match not found")
	}

	// 2. 检查当前订阅状态
	isSubscribed, err := s.matchRepo.CheckSubscription(req.StudentID, eventID)
	if err != nil {
		return "", err
	}

	// 3. 根据 operateType 执行逻辑
	if req.OperateType == 0 { // 新增订阅
		if isSubscribed {
			return "已订阅该比赛，无需重复操作", nil // Return message, no error
		}
		err = s.matchRepo.SubscribeToEvent(req.StudentID, eventID)
		if err != nil {
			return "", err
		}
		return "订阅成功", nil
	} else if req.OperateType == 1 { // 取消订阅
		if !isSubscribed {
			return "未订阅该比赛，无法取消", nil
		}
		err = s.matchRepo.UnsubscribeFromEvent(req.StudentID, eventID)
		if err != nil {
			return "", err
		}
		return "取消订阅成功", nil
	}

	return "", errors.New("invalid operateType")
}
