package model

// MatchStats 比赛详细数据
type MatchStats struct {
	MatchID           int64 `json:"match_id"`
	HomePossession    int   `json:"home_possession"`
	AwayPossession    int   `json:"away_possession"`
	HomeShots         int   `json:"home_shots"`
	AwayShots         int   `json:"away_shots"`
	HomeShotsOnTarget int   `json:"home_shots_on_target"`
	AwayShotsOnTarget int   `json:"away_shots_on_target"`
	HomeFouls         int   `json:"home_fouls"`
	AwayFouls         int   `json:"away_fouls"`
	HomeCorners       int   `json:"home_corners"`
	AwayCorners       int   `json:"away_corners"`
}
