package service

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type StandingsService struct{}

func NewStandingsService() *StandingsService {
	return &StandingsService{}
}

// RecalculateStandings 重新计算指定赛事的积分榜
func (s *StandingsService) RecalculateStandings(eventID int64) error {
	// 1. 获取该赛事所有已结束的比赛
	rows, err := db.DB.Query("SELECT team_a_id, team_b_id, score_team_a, score_team_b FROM matches WHERE event_id = ? AND status = 'finished'", eventID)
	if err != nil {
		return err
	}
	defer rows.Close()

	statsMap := make(map[int64]*model.Standings)

	// 辅助函数，获取或初始化统计数据
	getStat := func(teamID int64) *model.Standings {
		if _, ok := statsMap[teamID]; !ok {
			statsMap[teamID] = &model.Standings{
				EventID: eventID,
				TeamID:  teamID,
			}
		}
		return statsMap[teamID]
	}

	teamRows, err := db.DB.Query("SELECT team_id FROM event_teams WHERE event_id = ?", eventID)
	if err != nil {
		return err
	}
	for teamRows.Next() {
		var teamID int64
		if err := teamRows.Scan(&teamID); err != nil {
			_ = teamRows.Close()
			return err
		}
		_ = getStat(teamID)
	}
	if err := teamRows.Err(); err != nil {
		_ = teamRows.Close()
		return err
	}
	_ = teamRows.Close()

	for rows.Next() {
		var teamA, teamB int64
		var scoreA, scoreB int
		if err := rows.Scan(&teamA, &teamB, &scoreA, &scoreB); err != nil {
			return err
		}

		statA := getStat(teamA)
		statB := getStat(teamB)

		statA.Played++
		statB.Played++
		statA.GoalsFor += scoreA
		statA.GoalsAgainst += scoreB
		statB.GoalsFor += scoreB
		statB.GoalsAgainst += scoreA

		if scoreA > scoreB {
			statA.Won++
			statA.Points += 3
			statB.Lost++
		} else if scoreA < scoreB {
			statB.Won++
			statB.Points += 3
			statA.Lost++
		} else {
			statA.Drawn++
			statA.Points += 1
			statB.Drawn++
			statB.Points += 1
		}
	}

	// 转换为切片并排序
	var standings []*model.Standings
	for _, stat := range statsMap {
		standings = append(standings, stat)
	}

	// 排序规则：积分 > 净胜球 > 进球数
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		gdI := standings[i].GoalsFor - standings[i].GoalsAgainst
		gdJ := standings[j].GoalsFor - standings[j].GoalsAgainst
		if gdI != gdJ {
			return gdI > gdJ
		}
		return standings[i].GoalsFor > standings[j].GoalsFor
	})

	// 更新排名并写入数据库
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}

	// 删除旧数据
	if _, err := tx.Exec("DELETE FROM standings WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 插入新数据
	stmt, err := tx.Prepare("INSERT INTO standings (event_id, team_id, matches_played, wins, draws, losses, goals_for, goals_against, goal_difference, points) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for i, stat := range standings {
		stat.Rank = i + 1
		goalDiff := stat.GoalsFor - stat.GoalsAgainst
		if _, err := stmt.Exec(stat.EventID, stat.TeamID, stat.Played, stat.Won, stat.Drawn, stat.Lost, stat.GoalsFor, stat.GoalsAgainst, goalDiff, stat.Points); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetStandings 获取积分榜
func (s *StandingsService) GetStandings(eventID int64) ([]*model.Standings, error) {
	// 注意：这里假设 teams 表有 team_id, team_name, avatar_url 字段。
	// 这里使用 LEFT JOIN 确保即使 team 信息缺失也能查出积分（虽然业务上不应该发生）
	query := `
		SELECT s.event_id, s.team_id, COALESCE(t.team_name, 'Unknown'), COALESCE(t.avatar_url, ''), s.matches_played, s.wins, s.draws, s.losses, s.goals_for, s.goals_against, s.points
		FROM standings s
		LEFT JOIN teams t ON s.team_id = t.team_id
		WHERE s.event_id = ?
		ORDER BY s.points DESC, s.goal_difference DESC, s.goals_for DESC
	`
	rows, err := db.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Standings
	for rows.Next() {
		var st model.Standings
		if err := rows.Scan(&st.EventID, &st.TeamID, &st.TeamName, &st.TeamLogo, &st.Played, &st.Won, &st.Drawn, &st.Lost, &st.GoalsFor, &st.GoalsAgainst, &st.Points); err != nil {
			return nil, err
		}
		result = append(result, &st)
	}
	rows.Close()

	// 如果没有数据，尝试重新计算（初始化）
	if len(result) == 0 {
		if err := s.RecalculateStandings(eventID); err == nil {
			rows, err = db.DB.Query(query, eventID)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var st model.Standings
					if err := rows.Scan(&st.EventID, &st.TeamID, &st.TeamName, &st.TeamLogo, &st.Played, &st.Won, &st.Drawn, &st.Lost, &st.GoalsFor, &st.GoalsAgainst, &st.Points); err == nil {
						result = append(result, &st)
					}
				}
			}
		}
	}

	// 计算排名
	for i, st := range result {
		st.Rank = i + 1
	}

	return result, nil
}

