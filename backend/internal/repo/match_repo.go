package repo

import (
	"database/sql"
	"fmt"
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

// GetStandings 获取积分榜
// 优化逻辑：先根据 sportType 和 date 找到最匹配的一个赛事 (ongoing > upcoming > finished)，避免多个赛事的积分榜混合。
func (r *MatchRepo) GetStandings(sportType int, date string) ([]model.ResultItem, error) {
	// 1. 查找最相关的 event_id
	// 逻辑：查找该日期范围内正在进行或结束的赛事。
	// 输入 date 格式可能是 "YYYY-MM" 或 "YYYY-MM-DD"。统一截取前7位 "YYYY-MM" 进行比较。
	// 使用 DATE_FORMAT 进行月份级别的比较: start_date_month <= input_month <= end_date_month
	// 如果有多个，优先取 start_date 最近的
	queryEvent := `
		SELECT event_id FROM events 
		WHERE sport_id = ? 
		AND (DATE_FORMAT(start_date, '%Y-%m') <= ? AND DATE_FORMAT(end_date, '%Y-%m') >= ?)
		ORDER BY start_date DESC LIMIT 1
	`

	// 处理日期格式，确保是 YYYY-MM
	searchMonth := date
	if len(date) > 7 {
		searchMonth = date[:7]
	}

	var eventID int64
	err := db.QueryRow(queryEvent, sportType, searchMonth, searchMonth).Scan(&eventID)
	if err != nil {
		if err == sql.ErrNoRows {
			// 如果没有找到时间跨度内的赛事，尝试找在该月份开始的赛事作为兜底
			fallbackQuery := `SELECT event_id FROM events WHERE sport_id = ? AND start_date LIKE ? ORDER BY start_date DESC LIMIT 1`
			err = db.QueryRow(fallbackQuery, sportType, searchMonth+"%").Scan(&eventID)
			if err != nil {
				if err == sql.ErrNoRows {
					return []model.ResultItem{}, nil // 无相关赛事
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// 2. 根据 event_id 查询积分榜
	queryStandings := `
		SELECT 
			t.team_name, t.avatar_url,
			s.matches_played, s.wins, s.draws, s.losses, 
			s.goals_for, s.goals_against, s.points
		FROM standings s
		JOIN teams t ON s.team_id = t.team_id
		WHERE s.event_id = ?
		ORDER BY s.points DESC, s.goal_difference DESC
	`

	rows, err := db.Query(queryStandings, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ResultItem
	rank := 1
	for rows.Next() {
		var r model.ResultItem
		var goalsFor, goalsAgainst int
		var avatar sql.NullString

		if err := rows.Scan(
			&r.Name, &avatar,
			&r.Played, &r.Won, &r.Drawn, &r.Lost,
			&goalsFor, &goalsAgainst, &r.Points,
		); err != nil {
			return nil, err
		}
		if avatar.Valid {
			r.ImgURL = avatar.String
		}
		r.Rank = rank
		r.Goals = fmt.Sprintf("%d-%d", goalsFor, goalsAgainst)
		results = append(results, r)
		rank++
	}
	return results, nil
}

// GetSchedule 获取赛程 (status = 'upcoming' or 'ongoing')
func (r *MatchRepo) GetSchedule(sportType int, date string) ([]model.ResultItem, error) {
	query := `
		SELECT 
			m.match_name, m.match_time, m.status,
			ta.team_name, tb.team_name
		FROM matches m
		JOIN events e ON m.event_id = e.event_id
		JOIN teams ta ON m.team_a_id = ta.team_id
		JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE e.sport_id = ? AND m.match_time LIKE ? AND m.status IN ('upcoming', 'ongoing')
		ORDER BY m.match_time ASC
	`
	return r.queryMatchList(query, sportType, date)
}

// GetHistory 获取历史战绩 (status = 'finished')
func (r *MatchRepo) GetHistory(sportType int, date string) ([]model.ResultItem, error) {
	query := `
		SELECT 
			m.match_name, m.match_time, m.status,
			ta.team_name, tb.team_name,
			m.score_team_a, m.score_team_b
		FROM matches m
		JOIN events e ON m.event_id = e.event_id
		JOIN teams ta ON m.team_a_id = ta.team_id
		JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE e.sport_id = ? AND m.match_time LIKE ? AND m.status = 'finished'
		ORDER BY m.match_time DESC
	`
	// 复用逻辑 slightly differ, manually scan here for completeness
	datePattern := date + "%"
	rows, err := db.Query(query, sportType, datePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ResultItem
	for rows.Next() {
		var item model.ResultItem
		var teamA, teamB string
		var scoreA, scoreB int

		if err := rows.Scan(
			&item.Name, &item.MatchTime, &item.Status,
			&teamA, &teamB, &scoreA, &scoreB,
		); err != nil {
			return nil, err
		}
		// 格式化 Name 为 "TeamA vs TeamB"
		item.Name = fmt.Sprintf("%s vs %s", teamA, teamB)
		// 格式化 Goals 为 "ScoreA-ScoreB"
		item.Goals = fmt.Sprintf("%d-%d", scoreA, scoreB)
		results = append(results, item)
	}
	return results, nil
}

// queryMatchList 辅助方法查询赛程
func (r *MatchRepo) queryMatchList(query string, sportType int, date string) ([]model.ResultItem, error) {
	datePattern := date + "%"
	rows, err := db.Query(query, sportType, datePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ResultItem
	for rows.Next() {
		var item model.ResultItem
		var teamA, teamB string

		// Note: Schedule doesn't have scores usually
		if err := rows.Scan(
			&item.Name, &item.MatchTime, &item.Status,
			&teamA, &teamB,
		); err != nil {
			return nil, err
		}
		item.Name = fmt.Sprintf("%s vs %s", teamA, teamB)
		results = append(results, item)
	}
	return results, nil
}

// GetMatchesBySubscription 获取用户订阅的比赛
func (r *MatchRepo) GetMatchesBySubscription(studentID string) ([]model.SubscribedMatchItem, error) {
	// 逻辑：
	// 1. matches JOIN subscriptions ON event_id
	// 2. JOIN teams ta/tb to get names and avatars
	// 3. WHERE subscriptions.student_id = ?
	// 4. ORDER BY match_time ASC
	query := `
		SELECT 
			m.match_id, m.match_time, m.status, m.score_team_a, m.score_team_b,
			ta.team_name, ta.avatar_url,
			tb.team_name, tb.avatar_url
		FROM matches m
		JOIN subscriptions s ON m.event_id = s.event_id
		JOIN teams ta ON m.team_a_id = ta.team_id
		JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE s.student_id = ?
		ORDER BY m.match_time ASC
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.SubscribedMatchItem
	for rows.Next() {
		var item model.SubscribedMatchItem
		var scoreA, scoreB int
		var status string
		var matchTime string
		var taName, tbName string
		var taAvatar, tbAvatar sql.NullString

		if err := rows.Scan(
			&item.MatchID, &matchTime, &status, &scoreA, &scoreB,
			&taName, &taAvatar,
			&tbName, &tbAvatar,
		); err != nil {
			return nil, err
		}

		// 处理 Nullable Avatar
		if taAvatar.Valid {
			item.HomeTeam.Avatar = taAvatar.String
		}
		if tbAvatar.Valid {
			item.AwayTeam.Avatar = tbAvatar.String
		}

		item.HomeTeam.Name = taName
		item.AwayTeam.Name = tbName
		item.MatchTime = matchTime
		item.MatchVenue = "校体育馆" // Placeholder as requested

		// 简单的状态转换 (更复杂的逻辑可以放在 Service 层，或者这里处理)
		if status == "not_started" {
			item.MatchStatus = "VS"
			item.MatchState = "未开始"
		} else if status == "finished" {
			item.MatchStatus = fmt.Sprintf("%d-%d", scoreA, scoreB)
			item.MatchState = "已结束"
		} else {
			// ongoing or others
			item.MatchStatus = fmt.Sprintf("%d-%d", scoreA, scoreB)
			item.MatchState = "进行中"
		}

		results = append(results, item)
	}
	return results, nil
}

// GetEventIDByMatchID 根据比赛ID查找赛事ID
func (r *MatchRepo) GetEventIDByMatchID(matchID string) (int64, error) {
	query := `SELECT event_id FROM matches WHERE match_id = ?`
	var eventID int64
	err := db.QueryRow(query, matchID).Scan(&eventID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return eventID, err
}

// CheckSubscription 检查是否已订阅
func (r *MatchRepo) CheckSubscription(studentID string, eventID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM subscriptions WHERE student_id = ? AND event_id = ?`
	var count int
	err := db.QueryRow(query, studentID, eventID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SubscribeToEvent 订阅赛事
func (r *MatchRepo) SubscribeToEvent(studentID string, eventID int64) error {
	query := `INSERT INTO subscriptions (student_id, event_id) VALUES (?, ?)`
	_, err := db.Exec(query, studentID, eventID)
	return err
}

// UnsubscribeFromEvent 取消订阅赛事
func (r *MatchRepo) UnsubscribeFromEvent(studentID string, eventID int64) error {
	query := `DELETE FROM subscriptions WHERE student_id = ? AND event_id = ?`
	_, err := db.Exec(query, studentID, eventID)
	return err
}
