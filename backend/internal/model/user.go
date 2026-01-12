package model

import "time"

// User 用户数据模型
type User struct {
	StudentID   string    `json:"student_id"`
	Password    string    `json:"-"` // 密码不返回给前端
	Name        string    `json:"name"`
	College     string    `json:"college"`
	Grade       string    `json:"grade"`
	Role        string    `json:"role"`
	Roles       []string  `json:"roles"`
	AvatarURL   string    `json:"avatar_url"`
	ApplyStatus string    `json:"apply_status,omitempty"` // none, pending, approved
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	College   string `json:"college"`
	Grade     string `json:"grade"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name      string `json:"name"`
	College   string `json:"college"`
	Grade     string `json:"grade"`
	AvatarURL string `json:"avatar_url"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	StudentID   string `json:"student_id" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// UpdateAvatarResponse 更新头像响应
type UpdateAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
	Message   string `json:"message"`
}

// UpdateCollegeRequest 更新学院请求
type UpdateCollegeRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	College   string `json:"college" binding:"required"`
}

// UpdateCollegeResponse 更新学院响应
type UpdateCollegeResponse struct {
	StudentID string `json:"student_id"`
	College   string `json:"college"`
	Message   string `json:"message"`
}

// UpdateNameRequest 更新姓名请求
type UpdateNameRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
}

// UpdateNameResponse 更新姓名响应
type UpdateNameResponse struct {
	StudentID string `json:"student_id"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}

// ApplyAthleteRequest 申请成为运动员请求
type ApplyAthleteRequest struct {
	StudentID    string `json:"student_id" binding:"required"`
	SportType    string `json:"sport_type" binding:"required"`
	TeamID       int64  `json:"team_id" binding:"required"`
	JerseyNumber string `json:"jersey_number"`
	IsCaptain    bool   `json:"is_captain"`
}

// Athlete 运动员信息
type Athlete struct {
	ID           int64  `json:"athlete_id"`
	StudentID    string `json:"student_id"`
	SportType    string `json:"sport_type"`
	TeamID       int64  `json:"team_id"`
	TeamName     string `json:"team_name,omitempty"`
	JerseyNumber string `json:"jersey_number"`
	IsCaptain    bool   `json:"is_captain"`
}

// TeamMember 队伍成员信息
type TeamMember struct {
	ID         int64     `json:"team_member_id"`
	TeamID     int64     `json:"team_id"`
	AthleteID  int64     `json:"athlete_id"`
	JoinDate   time.Time `json:"join_date"`
	IsActive   bool      `json:"is_active"`
	IsApproved bool      `json:"is_approved"`
}

// PendingAthleteApplication 待审核运动员申请
type PendingAthleteApplication struct {
	TeamMemberID int64     `json:"team_member_id"`
	StudentID    string    `json:"student_id"`
	Name         string    `json:"name"`
	TeamID       int64     `json:"team_id"`
	TeamName     string    `json:"team_name"`
	SportType    string    `json:"sport_type"`
	ApplyTime    time.Time `json:"apply_time"`
}

// Collector 采集员
type Collector struct {
	ID         int64     `json:"collector_id"`
	StudentID  string    `json:"student_id"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
}

// Admin 管理员
type Admin struct {
	ID          int64  `json:"admin_id"`
	StudentID   string `json:"student_id"`
	Role        string `json:"role"`
	Permissions string `json:"permissions"`
}

// CollectorApplication 采集员申请信息（用于管理员查看）
type CollectorApplication struct {
	StudentID string    `json:"student_id"`
	Name      string    `json:"name"`
	College   string    `json:"college"`
	ApplyTime time.Time `json:"apply_time"`
}