func (s *StandingsService) GetStandingsOverview(eventID int64) (*model.EventStandingsOverview, error) {
	var ev model.Event
	var endDate sql.NullTime
	err := db.DB.QueryRow(
		"SELECT event_id, event_name, sport_id, COALESCE(season, ''), COALESCE(round, ''), COALESCE(format_type, ''), start_date, end_date, status, created_at FROM events WHERE event_id = ?",
		eventID,
	).Scan(&ev.ID, &ev.Name, &ev.SportID, &ev.Season, &ev.Round, &ev.Format, &ev.StartDate, &endDate, &ev.Status, &ev.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found")
	}
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		ev.EndDate = endDate.Time
	}

	ov := &model.EventStandingsOverview{Event: &ev, FormatType: ev.Format}
	if ev.Format == "points" {
		list, err := s.GetStandings(eventID)
		if err != nil {
			return nil, err
		}
		ov.League = list
		return ov, nil
	}
	if ev.Format == "knockout" {
		knockout, err := s.getKnockoutStages(eventID)
		if err != nil {
			return nil, err
		}
		ov.KnockoutStages = knockout
		return ov, nil
	}

	groups, err := s.getGroupStandings(eventID)
	if err != nil {
		return nil, err
	}
	ov.Groups = groups

	knockout, err := s.getKnockoutStages(eventID)
	if err != nil {
		return nil, err
	}
	ov.KnockoutStages = knockout

	return ov, nil
}

