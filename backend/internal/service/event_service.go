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
}

func NewEventService() *EventService {
	return &EventService{
		eventRepo: repo.NewEventRepo(),
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
	if req.Format != "points" && req.Format != "group_knockout" {
		return nil, errors.New("invalid format type: must be 'points' or 'group_knockout'")
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

	// 更新字段
	event.Name = req.EventName
	event.SportID = req.SportID
	event.Season = req.Season
	event.Round = req.Round
	event.Format = req.Format
	event.StartDate = req.StartDate
	event.EndDate = req.EndDate

	return s.eventRepo.UpdateEvent(event)
}

func isValidEventStatus(status string) bool {
	return status == "upcoming" || status == "ongoing" || status == "finished"
}
