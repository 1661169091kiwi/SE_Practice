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

// GetUserByStudentID 根据学号获取用户
func (r *UserRepo) GetUserByStudentID(studentID string) (*model.User, error) {
	query := `SELECT student_id, password, name, college, grade, avatar_url, created_at 
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