func (s *StandingsService) UpdateKnockoutAdvancement(eventID, finishedMatchID int64) error {
	var teamAID, teamBID int64
	var scoreA, scoreB int
	var status string
	err := db.DB.QueryRow(
		"SELECT team_a_id, team_b_id, score_team_a, score_team_b, status FROM matches WHERE match_id = ? AND event_id = ?",
		finishedMatchID,
		eventID,
	).Scan(&teamAID, &teamBID, &scoreA, &scoreB, &status)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	if status != "finished" {
		return nil
	}
	if scoreA == scoreB {
		return nil
	}

	winnerTeamID := teamAID
	if scoreB > scoreA {
		winnerTeamID = teamBID
	}

	var finishedKnockoutMatchID int64
	err = db.DB.QueryRow("SELECT knockout_match_id FROM knockout_matches WHERE match_id = ?", finishedMatchID).Scan(&finishedKnockoutMatchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // not a knockout match
		}
		return err
	}

	getWinner := func(knockoutMatchID int64) (int64, bool, error) {
		var matchID int64
		err := db.DB.QueryRow("SELECT match_id FROM knockout_matches WHERE knockout_match_id = ?", knockoutMatchID).Scan(&matchID)
		if err != nil {
			if err == sql.ErrNoRows {
				return 0, false, nil
			}
			return 0, false, err
		}
		if matchID == 0 {
			return 0, false, nil
		}

		var a, b int64
		var sa, sb int
		var st string
		err = db.DB.QueryRow(
			"SELECT team_a_id, team_b_id, score_team_a, score_team_b, status FROM matches WHERE match_id = ? AND event_id = ?",
			matchID,
			eventID,
		).Scan(&a, &b, &sa, &sb, &st)
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		if err != nil {
			return 0, false, err
		}
		if st != "finished" {
			return 0, false, nil
		}
		if sa == sb {
			return 0, false, nil
		}
		if sa > sb {
			return a, true, nil
		}
		return b, true, nil
	}

	rows, err := db.DB.Query(
		"SELECT km.knockout_match_id, km.stage_id, km.match_id, COALESCE(km.prev_match_a_id, 0), COALESCE(km.prev_match_b_id, 0), COALESCE(ks.stage_name, ''), COALESCE(km.match_order, 0) FROM knockout_matches km JOIN knockout_stages ks ON km.stage_id = ks.stage_id WHERE ks.event_id = ? AND (km.prev_match_a_id = ? OR km.prev_match_b_id = ?)",
		eventID,
		finishedKnockoutMatchID,
		finishedKnockoutMatchID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var knockoutMatchID int64
		var stageID int64
		var matchID sql.NullInt64
		var prevAID, prevBID int64
		var stageName string
		var matchOrder int
		if err := rows.Scan(&knockoutMatchID, &stageID, &matchID, &prevAID, &prevBID, &stageName, &matchOrder); err != nil {
			return err
		}

		var aWinner, bWinner int64
		if prevAID > 0 {
			if prevAID == finishedKnockoutMatchID {
				aWinner = winnerTeamID
			} else {
				w, ok, err := getWinner(prevAID)
				if err != nil {
					return err
				}
				if ok {
					aWinner = w
				}
			}
		}
		if prevBID > 0 {
			if prevBID == finishedKnockoutMatchID {
				bWinner = winnerTeamID
			} else {
				w, ok, err := getWinner(prevBID)
				if err != nil {
					return err
				}
				if ok {
					bWinner = w
				}
			}
		}

		if matchID.Valid && matchID.Int64 > 0 {
			if aWinner > 0 {
				if _, err := db.DB.Exec("UPDATE matches SET team_a_id = ? WHERE match_id = ?", aWinner, matchID.Int64); err != nil {
					return err
				}
			}
			if bWinner > 0 {
				if _, err := db.DB.Exec("UPDATE matches SET team_b_id = ? WHERE match_id = ?", bWinner, matchID.Int64); err != nil {
					return err
				}
			}
			continue
		}

		if aWinner <= 0 || bWinner <= 0 {
			continue
		}

		matchName := fmt.Sprintf("%s %d", stageName, matchOrder)

		// 获取 scheduled_time
		var scheduledTime sql.NullString
		_ = db.DB.QueryRow("SELECT scheduled_time FROM knockout_matches WHERE knockout_match_id = ?", knockoutMatchID).Scan(&scheduledTime)

		var args []interface{}
		args = append(args, eventID, matchName, stageName)

		if scheduledTime.Valid && scheduledTime.String != "" {
			args = append(args, scheduledTime.String)
		} else {
			// If no scheduled time, use NOW() but we need to pass it as part of query string logic or handle args carefully
			// To simplify, if scheduledTime is NOT valid, we use NOW() in SQL
		}

		query := ""
		if scheduledTime.Valid && scheduledTime.String != "" {
			query = "INSERT INTO matches (event_id, match_name, round, match_time, team_a_id, team_b_id, status) VALUES (?, ?, ?, ?, ?, ?, 'not_started')"
			args = append(args, aWinner, bWinner)
		} else {
			query = "INSERT INTO matches (event_id, match_name, round, match_time, team_a_id, team_b_id, status) VALUES (?, ?, ?, '1000-01-01 00:00:00', ?, ?, 'not_started')"
			args = append(args, aWinner, bWinner)
		}

		res, err := db.DB.Exec(query, args...)
		if err != nil {
			return err
		}
		newMatchID, _ := res.LastInsertId()
		upd, err := db.DB.Exec("UPDATE knockout_matches SET match_id = ? WHERE knockout_match_id = ? AND match_id IS NULL", newMatchID, knockoutMatchID)
		if err != nil {
			_, _ = db.DB.Exec("DELETE FROM matches WHERE match_id = ?", newMatchID)
			return err
		}
		ra, _ := upd.RowsAffected()
		if ra == 0 {
			_, _ = db.DB.Exec("DELETE FROM matches WHERE match_id = ?", newMatchID)
		}
	}
	return rows.Err()
}

