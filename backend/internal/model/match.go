package model

import "time"

// Match 比赛模型
type Match struct {
	ID           int64     `json:"match_id"`
	EventID      int64     `json:"event_id"`
	Name         string    `json:"match_name"`
	Round        string    `json:"round"`
	Time         time.Time `json:"match_time"`
	TeamAID      int64     `json:"team_a_id"`
	TeamBID      int64     `json:"team_b_id"`
	ScoreA       int       `json:"score_team_a"`
	ScoreB       int       `json:"score_team_b"`
	HalfScoreA   int       `json:"half_score_team_a"`
	HalfScoreB   int       `json:"half_score_team_b"`
	Status       string    `json:"status"` // not_started/ongoing/finished/cancelled
	Collector1ID int64     `json:"collector1_id"`
	Collector2ID int64     `json:"collector2_id"`
}

// CreateMatchRequest 创建比赛请求
type CreateMatchRequest struct {
	EventID   int64     `json:"event_id"`
	MatchName string    `json:"match_name"`
	Round     string    `json:"round"`
	MatchTime time.Time `json:"match_time"`
	TeamAID   int64     `json:"team_a_id"`
	TeamBID   int64     `json:"team_b_id"`
}

// MatchDetailResponse 比赛详情响应
type MatchDetailResponse struct {
	MatchID    int64       `json:"match_id"`
	EventID    int64       `json:"event_id"`
	MatchName  string      `json:"match_name"`
	Round      string      `json:"round"`
	MatchTime  time.Time   `json:"match_time"`
	TeamA      TeamBrief   `json:"team_a"`
	TeamB      TeamBrief   `json:"team_b"`
	ScoreA     int         `json:"score_team_a"`
	ScoreB     int         `json:"score_team_b"`
	Status     string      `json:"status"`
	Collectors []UserBrief `json:"collectors"`
}

// TeamBrief 队伍简要信息
type TeamBrief struct {
	ID   int64  `json:"team_id"`
	Name string `json:"team_name"`
}

// UserBrief 用户简要信息
type UserBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
