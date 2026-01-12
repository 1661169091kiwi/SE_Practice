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
	EventName        string            `json:"event_name"`
	SportID          int64             `json:"sport_id"`
	Season           string            `json:"season"`
	Round            string            `json:"round"`
	Format           string            `json:"format_type"`
	StartDate        time.Time         `json:"start_date"`
	EndDate          time.Time         `json:"end_date"`
	Teams            []EventTeamInput  `json:"teams,omitempty"`
	KnockoutSchedule map[string]string `json:"knockout_schedule,omitempty"`
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

type GroupStandings struct {
	Group     string       `json:"group"`
	Standings []*Standings `json:"standings"`
}

type KnockoutMatchView struct {
	KnockoutMatchID int64  `json:"knockout_match_id"`
	StageID         int64  `json:"stage_id"`
	MatchID         int64  `json:"match_id"`
	MatchName       string `json:"match_name"`
	Round           string `json:"round"`
	MatchTime       string `json:"match_time"`
	Status          string `json:"status"`
	ScoreA          int    `json:"score_team_a"`
	ScoreB          int    `json:"score_team_b"`
	TeamAID         int64  `json:"team_a_id"`
	TeamAName       string `json:"team_a_name"`
	TeamALogo       string `json:"team_a_logo"`
	TeamBID         int64  `json:"team_b_id"`
	TeamBName       string `json:"team_b_name"`
	TeamBLogo       string `json:"team_b_logo"`
	PrevMatchAID    int64  `json:"prev_match_a_id"`
	PrevMatchBID    int64  `json:"prev_match_b_id"`
	WinnerTeamID    int64  `json:"winner_team_id"`
	ScheduledTime   string `json:"scheduled_time,omitempty"`
}

type KnockoutStageView struct {
	StageID    int64               `json:"stage_id"`
	StageName  string              `json:"stage_name"`
	StageOrder int                 `json:"stage_order"`
	Matches    []KnockoutMatchView `json:"matches"`
}

type EventStandingsOverview struct {
	Event          *Event              `json:"event"`
	FormatType     string              `json:"format_type"`
	League         []*Standings        `json:"league,omitempty"`
	Groups         []GroupStandings    `json:"groups,omitempty"`
	KnockoutStages []KnockoutStageView `json:"knockout_stages,omitempty"`
}

type EventTeamInput struct {
	TeamID    int64  `json:"team_id"`
	GroupName string `json:"group_name,omitempty"`
	Slot      int    `json:"slot,omitempty"`
}

type EventTeam struct {
	EventID   int64  `json:"event_id"`
	TeamID    int64  `json:"team_id"`
	TeamName  string `json:"team_name"`
	GroupName string `json:"group_name"`
	Slot      int    `json:"slot"`
}