func (s *StandingsService) getGroupStandings(eventID int64) ([]model.GroupStandings, error) {
	rows, err := db.DB.Query(
		"SELECT m.round, m.team_a_id, m.team_b_id, m.score_team_a, m.score_team_b, m.status FROM matches m LEFT JOIN knockout_matches km ON km.match_id = m.match_id WHERE m.event_id = ? AND km.match_id IS NULL",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupStats := map[string]map[int64]*model.Standings{}
	teamIDs := map[int64]struct{}{}
	teamGroup := map[int64]string{}

	getStat := func(group string, teamID int64) *model.Standings {
		m, ok := groupStats[group]
		if !ok {
			m = map[int64]*model.Standings{}
			groupStats[group] = m
		}
		st, ok := m[teamID]
		if !ok {
			st = &model.Standings{EventID: eventID, TeamID: teamID}
			m[teamID] = st
		}
		return st
	}

	teamRows, err := db.DB.Query("SELECT team_id, COALESCE(group_name, '') FROM event_teams WHERE event_id = ?", eventID)
	if err != nil {
		return nil, err
	}
	for teamRows.Next() {
		var teamID int64
		var groupName string
		if err := teamRows.Scan(&teamID, &groupName); err != nil {
			_ = teamRows.Close()
			return nil, err
		}
		group := strings.TrimSpace(groupName)
		if group == "" {
			group = "小组赛"
		}
		teamGroup[teamID] = group
		_ = getStat(group, teamID)
		teamIDs[teamID] = struct{}{}
	}
	if err := teamRows.Err(); err != nil {
		_ = teamRows.Close()
		return nil, err
	}
	_ = teamRows.Close()

	for rows.Next() {
		var round string
		var teamA, teamB int64
		var scoreA, scoreB int
		var status string
		if err := rows.Scan(&round, &teamA, &teamB, &scoreA, &scoreB, &status); err != nil {
			return nil, err
		}
		group := strings.TrimSpace(round)
		if g, ok := teamGroup[teamA]; ok && strings.TrimSpace(g) != "" {
			group = g
		} else if g, ok := teamGroup[teamB]; ok && strings.TrimSpace(g) != "" {
			group = g
		}
		if strings.TrimSpace(group) == "" {
			group = "小组赛"
		}
		_ = getStat(group, teamA)
		_ = getStat(group, teamB)
		teamIDs[teamA] = struct{}{}
		teamIDs[teamB] = struct{}{}
		if status != "finished" {
			continue
		}
		statA := getStat(group, teamA)
		statB := getStat(group, teamB)

		statA.Played++
		statB.Played++
		statA.GoalsFor += scoreA
		statA.GoalsAgainst += scoreB
		statB.GoalsFor += scoreB
		statB.GoalsAgainst += scoreA

		if scoreA > scoreB {
			statA.Won++
			statA.Points += 3
			statB.Lost++
		} else if scoreA < scoreB {
			statB.Won++
			statB.Points += 3
			statA.Lost++
		} else {
			statA.Drawn++
			statA.Points += 1
			statB.Drawn++
			statB.Points += 1
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	teamInfo, err := fetchTeamInfo(teamIDs)
	if err != nil {
		return nil, err
	}

	var groups []model.GroupStandings
	for group, statsMap := range groupStats {
		var standings []*model.Standings
		for _, st := range statsMap {
			if info, ok := teamInfo[st.TeamID]; ok {
				st.TeamName = info.name
				st.TeamLogo = info.logo
			}
			standings = append(standings, st)
		}
		standings = sortStandings(standings)
		groups = append(groups, model.GroupStandings{Group: group, Standings: standings})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Group < groups[j].Group
	})
	return groups, nil
}

func (s *StandingsService) getKnockoutStages(eventID int64) ([]model.KnockoutStageView, error) {
	rows, err := db.DB.Query("SELECT stage_id, stage_name, stage_order FROM knockout_stages WHERE event_id = ? ORDER BY stage_order ASC", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []model.KnockoutStageView
	for rows.Next() {
		var st model.KnockoutStageView
		if err := rows.Scan(&st.StageID, &st.StageName, &st.StageOrder); err != nil {
			return nil, err
		}
		ms, err := fetchKnockoutMatchesByStage(st.StageID)
		if err != nil {
			return nil, err
		}
		st.Matches = ms
		stages = append(stages, st)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stages, nil
}

type teamBrief struct {
	name string
	logo string
}

func fetchTeamInfo(ids map[int64]struct{}) (map[int64]teamBrief, error) {
	if len(ids) == 0 {
		return map[int64]teamBrief{}, nil
	}
	var placeholders []string
	var args []interface{}
	for id := range ids {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	q := fmt.Sprintf("SELECT team_id, COALESCE(team_name, ''), COALESCE(avatar_url, '') FROM teams WHERE team_id IN (%s)", strings.Join(placeholders, ","))
	rows, err := db.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := map[int64]teamBrief{}
	for rows.Next() {
		var id int64
		var name, logo string
		if err := rows.Scan(&id, &name, &logo); err != nil {
			return nil, err
		}
		out[id] = teamBrief{name: name, logo: logo}
	}
	return out, rows.Err()
}

func sortStandings(standings []*model.Standings) []*model.Standings {
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		gdI := standings[i].GoalsFor - standings[i].GoalsAgainst
		gdJ := standings[j].GoalsFor - standings[j].GoalsAgainst
		if gdI != gdJ {
			return gdI > gdJ
		}
		return standings[i].GoalsFor > standings[j].GoalsFor
	})
	for i := range standings {
		standings[i].Rank = i + 1
	}
	return standings
}

func fetchKnockoutMatchesByStage(stageID int64) ([]model.KnockoutMatchView, error) {
	query := `
		SELECT km.knockout_match_id, km.stage_id, COALESCE(km.match_id, 0),
		       COALESCE(m.match_name, ''), COALESCE(m.round, ''), m.match_time, COALESCE(m.status, ''),
		       COALESCE(m.score_team_a, 0), COALESCE(m.score_team_b, 0),
		       COALESCE(m.team_a_id, 0), COALESCE(ta.team_name, ''), COALESCE(ta.avatar_url, ''),
		       COALESCE(m.team_b_id, 0), COALESCE(tb.team_name, ''), COALESCE(tb.avatar_url, ''),
		       COALESCE(km.prev_match_a_id, 0), COALESCE(km.prev_match_b_id, 0),
		       km.scheduled_time
		FROM knockout_matches km
		LEFT JOIN matches m ON km.match_id = m.match_id
		LEFT JOIN teams ta ON m.team_a_id = ta.team_id
		LEFT JOIN teams tb ON m.team_b_id = tb.team_id
		WHERE km.stage_id = ?
		ORDER BY COALESCE(km.match_order, 0) ASC, m.match_time ASC
	`
	rows, err := db.DB.Query(query, stageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.KnockoutMatchView
	for rows.Next() {
		var m model.KnockoutMatchView
		var matchTime sql.NullTime
		var scheduledTime sql.NullString
		if err := rows.Scan(
			&m.KnockoutMatchID, &m.StageID, &m.MatchID,
			&m.MatchName, &m.Round, &matchTime, &m.Status,
			&m.ScoreA, &m.ScoreB,
			&m.TeamAID, &m.TeamAName, &m.TeamALogo,
			&m.TeamBID, &m.TeamBName, &m.TeamBLogo,
			&m.PrevMatchAID, &m.PrevMatchBID,
			&scheduledTime,
		); err != nil {
			return nil, err
		}
		if matchTime.Valid {
			m.MatchTime = matchTime.Time.Format("2006-01-02 15:04")
		} else if scheduledTime.Valid {
			// If match not created yet, show scheduled time
			m.MatchTime = scheduledTime.String
			// Format if needed, but assuming DB string is OK or needs parsing
			// Usually DATETIME comes as string in some drivers or time.Time.
			// Here Scan to string is safe for DATETIME
			// Try to parse to pretty format if possible
			if t, err := time.Parse(time.RFC3339, scheduledTime.String); err == nil {
				m.MatchTime = t.Format("2006-01-02 15:04")
			} else if t, err := time.Parse("2006-01-02 15:04:05", scheduledTime.String); err == nil {
				m.MatchTime = t.Format("2006-01-02 15:04")
			}
			m.ScheduledTime = scheduledTime.String
		}
		if m.Status == "finished" && m.ScoreA != m.ScoreB {
			if m.ScoreA > m.ScoreB {
				m.WinnerTeamID = m.TeamAID
			} else {
				m.WinnerTeamID = m.TeamBID
			}
		}
		matches = append(matches, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return matches, nil
}
