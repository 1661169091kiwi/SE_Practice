package repo

import (
	"database/sql"
	"fmt"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"strings"
	"time"
)

type EventRepo struct{}

func NewEventRepo() *EventRepo {
	return &EventRepo{}
}

// CreateEvent 创建赛事
func (r *EventRepo) CreateEvent(event *model.Event) (int64, error) {
	query := `INSERT INTO events 
	          (event_name, sport_id, season, round, format_type, start_date, end_date, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query,
		event.Name,
		event.SportID,
		event.Season,
		event.Round,
		event.Format,
		event.StartDate,
		event.EndDate,
		event.Status,
	)
}

// GetEventByID 根据ID获取赛事
func (r *EventRepo) GetEventByID(id int64) (*model.Event, error) {
	query := `SELECT event_id, event_name, sport_id, COALESCE(season, ''), COALESCE(round, ''), COALESCE(format_type, ''), 
	                 start_date, end_date, status, created_at 
	          FROM events WHERE event_id = ?`

	var event model.Event
	var endDate sql.NullTime

	err := db.QueryRow(query, id).Scan(
		&event.ID,
		&event.Name,
		&event.SportID,
		&event.Season,
		&event.Round,
		&event.Format,
		&event.StartDate,
		&endDate,
		&event.Status,
		&event.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if endDate.Valid {
		event.EndDate = endDate.Time
	}
	return &event, nil
}

// ListEvents 列出赛事
func (r *EventRepo) ListEvents(status string) ([]model.Event, error) {
	query := `SELECT event_id, event_name, sport_id, COALESCE(season, ''), COALESCE(round, ''), COALESCE(format_type, ''), 
	                 start_date, end_date, status, created_at 
	          FROM events`

	var args []interface{}
	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}
	query += " ORDER BY start_date DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var e model.Event
		var endDate sql.NullTime
		if err := rows.Scan(
			&e.ID, &e.Name, &e.SportID, &e.Season, &e.Round,
			&e.Format, &e.StartDate, &endDate, &e.Status, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		if endDate.Valid {
			e.EndDate = endDate.Time
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// DeleteEvent 删除赛事
func (r *EventRepo) DeleteEvent(eventID int64) error {
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}

	// 1. Delete subscriptions
	if _, err := tx.Exec("DELETE FROM subscriptions WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 2. Delete dependent data for matches
	// comment_likes (via player_ratings)
	if _, err := tx.Exec("DELETE cl FROM comment_likes cl JOIN player_ratings pr ON cl.rating_id = pr.rating_id JOIN matches m ON pr.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// player_ratings
	if _, err := tx.Exec("DELETE pr FROM player_ratings pr JOIN matches m ON pr.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// match_comments
	if _, err := tx.Exec("DELETE mc FROM match_comments mc JOIN matches m ON mc.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// match_lineups
	if _, err := tx.Exec("DELETE ml FROM match_lineups ml JOIN matches m ON ml.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// match_lineups_external
	if _, err := tx.Exec("DELETE mle FROM match_lineups_external mle JOIN matches m ON mle.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		// Ignore if table doesn't exist, but it should
	}
	// match_events
	if _, err := tx.Exec("DELETE me FROM match_events me JOIN matches m ON me.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// injury_suspensions
	if _, err := tx.Exec("DELETE ins FROM injury_suspensions ins JOIN matches m ON ins.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// match_stats
	if _, err := tx.Exec("DELETE ms FROM match_stats ms JOIN matches m ON ms.match_id = m.match_id WHERE m.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 3. Delete knockout data
	// knockout_matches (both instantiated and virtual)
	if _, err := tx.Exec("DELETE km FROM knockout_matches km JOIN knockout_stages ks ON km.stage_id = ks.stage_id WHERE ks.event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}
	// knockout_stages
	if _, err := tx.Exec("DELETE FROM knockout_stages WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 4. Delete matches
	if _, err := tx.Exec("DELETE FROM matches WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 5. Delete event teams
	if _, err := tx.Exec("DELETE FROM event_teams WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 6. Delete standings
	if _, err := tx.Exec("DELETE FROM standings WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 7. Delete event
	if _, err := tx.Exec("DELETE FROM events WHERE event_id = ?", eventID); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UpdateEventStatus 更新赛事状态
func (r *EventRepo) UpdateEventStatus(id int64, status string) error {
	query := "UPDATE events SET status = ? WHERE event_id = ?"
	_, err := db.Exec(query, status, id)
	return err
}

// UpdateEvent 更新赛事信息
func (r *EventRepo) UpdateEvent(event *model.Event) error {
	query := `UPDATE events SET 
              event_name=?, 
              sport_id=?, 
              season=?, 
              round=?, 
              format_type=?, 
              start_date=?, 
              end_date=? 
              WHERE event_id=?`
	_, err := db.Exec(query,
		event.Name,
		event.SportID,
		event.Season,
		event.Round,
		event.Format,
		event.StartDate,
		event.EndDate,
		event.ID,
	)
	return err
}

func (r *EventRepo) ReplaceEventTeams(eventID int64, teams []model.EventTeamInput) error {
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM event_teams WHERE event_id = ?", eventID); err != nil {
		return err
	}

	if len(teams) == 0 {
		return tx.Commit()
	}

	stmt, err := tx.Prepare("INSERT INTO event_teams (event_id, team_id, group_name, slot) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, t := range teams {
		var group any
		if strings.TrimSpace(t.GroupName) == "" {
			group = nil
		} else {
			group = strings.TrimSpace(t.GroupName)
		}
		var slot any
		if t.Slot <= 0 {
			slot = nil
		} else {
			slot = t.Slot
		}
		if _, err := stmt.Exec(eventID, t.TeamID, group, slot); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *EventRepo) ListEventTeams(eventID int64) ([]model.EventTeam, error) {
	query := `
		SELECT et.event_id, et.team_id, COALESCE(t.team_name, ''), COALESCE(et.group_name, ''), COALESCE(et.slot, 0)
		FROM event_teams et
		JOIN teams t ON et.team_id = t.team_id
		WHERE et.event_id = ?
		ORDER BY COALESCE(et.slot, 1000000) ASC, et.team_id ASC
	`
	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.EventTeam
	for rows.Next() {
		var it model.EventTeam
		if err := rows.Scan(&it.EventID, &it.TeamID, &it.TeamName, &it.GroupName, &it.Slot); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *EventRepo) IsTeamInEvent(eventID, teamID int64) (bool, error) {
	var c int
	err := db.QueryRow("SELECT COUNT(*) FROM event_teams WHERE event_id = ? AND team_id = ?", eventID, teamID).Scan(&c)
	if err != nil {
		return false, err
	}
	return c > 0, nil
}

func (r *EventRepo) CountEventTeams(eventID int64) (int, error) {
	var c int
	if err := db.QueryRow("SELECT COUNT(*) FROM event_teams WHERE event_id = ?", eventID).Scan(&c); err != nil {
		return 0, err
	}
	return c, nil
}

func (r *EventRepo) SetupKnockoutBracket(eventID int64, schedule map[string]string) error {
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var format string
	if err := tx.QueryRow("SELECT COALESCE(format_type, '') FROM events WHERE event_id = ?", eventID).Scan(&format); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("event not found")
		}
		return err
	}
	if format != "knockout" && format != "group_knockout" {
		return fmt.Errorf("invalid event format")
	}

	// cleanup old bracket (and its generated matches)
	rows, err := tx.Query(`
		SELECT km.match_id
		FROM knockout_matches km
		JOIN knockout_stages ks ON km.stage_id = ks.stage_id
		WHERE ks.event_id = ? AND km.match_id IS NOT NULL
	`, eventID)
	if err != nil {
		return err
	}
	var matchIDs []int64
	for rows.Next() {
		var mid int64
		if err := rows.Scan(&mid); err != nil {
			_ = rows.Close()
			return err
		}
		matchIDs = append(matchIDs, mid)
	}
	_ = rows.Close()

	if _, err := tx.Exec("DELETE km FROM knockout_matches km JOIN knockout_stages ks ON km.stage_id = ks.stage_id WHERE ks.event_id = ?", eventID); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM knockout_stages WHERE event_id = ?", eventID); err != nil {
		return err
	}
	if len(matchIDs) > 0 {
		var placeholders []string
		var args []any
		for _, id := range matchIDs {
			placeholders = append(placeholders, "?")
			args = append(args, id)
		}
		if _, err := tx.Exec("DELETE FROM matches WHERE match_id IN ("+strings.Join(placeholders, ",")+")", args...); err != nil {
			return err
		}
	}

	teams, err := r.listEventTeamsTx(tx, eventID)
	if err != nil {
		return err
	}
	if format == "group_knockout" {
		var picked []model.EventTeam
		for _, t := range teams {
			if t.Slot > 0 {
				picked = append(picked, t)
			}
		}
		teams = picked
	}
	if len(teams) < 2 {
		return fmt.Errorf("not enough teams")
	}

	// validate slots
	slotToTeam := map[int]int64{}
	maxSlot := 0
	for _, t := range teams {
		if t.Slot <= 0 {
			return fmt.Errorf("knockout slot required")
		}
		if _, ok := slotToTeam[t.Slot]; ok {
			return fmt.Errorf("duplicate slot")
		}
		slotToTeam[t.Slot] = t.TeamID
		if t.Slot > maxSlot {
			maxSlot = t.Slot
		}
	}
	if maxSlot != len(teams) {
		return fmt.Errorf("slots must be 1..N without gaps")
	}
	if (maxSlot & (maxSlot - 1)) != 0 {
		return fmt.Errorf("team count must be power of two")
	}

	roundCount := 0
	for n := maxSlot; n > 1; n >>= 1 {
		roundCount++
	}

	stageNames := make([]string, 0, roundCount)
	for n := maxSlot; n >= 2; n >>= 1 {
		stageNames = append(stageNames, knockoutStageName(n))
	}

	// create stages
	stageIDs := make([]int64, 0, len(stageNames))
	for i, name := range stageNames {
		res, err := tx.Exec("INSERT INTO knockout_stages (event_id, stage_name, stage_order) VALUES (?, ?, ?)", eventID, name, i+1)
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		stageIDs = append(stageIDs, id)
	}

	// round 1: create real matches and bind to knockout_matches
	type pair struct{ a, b int64 }
	var pairs []pair
	for i := 1; i <= maxSlot/2; i++ {
		ta := slotToTeam[i]
		tb := slotToTeam[maxSlot+1-i]
		pairs = append(pairs, pair{a: ta, b: tb})
	}

	prevMatchIDs := make([]int64, 0, len(pairs))
	for i, p := range pairs {
		matchName := fmt.Sprintf("%s %d", stageNames[0], i+1)

		key := fmt.Sprintf("r1-m%d", i+1)
		timeStr := schedule[key]
		var matchTime interface{} = time.Now()
		var scheduledTime interface{} = nil
		if timeStr != "" {
			if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
				matchTime = t
				scheduledTime = t
			}
		}

		res, err := tx.Exec(
			"INSERT INTO matches (event_id, match_name, round, match_time, team_a_id, team_b_id, status) VALUES (?, ?, ?, ?, ?, ?, 'not_started')",
			eventID, matchName, stageNames[0], matchTime, p.a, p.b,
		)
		if err != nil {
			return err
		}
		mid, _ := res.LastInsertId()

		res2, err := tx.Exec("INSERT INTO knockout_matches (stage_id, match_id, match_order, prev_match_a_id, prev_match_b_id, scheduled_time) VALUES (?, ?, ?, NULL, NULL, ?)", stageIDs[0], mid, i+1, scheduledTime)
		if err != nil {
			return err
		}
		kmid, _ := res2.LastInsertId()
		prevMatchIDs = append(prevMatchIDs, kmid)
	}

	// following rounds: create knockout_matches placeholders with prev links
	for round := 1; round < len(stageIDs); round++ {
		var nextPrev []int64
		for i := 0; i < len(prevMatchIDs); i += 2 {
			prevA := prevMatchIDs[i]
			prevB := prevMatchIDs[i+1]

			key := fmt.Sprintf("r%d-m%d", round+1, (i/2)+1)
			timeStr := schedule[key]
			var scheduledTime interface{} = nil
			if timeStr != "" {
				if t, err := time.Parse(time.RFC3339, timeStr); err == nil {
					scheduledTime = t
				}
			}

			res, err := tx.Exec("INSERT INTO knockout_matches (stage_id, match_id, match_order, prev_match_a_id, prev_match_b_id, scheduled_time) VALUES (?, NULL, ?, ?, ?, ?)", stageIDs[round], (i/2)+1, prevA, prevB, scheduledTime)
			if err != nil {
				return err
			}
			kmid, _ := res.LastInsertId()
			nextPrev = append(nextPrev, kmid)
		}
		prevMatchIDs = nextPrev
	}

	return tx.Commit()
}

// UpdateKnockoutSchedule 更新淘汰赛赛程
func (r *EventRepo) UpdateKnockoutSchedule(updates map[int64]string) error {
	tx, err := db.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE knockout_matches SET scheduled_time = ? WHERE knockout_match_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	updateMatchStmt, err := tx.Prepare("UPDATE matches SET match_time = ? WHERE match_id = ?")
	if err != nil {
		return err
	}
	defer updateMatchStmt.Close()

	for kmid, timeStr := range updates {
		parsedTime, err := time.Parse(time.RFC3339, timeStr)
		if err != nil {
			continue // Skip invalid time strings
		}
		if _, err := stmt.Exec(parsedTime, kmid); err != nil {
			return err
		}

		// Check if there is an existing match_id
		var matchID sql.NullInt64
		if err := tx.QueryRow("SELECT match_id FROM knockout_matches WHERE knockout_match_id = ?", kmid).Scan(&matchID); err != nil {
			continue
		}
		if matchID.Valid {
			if _, err := updateMatchStmt.Exec(parsedTime, matchID.Int64); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *EventRepo) listEventTeamsTx(tx *sql.Tx, eventID int64) ([]model.EventTeam, error) {
	rows, err := tx.Query("SELECT event_id, team_id, COALESCE(group_name, ''), COALESCE(slot, 0) FROM event_teams WHERE event_id = ?", eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.EventTeam
	for rows.Next() {
		var it model.EventTeam
		if err := rows.Scan(&it.EventID, &it.TeamID, &it.GroupName, &it.Slot); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func knockoutStageName(n int) string {
	switch n {
	case 2:
		return "决赛"
	case 4:
		return "半决赛"
	case 8:
		return "四分之一决赛"
	case 16:
		return "八分之一决赛"
	case 32:
		return "十六分之一决赛"
	default:
		return fmt.Sprintf("第%d轮", n)
	}
}
