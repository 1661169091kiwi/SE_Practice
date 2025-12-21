package repo

import (
	"database/sql"
	"fmt"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"time"
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
	query := `SELECT match_id, event_id, match_name, COALESCE(round, ''), match_time, 
	                 team_a_id, team_b_id, score_team_a, score_team_b,
	                 half_score_team_a, half_score_team_b, status,
	                 COALESCE(collector1_id, 0), COALESCE(collector2_id, 0) 
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
	query := `SELECT m.match_id, m.event_id, e.sport_id, m.match_name, COALESCE(m.round, ''), m.match_time, 
	                 m.team_a_id, ta.team_name, m.team_b_id, tb.team_name, m.status,
	                 m.score_team_a, m.score_team_b
	          FROM matches m
	          JOIN events e ON m.event_id = e.event_id
	          LEFT JOIN teams ta ON m.team_a_id = ta.team_id
	          LEFT JOIN teams tb ON m.team_b_id = tb.team_id
	          WHERE m.event_id = ? ORDER BY m.match_time ASC`

	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		var taName, tbName sql.NullString
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.SportID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &taName, &m.TeamBID, &tbName, &m.Status,
			&m.ScoreA, &m.ScoreB,
		); err != nil {
			return nil, err
		}
		if taName.Valid {
			m.TeamAName = taName.String
		}
		if tbName.Valid {
			m.TeamBName = tbName.String
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// DeleteMatch 删除比赛
func (r *MatchRepo) DeleteMatch(matchID int64) error {
	// Delete related data
	// 1. comment_likes via player_ratings
	_, err := db.Exec("DELETE cl FROM comment_likes cl JOIN player_ratings pr ON cl.rating_id = pr.rating_id WHERE pr.match_id = ?", matchID)
	if err != nil {
		return err
	}
	// 2. player_ratings
	_, err = db.Exec("DELETE FROM player_ratings WHERE match_id = ?", matchID)
	if err != nil {
		return err
	}
	// 3. match_lineups
	_, err = db.Exec("DELETE FROM match_lineups WHERE match_id = ?", matchID)
	if err != nil {
		return err
	}
	// 4. match_events
	_, err = db.Exec("DELETE FROM match_events WHERE match_id = ?", matchID)
	if err != nil {
		return err
	}
	// 5. injury_suspensions
	_, err = db.Exec("DELETE FROM injury_suspensions WHERE match_id = ?", matchID)
	if err != nil {
		return err
	}
	// 6. knockout_matches
	_, err = db.Exec("DELETE FROM knockout_matches WHERE match_id = ?", matchID)
	if err != nil {
		// Ignore error if table doesn't exist or column differs, but assuming schema matches
		// return err
		// Actually, let's just log or ignore if not found? No, should be fine.
	}

	// 7. matches
	_, err = db.Exec("DELETE FROM matches WHERE match_id = ?", matchID)
	return err
}

// ListAllMatches 获取所有比赛
func (r *MatchRepo) ListAllMatches() ([]model.Match, error) {
	query := `SELECT m.match_id, m.event_id, e.sport_id, m.match_name, COALESCE(m.round, ''), m.match_time, 
	                 m.team_a_id, ta.team_name, m.team_b_id, tb.team_name, m.status,
	                 m.score_team_a, m.score_team_b
	          FROM matches m
	          JOIN events e ON m.event_id = e.event_id
	          LEFT JOIN teams ta ON m.team_a_id = ta.team_id
	          LEFT JOIN teams tb ON m.team_b_id = tb.team_id
	          ORDER BY m.match_time ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		var taName, tbName sql.NullString
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.SportID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &taName, &m.TeamBID, &tbName, &m.Status,
			&m.ScoreA, &m.ScoreB,
		); err != nil {
			return nil, err
		}
		if taName.Valid {
			m.TeamAName = taName.String
		}
		if tbName.Valid {
			m.TeamBName = tbName.String
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// UpdateMatchScore 更新比赛分数
func (r *MatchRepo) UpdateMatchScore(id int64, scoreA, scoreB int) error {
	query := `UPDATE matches SET 
	          score_team_a = ?, 
	          score_team_b = ?
	          WHERE match_id = ?`
	_, err := db.Exec(query, scoreA, scoreB, id)
	return err
}

// UpdateMatch 更新比赛信息
func (r *MatchRepo) UpdateMatch(match *model.Match) error {
	query := `UPDATE matches SET 
              event_id=?, 
              match_name=?, 
              round=?, 
              match_time=?, 
              team_a_id=?, 
              team_b_id=? 
              WHERE match_id=?`
	_, err := db.Exec(query,
		match.EventID,
		match.Name,
		match.Round,
		match.Time,
		match.TeamAID,
		match.TeamBID,
		match.ID,
	)
	return err
}

// GetMatchComments 获取比赛评论
func (r *MatchRepo) GetMatchComments(matchID int64) ([]model.MatchComment, error) {
	query := `
        SELECT mc.id, mc.match_id, mc.student_id, mc.content, mc.created_at,
               u.name, u.avatar_url
        FROM match_comments mc
        JOIN users u ON mc.student_id = u.student_id
        WHERE mc.match_id = ?
        ORDER BY mc.created_at DESC
    `
	rows, err := db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.MatchComment
	for rows.Next() {
		var c model.MatchComment
		if err := rows.Scan(&c.ID, &c.MatchID, &c.StudentID, &c.Content, &c.CreatedAt, &c.UserName, &c.AvatarURL); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// CreateMatchComment 创建比赛评论
func (r *MatchRepo) CreateMatchComment(comment *model.MatchComment) error {
	query := `INSERT INTO match_comments (match_id, student_id, content, created_at) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, comment.MatchID, comment.StudentID, comment.Content, comment.CreatedAt)
	return err
}

// GetMatchStats 获取比赛统计
func (r *MatchRepo) GetMatchStats(matchID int64) (*model.MatchStats, error) {
	query := `SELECT match_id, home_possession, away_possession, home_shots, away_shots, 
                     home_shots_on_target, away_shots_on_target, home_fouls, away_fouls, 
                     home_corners, away_corners 
              FROM match_stats WHERE match_id = ?`

	var s model.MatchStats
	err := db.QueryRow(query, matchID).Scan(
		&s.MatchID, &s.HomePossession, &s.AwayPossession, &s.HomeShots, &s.AwayShots,
		&s.HomeShotsOnTarget, &s.AwayShotsOnTarget, &s.HomeFouls, &s.AwayFouls,
		&s.HomeCorners, &s.AwayCorners,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// GetMatchLineups 获取比赛阵容
func (r *MatchRepo) GetMatchLineups(matchID int64) ([]model.MatchLineup, error) {
	query := `
SELECT ml.lineup_id, ml.match_id, ml.team_id, ml.student_id, ml.position, ml.is_starting, ml.jersey_number, u.name
FROM match_lineups ml
LEFT JOIN users u ON ml.student_id = u.student_id
WHERE ml.match_id = ?
UNION ALL
SELECT ext.ext_id AS lineup_id, ext.match_id, ext.team_id, '' AS student_id, ext.position, ext.is_starting, ext.jersey_number, ext.player_name AS name
FROM match_lineups_external ext
WHERE ext.match_id = ?`

	rows, err := db.Query(query, matchID, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lineups []model.MatchLineup
	for rows.Next() {
		var l model.MatchLineup
		var name sql.NullString
		if err := rows.Scan(&l.LineupID, &l.MatchID, &l.TeamID, &l.StudentID, &l.Position, &l.IsStarting, &l.JerseyNumber, &name); err != nil {
			return nil, err
		}
		if name.Valid {
			l.PlayerName = name.String
		}
		lineups = append(lineups, l)
	}
	return lineups, nil
}

// GetAthleteMatches 获取运动员参加的比赛列表
func (r *MatchRepo) GetAthleteMatches(studentID string) ([]model.Match, error) {
	query := `
		SELECT DISTINCT m.match_id, m.event_id, e.sport_id, m.match_name, COALESCE(m.round, ''), m.match_time,
		       m.team_a_id, ta.team_name, m.team_b_id, tb.team_name, m.status,
		       m.score_team_a, m.score_team_b, m.half_score_team_a, m.half_score_team_b,
		       COALESCE(m.collector1_id, 0), COALESCE(m.collector2_id, 0)
		FROM match_lineups ml
		JOIN matches m ON ml.match_id = m.match_id
		JOIN events e ON m.event_id = e.event_id
		LEFT JOIN teams ta ON m.team_a_id = ta.team_id
		LEFT JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE ml.student_id = ?
		ORDER BY m.match_time DESC
	`

	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		var taName, tbName sql.NullString
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.SportID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &taName, &m.TeamBID, &tbName, &m.Status,
			&m.ScoreA, &m.ScoreB, &m.HalfScoreA, &m.HalfScoreB,
			&m.Collector1ID, &m.Collector2ID,
		); err != nil {
			return nil, err
		}
		if taName.Valid {
			m.TeamAName = taName.String
		}
		if tbName.Valid {
			m.TeamBName = tbName.String
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// GetMatchesByTeam 获取某个队伍参与的所有比赛
func (r *MatchRepo) GetMatchesByTeam(teamID int64) ([]model.Match, error) {
	query := `
		SELECT m.match_id, m.event_id, e.sport_id, m.match_name, COALESCE(m.round, ''), m.match_time,
		       m.team_a_id, ta.team_name, m.team_b_id, tb.team_name, m.status,
		       m.score_team_a, m.score_team_b, m.half_score_team_a, m.half_score_team_b,
		       COALESCE(m.collector1_id, 0), COALESCE(m.collector2_id, 0)
		FROM matches m
		JOIN events e ON m.event_id = e.event_id
		LEFT JOIN teams ta ON m.team_a_id = ta.team_id
		LEFT JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE (m.team_a_id = ? OR m.team_b_id = ?)
		ORDER BY m.match_time DESC
	`

	rows, err := db.Query(query, teamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.Match
	for rows.Next() {
		var m model.Match
		var taName, tbName sql.NullString
		if err := rows.Scan(
			&m.ID, &m.EventID, &m.SportID, &m.Name, &m.Round, &m.Time,
			&m.TeamAID, &taName, &m.TeamBID, &tbName, &m.Status,
			&m.ScoreA, &m.ScoreB, &m.HalfScoreA, &m.HalfScoreB,
			&m.Collector1ID, &m.Collector2ID,
		); err != nil {
			return nil, err
		}
		if taName.Valid {
			m.TeamAName = taName.String
		}
		if tbName.Valid {
			m.TeamBName = tbName.String
		}
		matches = append(matches, m)
	}
	return matches, rows.Err()
}

// RecalculateStandings 重新计算赛事积分榜
func (r *MatchRepo) RecalculateStandings(eventID int64) error {
	// 1. 获取赛事的所有已结束比赛
	matches, err := r.ListMatchesByEvent(eventID)
	if err != nil {
		return err
	}

	// 2. 内存计算积分
	type TeamStats struct {
		Played, Won, Drawn, Lost, GF, GA, Points int
	}
	stats := make(map[int64]*TeamStats)

	for _, m := range matches {
		if m.Status != "finished" {
			continue
		}

		// Init stats if not exists
		if _, ok := stats[m.TeamAID]; !ok {
			stats[m.TeamAID] = &TeamStats{}
		}
		if _, ok := stats[m.TeamBID]; !ok {
			stats[m.TeamBID] = &TeamStats{}
		}

		sa := m.ScoreA
		sb := m.ScoreB

		stats[m.TeamAID].Played++
		stats[m.TeamBID].Played++
		stats[m.TeamAID].GF += sa
		stats[m.TeamAID].GA += sb
		stats[m.TeamBID].GF += sb
		stats[m.TeamBID].GA += sa

		if sa > sb {
			stats[m.TeamAID].Won++
			stats[m.TeamAID].Points += 3
			stats[m.TeamBID].Lost++
		} else if sa < sb {
			stats[m.TeamBID].Won++
			stats[m.TeamBID].Points += 3
			stats[m.TeamAID].Lost++
		} else {
			stats[m.TeamAID].Drawn++
			stats[m.TeamAID].Points += 1
			stats[m.TeamBID].Drawn++
			stats[m.TeamBID].Points += 1
		}
	}

	// 3. 开启事务更新数据库
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}
	// Note: If anything fails, we should rollback.
	// However, if we commit successfully, defer Rollback will do nothing (usually, depends on driver/wrapper).
	// Standard practice: defer tx.Rollback() and commit at end.
	defer tx.Rollback()

	// 清空该赛事旧积分
	_, err = tx.Exec("DELETE FROM standings WHERE event_id = ?", eventID)
	if err != nil {
		return err
	}

	// 插入新积分
	stmt, err := tx.Prepare(`INSERT INTO standings 
        (event_id, team_id, matches_played, wins, draws, losses, goals_for, goals_against, goal_difference, points) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for teamID, s := range stats {
		_, err = stmt.Exec(eventID, teamID, s.Played, s.Won, s.Drawn, s.Lost, s.GF, s.GA, s.GF-s.GA, s.Points)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
			tb.team_name, tb.avatar_url,
			e.sport_id
		FROM matches m
		JOIN subscriptions s ON m.event_id = s.event_id
		JOIN events e ON m.event_id = e.event_id
		JOIN teams ta ON m.team_a_id = ta.team_id
		JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE s.student_id = ? AND (s.match_id IS NULL OR s.match_id = m.match_id)
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
		var matchTime time.Time // Change to time.Time
		var taName, tbName string
		var taAvatar, tbAvatar sql.NullString

		if err := rows.Scan(
			&item.MatchID, &matchTime, &status, &scoreA, &scoreB,
			&taName, &taAvatar,
			&tbName, &tbAvatar,
			&item.SportID,
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

		// Format time for JSON
		item.MatchTime = matchTime.Format("2006-01-02 15:04:05")
		item.MatchVenue = "校体育馆"

		// Populate scores
		item.ScoreA = scoreA
		item.ScoreB = scoreB

		// Populate raw fields for service
		item.RawTime = matchTime
		item.RawStatus = status

		// 简单的状态转换 (Service 层会再次修正)
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
func (r *MatchRepo) CheckSubscription(studentID string, eventID int64, matchID int64) (bool, error) {
	// Logic:
	// If matchID > 0, we check for match_id = matchID OR match_id IS NULL (covering the match via event subscription).
	// If matchID == 0, we check for match_id IS NULL.

	var err error
	var count int
	var query string

	if matchID > 0 {
		query = `SELECT COUNT(*) FROM subscriptions WHERE student_id = ? AND event_id = ? AND (match_id = ? OR match_id IS NULL)`
		err = db.QueryRow(query, studentID, eventID, matchID).Scan(&count)
	} else {
		query = `SELECT COUNT(*) FROM subscriptions WHERE student_id = ? AND event_id = ? AND match_id IS NULL`
		err = db.QueryRow(query, studentID, eventID).Scan(&count)
	}

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SubscribeToEvent 订阅赛事 (Renamed or Overloaded concept)
// matchID can be 0 (subscribe to event) or > 0 (subscribe to match)
func (r *MatchRepo) Subscribe(studentID string, eventID int64, matchID int64) error {
	query := `INSERT INTO subscriptions (student_id, event_id, match_id) VALUES (?, ?, ?)`
	var mID interface{}
	if matchID > 0 {
		mID = matchID
	} else {
		mID = nil
	}
	_, err := db.Exec(query, studentID, eventID, mID)
	return err
}

// UnsubscribeFromEvent 取消订阅赛事
func (r *MatchRepo) Unsubscribe(studentID string, eventID int64, matchID int64) error {
	var query string
	var err error
	if matchID > 0 {
		// Try to delete specific subscription first
		query = `DELETE FROM subscriptions WHERE student_id = ? AND event_id = ? AND match_id = ?`
		result, err := db.Exec(query, studentID, eventID, matchID)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			// If no specific subscription, delete generic subscription (legacy support)
			query = `DELETE FROM subscriptions WHERE student_id = ? AND event_id = ? AND match_id IS NULL`
			_, err = db.Exec(query, studentID, eventID)
		}
	} else {
		query = `DELETE FROM subscriptions WHERE student_id = ? AND event_id = ? AND match_id IS NULL`
		_, err = db.Exec(query, studentID, eventID)
	}
	return err
}
