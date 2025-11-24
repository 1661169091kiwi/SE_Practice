package repo

import (
	"database/sql"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type MatchRepo struct{}

func NewMatchRepo() *MatchRepo {
	return &MatchRepo{}
}

// CreateMatch 创建比赛
func (r *MatchRepo) CreateMatch(match *model.Match) (int64, error) {
	query := `INSERT INTO matches 
	          (event_id, match_name, round, match_time, team_a_id, team_b_id, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query,
		match.EventID,
		match.Name,
		match.Round,
		match.Time,
		match.TeamAID,
		match.TeamBID,
		match.Status,
	)
}

// GetMatchByID 获取比赛详情
func (r *MatchRepo) GetMatchByID(id int64) (*model.Match, error) {
	query := `SELECT match_id, event_id, match_name, round, match_time, 
	                 team_a_id, team_b_id, score_team_a, score_team_b,
	                 half_score_team_a, half_score_team_b, status,
	                 collector1_id, collector2_id 
	          FROM matches WHERE match_id = ?`

	var match model.Match
	err := db.QueryRow(query, id).Scan(
		&match.ID,
		&match.EventID,
		&match.Name,
		&match.Round,
		&match.Time,
		&match.TeamAID,
		&match.TeamBID,
		&match.ScoreA,
		&match.ScoreB,
		&match.HalfScoreA,
		&match.HalfScoreB,
		&match.Status,
		&match.Collector1ID,
		&match.Collector2ID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &match, err
}

// ListMatchesByEvent 按赛事ID列出比赛
func (r *MatchRepo) ListMatchesByEvent(eventID int64) ([]model.Match, error) {
	query := `SELECT match_id, event_id, match_name, round, match_time, 
	                 team_a_id, team_b_id, status 
	          FROM matches WHERE event_id = ? ORDER BY match_time ASC`

	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &m.TeamBID, &m.Status,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// ListAllMatches 获取所有比赛
func (r *MatchRepo) ListAllMatches() ([]model.Match, error) {
	query := `SELECT match_id, event_id, match_name, round, match_time, 
	                 team_a_id, team_b_id, status 
	          FROM matches ORDER BY match_time ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &m.TeamBID, &m.Status,
		); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// UpdateMatchScore 更新比赛分数
func (r *MatchRepo) UpdateMatchScore(id int64, scoreA, scoreB int) error {
	query := `UPDATE matches SET 
	          score_team_a = ?, 
	          score_team_b = ?,
	          status = 'finished' 
	          WHERE match_id = ?`
	_, err := db.Exec(query, scoreA, scoreB, id)
	return err
}
