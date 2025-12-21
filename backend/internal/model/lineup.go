package model

type MatchLineup struct {
	LineupID     int64  `json:"lineup_id"`
	MatchID      int64  `json:"match_id"`
	TeamID       int64  `json:"team_id"`
	StudentID    string `json:"student_id"`
	Position     string `json:"position"`
	IsStarting   bool   `json:"is_starting"`
	JerseyNumber string `json:"jersey_number"`
	PlayerName   string `json:"player_name,omitempty"`
}
