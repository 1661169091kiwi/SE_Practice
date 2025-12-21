package repo

import (
	"database/sql"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type TeamRepo struct{}

func NewTeamRepo() *TeamRepo {
	return &TeamRepo{}
}

// CreateTeam 创建队伍
func (r *TeamRepo) CreateTeam(team *model.Team) (int64, error) {
	query := `INSERT INTO teams (team_name, sport_id, college, team_type, avatar_url, description, created_by, is_approved) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query, team.TeamName, team.SportID, team.College, team.TeamType, team.AvatarURL, team.Description, team.CreatedBy, team.IsApproved)
}

// GetTeamByID 获取队伍信息
func (r *TeamRepo) GetTeamByID(teamID int64) (*model.Team, error) {
	query := `SELECT team_id, team_name, sport_id, college, team_type, COALESCE(avatar_url, ''), COALESCE(description, ''), created_by, created_at, is_approved 
	          FROM teams WHERE team_id = ?`
	var team model.Team
	err := db.QueryRow(query, teamID).Scan(
		&team.ID, &team.TeamName, &team.SportID, &team.College,
		&team.TeamType, &team.AvatarURL, &team.Description, &team.CreatedBy, &team.CreatedAt, &team.IsApproved,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &team, err
}

// ListTeams 列出所有已批准的队伍
func (r *TeamRepo) ListTeams(sportID int64) ([]model.Team, error) {
	query := `SELECT team_id, team_name, sport_id, college, team_type, COALESCE(avatar_url, ''), COALESCE(description, ''), created_by, created_at, is_approved 
	          FROM teams WHERE is_approved = TRUE`

	var args []interface{}
	if sportID > 0 {
		query += " AND sport_id = ?"
		args = append(args, sportID)
	}
	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var t model.Team
		if err := rows.Scan(
			&t.ID, &t.TeamName, &t.SportID, &t.College,
			&t.TeamType, &t.AvatarURL, &t.Description, &t.CreatedBy, &t.CreatedAt, &t.IsApproved,
		); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

// ListPendingTeams 获取所有待审核的队伍
func (r *TeamRepo) ListPendingTeams() ([]model.Team, error) {
	query := `SELECT team_id, team_name, sport_id, college, team_type, COALESCE(avatar_url, ''), COALESCE(description, ''), created_by, created_at, is_approved 
	          FROM teams WHERE is_approved = FALSE ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var t model.Team
		if err := rows.Scan(
			&t.ID, &t.TeamName, &t.SportID, &t.College,
			&t.TeamType, &t.AvatarURL, &t.Description, &t.CreatedBy, &t.CreatedAt, &t.IsApproved,
		); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, rows.Err()
}

// GetTeamsByStudentID 获取学生加入的队伍
func (r *TeamRepo) GetTeamsByStudentID(studentID string) ([]model.Team, error) {
	query := `
		SELECT t.team_id, t.team_name, t.sport_id, t.college, t.team_type, 
		       COALESCE(t.avatar_url, ''), COALESCE(t.description, ''), t.created_by, t.created_at, t.is_approved
		FROM teams t
		JOIN team_members tm ON t.team_id = tm.team_id
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		WHERE a.student_id = ? AND tm.is_approved = TRUE AND tm.is_active = TRUE
		
		UNION
		
		SELECT t.team_id, t.team_name, t.sport_id, t.college, t.team_type, 
		       COALESCE(t.avatar_url, ''), COALESCE(t.description, ''), t.created_by, t.created_at, t.is_approved
		FROM teams t
		WHERE t.created_by = ?
	`
	rows, err := db.Query(query, studentID, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []model.Team
	for rows.Next() {
		var team model.Team
		if err := rows.Scan(
			&team.ID, &team.TeamName, &team.SportID, &team.College,
			&team.TeamType, &team.AvatarURL, &team.Description, &team.CreatedBy, &team.CreatedAt, &team.IsApproved,
		); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

// ApproveTeam 批准队伍
func (r *TeamRepo) ApproveTeam(teamID int64) error {
	query := `UPDATE teams SET is_approved = TRUE WHERE team_id = ?`
	_, err := db.Exec(query, teamID)
	return err
}

// UpdateTeamName 更新队伍名称
func (r *TeamRepo) UpdateTeamName(teamID int64, name string) error {
	query := `UPDATE teams SET team_name = ? WHERE team_id = ?`
	_, err := db.Exec(query, name, teamID)
	return err
}

// DeleteTeam 删除队伍
func (r *TeamRepo) DeleteTeam(teamID int64) error {
	// First delete team members
	_, err := db.Exec("DELETE FROM team_members WHERE team_id = ?", teamID)
	if err != nil {
		return err
	}
	// Then delete team
	_, err = db.Exec("DELETE FROM teams WHERE team_id = ?", teamID)
	return err
}

// GetTeamMembers 获取队伍成员列表
func (r *TeamRepo) GetTeamMembers(teamID int64) ([]model.TeamMemberDetail, error) {
	query := `
		SELECT tm.team_member_id, tm.team_id, tm.athlete_id, 
		       u.student_id, u.name, u.college, a.sport_type, COALESCE(tm.join_date, ''), COALESCE(a.jersey_number, '')
		FROM team_members tm
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		JOIN users u ON a.student_id = u.student_id
		WHERE tm.team_id = ? AND tm.is_active = TRUE
	`
	rows, err := db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMemberDetail
	for rows.Next() {
		var m model.TeamMemberDetail
		if err := rows.Scan(
			&m.TeamMemberID, &m.TeamID, &m.AthleteID,
			&m.StudentID, &m.Name, &m.College, &m.SportType, &m.JoinDate, &m.JerseyNumber,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// GetPendingTeamMembers 获取待审核队伍成员
func (r *TeamRepo) GetPendingTeamMembers(teamID int64) ([]model.TeamMemberDetail, error) {
	query := `
		SELECT tm.team_member_id, tm.team_id, tm.athlete_id, 
		       u.student_id, u.name, u.college, a.sport_type, COALESCE(tm.join_date, ''), COALESCE(a.jersey_number, '')
		FROM team_members tm
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		JOIN users u ON a.student_id = u.student_id
		WHERE tm.team_id = ? AND tm.is_approved = FALSE
	`
	rows, err := db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMemberDetail
	for rows.Next() {
		var m model.TeamMemberDetail
		if err := rows.Scan(
			&m.TeamMemberID, &m.TeamID, &m.AthleteID,
			&m.StudentID, &m.Name, &m.College, &m.SportType, &m.JoinDate, &m.JerseyNumber,
		); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// ApproveMember 批准队员
func (r *TeamRepo) ApproveMember(teamMemberID int64) error {
	query := `UPDATE team_members SET is_approved = TRUE, is_active = TRUE WHERE team_member_id = ?`
	_, err := db.Exec(query, teamMemberID)
	return err
}

// RemoveMember 移除队伍成员
func (r *TeamRepo) RemoveMember(teamMemberID int64) error {
	_, err := db.Exec("DELETE FROM team_members WHERE team_member_id = ?", teamMemberID)
	return err
}

// GetTeamMemberByID 根据ID获取成员信息
func (r *TeamRepo) GetTeamMemberByID(memberID int64) (*model.TeamMemberDetail, error) {
	query := `SELECT team_member_id, team_id, athlete_id, '' as student_id, '' as name, '' as college, '' as sport_type, COALESCE(join_date, '')
	          FROM team_members WHERE team_member_id = ?`
	// Note: We only need TeamID for permission check, so simple query is enough.
	// If we need full details, we would join.
	var m model.TeamMemberDetail
	err := db.QueryRow(query, memberID).Scan(
		&m.TeamMemberID, &m.TeamID, &m.AthleteID, &m.StudentID, &m.Name, &m.College, &m.SportType, &m.JoinDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}
