package model

type MatchScores struct {
	TeamA int `json:"teamA"`
	TeamB int `json:"teamB"`
}

type PlayerBrief struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Number     string `json:"number"`
	Position   string `json:"position"`
	IsStarting bool   `json:"isStarting"`
}

type MatchLineups struct {
	TeamA []PlayerBrief `json:"teamA"`
	TeamB []PlayerBrief `json:"teamB"`
}

type CollectorTeamBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type CollectorMatchData struct {
	Scores  MatchScores  `json:"scores"`
	Lineups MatchLineups `json:"lineups"`
}

type CollectorMatchDetailResponse struct {
	ID              int64                `json:"id"`
	EventName       string               `json:"eventName"`
	Name            string               `json:"name"`
	TeamA           string               `json:"teamA"`
	TeamB           string               `json:"teamB"`
	TeamAID         int64                `json:"teamAId"`
	TeamBID         int64                `json:"teamBId"`
	Time            string               `json:"time"`
	Venue           string               `json:"venue"`
	SportType       string               `json:"sportType"`
	Status          string               `json:"status"`
	AutoFinishAt    string               `json:"autoFinishAt"`
	CandidatesTeamA []PlayerBrief        `json:"candidatesTeamA"`
	CandidatesTeamB []PlayerBrief        `json:"candidatesTeamB"`
	AvailableTeams  []CollectorTeamBrief `json:"availableTeams"`
	CurrentData     CollectorMatchData   `json:"currentData"`
}

type LineupItem struct {
	StudentID    string `json:"studentId"`
	Name         string `json:"name,omitempty"`
	Position     string `json:"position,omitempty"`
	JerseyNumber string `json:"jerseyNumber,omitempty"`
	IsStarting   bool   `json:"isStarting"`
}

type CollectorSubmitDataRequest struct {
	ScoreTeamA int          `json:"scoreTeamA"`
	ScoreTeamB int          `json:"scoreTeamB"`
	TeamA      []LineupItem `json:"teamA"`
	TeamB      []LineupItem `json:"teamB"`
}
