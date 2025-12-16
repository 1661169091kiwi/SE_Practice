package model

// Standings 积分榜条目
type Standings struct {
	EventID      int64  `json:"event_id"`
	TeamID       int64  `json:"team_id"`
	TeamName     string `json:"team_name,omitempty"` // 关联查询出的队名
	TeamLogo     string `json:"team_logo,omitempty"` // 关联查询出的队标
	Played       int    `json:"played"`
	Won          int    `json:"won"`
	Drawn        int    `json:"drawn"`
	Lost         int    `json:"lost"`
	GoalsFor     int    `json:"goals_for"`
	GoalsAgainst int    `json:"goals_against"`
	Points       int    `json:"points"`
	Rank         int    `json:"rank"`
}
