package model

import "time"

// User 用户数据模型
type User struct {
	StudentID string    `json:"student_id"`
	Password  string    `json:"-"` // 密码不返回给前端
	Name      string    `json:"name"`
	College   string    `json:"college"`
	Grade     string    `json:"grade"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
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