package model

import "time"

// Team 队伍模型
type Team struct {
	ID          int64     `json:"team_id"`
	TeamName    string    `json:"team_name"`
	SportID     int64     `json:"sport_id"`
	College     string    `json:"college"`
	TeamType    string    `json:"team_type"` // college/school
	AvatarURL   string    `json:"avatar_url"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"` // student_id of creator
	CreatedAt   time.Time `json:"created_at"`
	IsApproved  bool      `json:"is_approved"`
}

// CreateTeamRequest 创建队伍请求
type CreateTeamRequest struct {
	TeamName    string `json:"team_name"`
	SportID     int64  `json:"sport_id"`
	College     string `json:"college"`
	TeamType    string `json:"team_type"`
	AvatarURL   string `json:"avatar_url"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"` // Optional, can be taken from context
}

// TeamMemberDetail 队伍成员详情
type TeamMemberDetail struct {
	TeamMemberID int64  `json:"team_member_id"`
	TeamID       int64  `json:"team_id"`
	AthleteID    int64  `json:"athlete_id"`
	StudentID    string `json:"student_id"`
	Name         string `json:"name"`
	College      string `json:"college"`
	SportType    string `json:"sport_type"`
	JoinDate     string `json:"join_date"`
	JerseyNumber string `json:"jersey_number"`
}
