package repo

import (
	"database/sql"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

// CreateUser 创建用户
func (r *UserRepo) CreateUser(user *model.User) error {
	query := `INSERT INTO users (student_id, password, name, college, grade, avatar_url) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.StudentID, user.Password, user.Name, user.College, user.Grade, user.AvatarURL)
	return err
}

// UpdateUserProfile 更新用户基础资料
func (r *UserRepo) UpdateUserProfile(studentID string, req *model.UpdateUserRequest) error {
	query := `UPDATE users 
	          SET name = ?, college = ?, grade = ?, avatar_url = ?
	          WHERE student_id = ?`
	_, err := db.Exec(query, req.Name, req.College, req.Grade, req.AvatarURL, studentID)
	return err
}

// GetUserByStudentID 根据学号获取用户
func (r *UserRepo) GetUserByStudentID(studentID string) (*model.User, error) {
	query := `SELECT student_id, password, name, college, grade, COALESCE(avatar_url, ''), created_at 
	          FROM users WHERE student_id = ?`

	var user model.User
	err := db.QueryRow(query, studentID).Scan(
		&user.StudentID, &user.Password, &user.Name,
		&user.College, &user.Grade, &user.AvatarURL, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// UserExists 检查用户是否存在
func (r *UserRepo) UserExists(studentID string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE student_id = ?`

	var count int
	err := db.QueryRow(query, studentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// UpdatePassword 更新用户密码
func (r *UserRepo) UpdatePassword(studentID string, newHashedPassword string) error {
	query := `UPDATE users SET password = ? WHERE student_id = ?`
	_, err := db.Exec(query, newHashedPassword, studentID)
	return err
}

// CreateTeam 创建队伍
func (r *UserRepo) CreateTeam(team *model.Team) (int64, error) {
	query := `INSERT INTO teams (team_name, sport_id, college, team_type, avatar_url, description, created_by) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query, team.TeamName, team.SportID, team.College, team.TeamType, team.AvatarURL, team.Description, team.CreatedBy)
}

// GetTeamByID 获取队伍信息
func (r *UserRepo) GetTeamByID(teamID int64) (*model.Team, error) {
	query := `SELECT team_id, team_name, sport_id, college, team_type, COALESCE(avatar_url, ''), description, created_by, created_at 
	          FROM teams WHERE team_id = ?`
	var team model.Team
	err := db.QueryRow(query, teamID).Scan(
		&team.ID, &team.TeamName, &team.SportID, &team.College,
		&team.TeamType, &team.AvatarURL, &team.Description, &team.CreatedBy, &team.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &team, err
}

// CreateAthlete 创建运动员档案
func (r *UserRepo) CreateAthlete(athlete *model.Athlete) (int64, error) {
	query := `INSERT INTO athletes (student_id, sport_type, team_id, jersey_number, is_captain) 
	          VALUES (?, ?, ?, ?, ?)`
	return db.Insert(query, athlete.StudentID, athlete.SportType, athlete.TeamID, athlete.JerseyNumber, athlete.IsCaptain)
}

// ListAthletesByStudentID 查询用户的运动员身份
func (r *UserRepo) ListAthletesByStudentID(studentID string) ([]model.Athlete, error) {
	query := `SELECT a.athlete_id, a.student_id, a.sport_type, a.team_id, 
	                 a.jersey_number, a.is_captain, t.team_name
	          FROM athletes a
	          LEFT JOIN teams t ON a.team_id = t.team_id
	          WHERE a.student_id = ?`
	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var athletes []model.Athlete
	for rows.Next() {
		var a model.Athlete
		if err := rows.Scan(&a.ID, &a.StudentID, &a.SportType, &a.TeamID, &a.JerseyNumber, &a.IsCaptain, &a.TeamName); err != nil {
			return nil, err
		}
		athletes = append(athletes, a)
	}
	return athletes, rows.Err()
}

// CreateTeamMember 绑定运动员与队伍
func (r *UserRepo) CreateTeamMember(teamID, athleteID int64, isActive, isApproved bool) (int64, error) {
	query := `INSERT INTO team_members (team_id, athlete_id, join_date, is_active, is_approved)
	          VALUES (?, ?, CURDATE(), ?, ?)`
	return db.Insert(query, teamID, athleteID, isActive, isApproved)
}

// ListTeamMembers 查询队伍成员
func (r *UserRepo) ListTeamMembers(teamID int64) ([]model.TeamMember, error) {
	query := `SELECT team_member_id, team_id, athlete_id, join_date, is_active, is_approved
	          FROM team_members WHERE team_id = ?`
	rows, err := db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.TeamMember
	for rows.Next() {
		var m model.TeamMember
		if err := rows.Scan(&m.ID, &m.TeamID, &m.AthleteID, &m.JoinDate, &m.IsActive, &m.IsApproved); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// ListPendingAthletes 获取待审核的运动员申请
func (r *UserRepo) ListPendingAthletes() ([]model.PendingAthleteApplication, error) {
	query := `
		SELECT tm.team_member_id, a.student_id, u.name, t.team_id, t.team_name, a.sport_type, tm.join_date
		FROM team_members tm
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		JOIN users u ON a.student_id = u.student_id
		JOIN teams t ON tm.team_id = t.team_id
		WHERE tm.is_approved = FALSE
		ORDER BY tm.join_date DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []model.PendingAthleteApplication
	for rows.Next() {
		var app model.PendingAthleteApplication
		if err := rows.Scan(&app.TeamMemberID, &app.StudentID, &app.Name, &app.TeamID, &app.TeamName, &app.SportType, &app.ApplyTime); err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

// IsApprovedAthlete 检查是否为已批准的运动员
func (r *UserRepo) IsApprovedAthlete(studentID string) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM team_members tm
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		WHERE a.student_id = ? AND tm.is_approved = TRUE AND tm.is_active = TRUE
	`
	var count int
	err := db.QueryRow(query, studentID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ApproveTeamMember 批准运动员加入队伍
func (r *UserRepo) ApproveTeamMember(teamMemberID int64) error {
	query := `UPDATE team_members SET is_approved = TRUE, is_active = TRUE WHERE team_member_id = ?`
	_, err := db.Exec(query, teamMemberID)
	return err
}

// CreateCollector 注册采集员
func (r *UserRepo) CreateCollector(studentID string, isApproved bool) (int64, error) {
	query := `INSERT INTO collectors (student_id, is_approved) VALUES (?, ?)`
	return db.Insert(query, studentID, isApproved)
}

// GetCollectorByStudentID 查询采集员信息
func (r *UserRepo) GetCollectorByStudentID(studentID string) (*model.Collector, error) {
	query := `SELECT collector_id, student_id, is_approved, created_at FROM collectors WHERE student_id = ?`
	var collector model.Collector
	err := db.QueryRow(query, studentID).Scan(&collector.ID, &collector.StudentID, &collector.IsApproved, &collector.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &collector, err
}

// CreateAdmin 创建管理员账号
func (r *UserRepo) CreateAdmin(admin *model.Admin) (int64, error) {
	query := `INSERT INTO admins (student_id, role, permissions) VALUES (?, ?, ?)`
	return db.Insert(query, admin.StudentID, admin.Role, admin.Permissions)
}

// GetAdminByStudentID 查询管理员信息
func (r *UserRepo) GetAdminByStudentID(studentID string) (*model.Admin, error) {
	query := `SELECT admin_id, student_id, role, COALESCE(permissions, '') FROM admins WHERE student_id = ?`
	var admin model.Admin
	err := db.QueryRow(query, studentID).Scan(&admin.ID, &admin.StudentID, &admin.Role, &admin.Permissions)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &admin, err
}

// ListPendingCollectors 获取所有待审核的采集员申请
func (r *UserRepo) ListPendingCollectors() ([]model.CollectorApplication, error) {
	query := `
		SELECT c.student_id, u.name, u.college, c.created_at
		FROM collectors c
		JOIN users u ON c.student_id = u.student_id
		WHERE c.is_approved = FALSE
		ORDER BY c.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []model.CollectorApplication
	for rows.Next() {
		var app model.CollectorApplication
		if err := rows.Scan(&app.StudentID, &app.Name, &app.College, &app.ApplyTime); err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}
	return applications, nil
}

// ApproveCollector 批准采集员申请
func (r *UserRepo) ApproveCollector(studentID string) error {
	query := `UPDATE collectors SET is_approved = TRUE WHERE student_id = ?`
	_, err := db.Exec(query, studentID)
	return err
}

// ListCollectors 获取所有已批准的采集员
func (r *UserRepo) ListCollectors() ([]model.CollectorApplication, error) {
	query := `
		SELECT c.student_id, u.name, u.college, c.created_at
		FROM collectors c
		JOIN users u ON c.student_id = u.student_id
		WHERE c.is_approved = TRUE
		ORDER BY c.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collectors []model.CollectorApplication
	for rows.Next() {
		var c model.CollectorApplication
		if err := rows.Scan(&c.StudentID, &c.Name, &c.College, &c.ApplyTime); err != nil {
			return nil, err
		}
		collectors = append(collectors, c)
	}
	return collectors, nil
}

// DeleteCollector 删除采集员
func (r *UserRepo) DeleteCollector(studentID string) error {
	query := `DELETE FROM collectors WHERE student_id = ?`
	_, err := db.Exec(query, studentID)
	return err
}
