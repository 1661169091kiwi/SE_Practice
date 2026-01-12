package service

import (
	"errors"

	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidStatus = errors.New("invalid event status")
)

type EventService struct {
	eventRepo *repo.EventRepo
	teamRepo  *repo.UserRepo
}

func NewEventService() *EventService {
	return &EventService{
		eventRepo: repo.NewEventRepo(),
		teamRepo:  repo.NewUserRepo(),
	}
}

// CreateEvent 创建赛事
func (s *EventService) CreateEvent(req *model.CreateEventRequest) (*model.EventResponse, error) {
	// 验证参数
	if req.EventName == "" || req.SportID <= 0 {
		return nil, errors.New("event name and sport id are required")
	}
	if req.StartDate.After(req.EndDate) {
		return nil, errors.New("start date must be before end date")
	}

	// 设置默认 Format 或验证 Format
	if req.Format == "" {
		req.Format = "points" // 默认积分制
	}
	if req.Format != "points" && req.Format != "knockout" && req.Format != "group_knockout" {
		return nil, errors.New("invalid format type: must be 'points' or 'knockout' or 'group_knockout'")
	}
	if len(req.Teams) == 0 {
		return nil, errors.New("teams are required")
	}

	// 构建赛事模型
	event := &model.Event{
		Name:      req.EventName,
		SportID:   req.SportID,
		Season:    req.Season,
		Round:     req.Round,
		Format:    req.Format,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "upcoming", // 默认未开始
	}

	// 保存到数据库
	id, err := s.eventRepo.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	if err := s.replaceEventTeams(id, req.SportID, req.Teams); err != nil {
		return nil, err
	}
	if err := s.maybeSetupKnockoutBracket(id, req.Format, req.Teams, req.KnockoutSchedule); err != nil {
		_ = s.eventRepo.DeleteEvent(id)
		return nil, err
	}

	// 构建响应
	return &model.EventResponse{
		EventID:   id,
		EventName: event.Name,
		SportID:   event.SportID,
		Status:    event.Status,
		StartDate: event.StartDate,
		EndDate:   event.EndDate,
	}, nil
}

// GetEventDetail 获取赛事详情
func (s *EventService) GetEventDetail(id int64) (*model.Event, error) {
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}
	return event, nil
}

// ListEvents 列出赛事
func (s *EventService) ListEvents(status string) ([]model.Event, error) {
	// 验证状态合法性
	validStatus := map[string]bool{
		"":         true, // 空表示查询所有
		"upcoming": true,
		"ongoing":  true,
		"finished": true,
	}
	if !validStatus[status] {
		return nil, ErrInvalidStatus
	}

	return s.eventRepo.ListEvents(status)
}

// DeleteEvent 删除赛事
func (s *EventService) DeleteEvent(id int64) error {
	return s.eventRepo.DeleteEvent(id)
}

// UpdateEventStatus 更新赛事状态
func (s *EventService) UpdateEventStatus(id int64, status string) error {
	if !isValidEventStatus(status) {
		return ErrInvalidStatus
	}

	// 检查赛事是否存在
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotFound
	}

	return s.eventRepo.UpdateEventStatus(id, status)
}

// UpdateEvent 更新赛事信息
func (s *EventService) UpdateEvent(id int64, req *model.CreateEventRequest) error {
	// 检查赛事是否存在
	event, err := s.eventRepo.GetEventByID(id)
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotFound
	}

	// 验证参数
	if req.EventName == "" || req.SportID <= 0 {
		return errors.New("event name and sport id are required")
	}
	if req.StartDate.After(req.EndDate) {
		return errors.New("start date must be before end date")
	}
	if req.Format == "" {
		req.Format = event.Format
	}
	if req.Format != "points" && req.Format != "knockout" && req.Format != "group_knockout" {
		return errors.New("invalid format type: must be 'points' or 'knockout' or 'group_knockout'")
	}

	// 更新字段
	event.Name = req.EventName
	event.SportID = req.SportID
	event.Season = req.Season
	event.Round = req.Round
	event.Format = req.Format
	event.StartDate = req.StartDate
	event.EndDate = req.EndDate

	if err := s.eventRepo.UpdateEvent(event); err != nil {
		return err
	}
	if req.Teams != nil {
		if err := s.replaceEventTeams(id, req.SportID, req.Teams); err != nil {
			return err
		}
		if err := s.maybeSetupKnockoutBracket(id, req.Format, req.Teams, req.KnockoutSchedule); err != nil {
			return err
		}
	}
	return nil
}

func isValidEventStatus(status string) bool {
	return status == "upcoming" || status == "ongoing" || status == "finished"
}

func (s *EventService) ListEventTeams(eventID int64) ([]model.EventTeam, error) {
	return s.eventRepo.ListEventTeams(eventID)
}

func (s *EventService) ReplaceEventTeams(eventID int64, teams []model.EventTeamInput) error {
	ev, err := s.eventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}
	if ev == nil {
		return ErrEventNotFound
	}
	if err := s.replaceEventTeams(eventID, ev.SportID, teams); err != nil {
		return err
	}
	return s.maybeSetupKnockoutBracket(eventID, ev.Format, teams, nil)
}

func (s *EventService) replaceEventTeams(eventID, sportID int64, teams []model.EventTeamInput) error {
	if len(teams) == 0 {
		return errors.New("teams are required")
	}
	seen := map[int64]struct{}{}
	for _, t := range teams {
		if t.TeamID <= 0 {
			return errors.New("invalid team_id")
		}
		if _, ok := seen[t.TeamID]; ok {
			return errors.New("duplicate team_id")
		}
		seen[t.TeamID] = struct{}{}
		team, err := s.teamRepo.GetTeamByID(t.TeamID)
		if err != nil {
			return err
		}
		if team == nil {
			return errors.New("team not found")
		}
		if !team.IsApproved {
			return errors.New("team is not approved")
		}
		if team.SportID != sportID {
			return errors.New("team sport_id mismatch")
		}
	}
	return s.eventRepo.ReplaceEventTeams(eventID, teams)
}

func (s *EventService) maybeSetupKnockoutBracket(eventID int64, format string, teams []model.EventTeamInput, schedule map[string]string) error {
	if format != "knockout" && format != "group_knockout" {
		return nil
	}
	if len(teams) == 0 {
		return errors.New("teams are required")
	}

	anySlot := false
	allSlot := true
	for _, t := range teams {
		if t.Slot > 0 {
			anySlot = true
		} else {
			allSlot = false
		}
	}

	if format == "knockout" {
		if !allSlot {
			return errors.New("knockout slot required")
		}
		return s.eventRepo.SetupKnockoutBracket(eventID, schedule)
	}

	if !anySlot {
		return nil
	}
	return s.eventRepo.SetupKnockoutBracket(eventID, schedule)
}

func (s *EventService) Repo() *repo.EventRepo {
	return s.eventRepo
}
