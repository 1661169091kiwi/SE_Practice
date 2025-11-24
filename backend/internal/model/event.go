package model

import "time"

// Event 赛事模型
type Event struct {
	ID        int64     `json:"event_id"`
	Name      string    `json:"event_name"`
	SportID   int64     `json:"sport_id"`
	Season    string    `json:"season"`
	Round     string    `json:"round"`
	Format    string    `json:"format_type"` // points/group_knockout
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"` // upcoming/ongoing/finished
	CreatedAt time.Time `json:"created_at"`
}

// CreateEventRequest 创建赛事请求
type CreateEventRequest struct {
	EventName string    `json:"event_name"`
	SportID   int64     `json:"sport_id"`
	Season    string    `json:"season"`
	Round     string    `json:"round"`
	Format    string    `json:"format_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// EventResponse 赛事响应模型
type EventResponse struct {
	EventID   int64     `json:"event_id"`
	EventName string    `json:"event_name"`
	SportID   int64     `json:"sport_id"`
	Status    string    `json:"status"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
