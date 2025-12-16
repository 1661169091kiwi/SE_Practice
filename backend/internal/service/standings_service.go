package service

import (
	"sort"

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
	stmt, err := tx.Prepare("INSERT INTO standings (event_id, team_id, matches_played, wins, draws, losses, goals_for, goals_against, goal_difference, points, `rank`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for i, stat := range standings {
		stat.Rank = i + 1
		goalDiff := stat.GoalsFor - stat.GoalsAgainst
		if _, err := stmt.Exec(stat.EventID, stat.TeamID, stat.Played, stat.Won, stat.Drawn, stat.Lost, stat.GoalsFor, stat.GoalsAgainst, goalDiff, stat.Points, stat.Rank); err != nil {
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
		SELECT s.event_id, s.team_id, COALESCE(t.team_name, 'Unknown'), COALESCE(t.avatar_url, ''), s.matches_played, s.wins, s.draws, s.losses, s.goals_for, s.goals_against, s.points, s.rank
		FROM standings s
		LEFT JOIN teams t ON s.team_id = t.team_id
		WHERE s.event_id = ?
		ORDER BY s.rank ASC
	`
	rows, err := db.DB.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Standings
	for rows.Next() {
		var st model.Standings
		if err := rows.Scan(&st.EventID, &st.TeamID, &st.TeamName, &st.TeamLogo, &st.Played, &st.Won, &st.Drawn, &st.Lost, &st.GoalsFor, &st.GoalsAgainst, &st.Points, &st.Rank); err != nil {
			return nil, err
		}
		result = append(result, &st)
	}
	return result, nil
}
