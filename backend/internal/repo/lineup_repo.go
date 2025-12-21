package repo

import (
	"database/sql"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type LineupRepo struct{}

func NewLineupRepo() *LineupRepo {
	return &LineupRepo{}
}

func (r *LineupRepo) EnsureExternalTable() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS match_lineups_external (
ext_id INT AUTO_INCREMENT PRIMARY KEY,
match_id INT NOT NULL,
team_id INT NOT NULL,
player_name VARCHAR(50) NOT NULL,
position VARCHAR(50),
is_starting BOOLEAN DEFAULT TRUE,
jersey_number VARCHAR(10),
FOREIGN KEY (match_id) REFERENCES matches(match_id),
FOREIGN KEY (team_id) REFERENCES teams(team_id)
)`)
	return err
}

func (r *LineupRepo) Upsert(item model.MatchLineup) error {
	res, err := db.Update(`UPDATE match_lineups SET position=?, is_starting=?, jersey_number=?
WHERE match_id=? AND team_id=? AND student_id=?`,
		item.Position, item.IsStarting, item.JerseyNumber, item.MatchID, item.TeamID, item.StudentID)
	if err != nil {
		return err
	}
	if res == 0 {
		_, err = db.Exec(`INSERT INTO match_lineups (match_id, team_id, student_id, position, is_starting, jersey_number)
VALUES (?, ?, ?, ?, ?, ?)`,
			item.MatchID, item.TeamID, item.StudentID, item.Position, item.IsStarting, item.JerseyNumber)
		return err
	}
	return nil
}

func (r *LineupRepo) UpsertExternal(matchID int64, teamID int64, name string, position string, jersey string, isStarting bool) error {
	if err := r.EnsureExternalTable(); err != nil {
		return err
	}
	res, err := db.Update(`UPDATE match_lineups_external SET position=?, is_starting=?, jersey_number=?
WHERE match_id=? AND team_id=? AND player_name=?`,
		position, isStarting, jersey, matchID, teamID, name)
	if err != nil {
		return err
	}
	if res == 0 {
		_, err = db.Exec(`INSERT INTO match_lineups_external (match_id, team_id, player_name, position, is_starting, jersey_number)
VALUES (?, ?, ?, ?, ?, ?)`, matchID, teamID, name, position, isStarting, jersey)
		return err
	}
	return nil
}

func (r *LineupRepo) ListByMatch(matchID int64) ([]model.MatchLineup, error) {
	rows, err := db.Query(`SELECT lineup_id, match_id, team_id, student_id, position, is_starting, jersey_number
FROM match_lineups WHERE match_id = ? ORDER BY team_id ASC, is_starting DESC`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []model.MatchLineup
	for rows.Next() {
		var ml model.MatchLineup
		if err := rows.Scan(&ml.LineupID, &ml.MatchID, &ml.TeamID, &ml.StudentID, &ml.Position, &ml.IsStarting, &ml.JerseyNumber); err != nil {
			return nil, err
		}
		list = append(list, ml)
	}
	return list, rows.Err()
}

func (r *LineupRepo) UpsertTx(tx *sql.Tx, item model.MatchLineup) error {
	res, err := tx.Exec(`UPDATE match_lineups SET position=?, is_starting=?, jersey_number=? WHERE match_id=? AND team_id=? AND student_id=?`,
		item.Position, item.IsStarting, item.JerseyNumber, item.MatchID, item.TeamID, item.StudentID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		_, err = tx.Exec(`INSERT INTO match_lineups (match_id, team_id, student_id, position, is_starting, jersey_number) VALUES (?, ?, ?, ?, ?, ?)`,
			item.MatchID, item.TeamID, item.StudentID, item.Position, item.IsStarting, item.JerseyNumber)
		return err
	}
	return nil
}

func (r *LineupRepo) UpsertExternalTx(tx *sql.Tx, matchID int64, teamID int64, name string, position string, jersey string, isStarting bool) error {
	res, err := tx.Exec(`UPDATE match_lineups_external SET position=?, is_starting=?, jersey_number=? WHERE match_id=? AND team_id=? AND player_name=? AND COALESCE(position,'') = COALESCE(?, '') AND COALESCE(jersey_number,'') = COALESCE(?, '')`,
		position, isStarting, jersey, matchID, teamID, name, position, jersey)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		_, err = tx.Exec(`INSERT INTO match_lineups_external (match_id, team_id, player_name, position, is_starting, jersey_number) VALUES (?, ?, ?, ?, ?, ?)`,
			matchID, teamID, name, position, isStarting, jersey)
		return err
	}
	return nil
}

func (r *LineupRepo) DeleteTx(tx *sql.Tx, matchID int64, teamID int64, studentID string) error {
	_, err := tx.Exec(`DELETE FROM match_lineups WHERE match_id=? AND team_id=? AND student_id=?`, matchID, teamID, studentID)
	return err
}

func (r *LineupRepo) DeleteExternalTx(tx *sql.Tx, matchID int64, teamID int64, name string, position string, jersey string) error {
	_, err := tx.Exec(`DELETE FROM match_lineups_external WHERE match_id=? AND team_id=? AND player_name=? AND COALESCE(position,'') = COALESCE(?, '') AND COALESCE(jersey_number,'') = COALESCE(?, '')`,
		matchID, teamID, name, position, jersey)
	return err
}
