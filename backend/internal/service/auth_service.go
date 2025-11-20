package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

// AuthService 负责登录/注册等认证相关业务
type AuthService struct {
	userRepo *repo.UserRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repo.NewUserRepo(),
	}
}

var (
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidCredential = errors.New("invalid student_id or password")
)

// 使用标准库做一个简单的哈希（教学/作业用，生产环境建议用 bcrypt / argon2 等）
func hashPassword(pw string) string {
	sum := sha256.Sum256([]byte(pw))
	return hex.EncodeToString(sum[:])
}

func checkPassword(hashed, pw string) bool {
	return hashed == hashPassword(pw)
}

// Register 用户注册：检查学号是否存在 -> 加密密码 -> 写入 users 表
func (s *AuthService) Register(req *model.RegisterRequest) (*model.User, error) {
	// 1. 检查是否已存在
	exists, err := s.userRepo.UserExists(req.StudentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserExists
	}

	// 2. 加密密码（这里用 sha256 简单哈希）
	hashed := hashPassword(req.Password)

	// 3. 组装 User 模型并写库
	user := &model.User{
		StudentID: req.StudentID,
		Password:  hashed,
		Name:      req.Name,
		College:   req.College,
		Grade:     req.Grade,
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// 4. 再查一遍带 created_at 返回
	dbUser, err := s.userRepo.GetUserByStudentID(req.StudentID)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

// Login 用户登录：校验密码，返回基础信息（后续可加 JWT）
func (s *AuthService) Login(req *model.LoginRequest) (*model.User, error) {
	user, err := s.userRepo.GetUserByStudentID(req.StudentID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredential
	}

	if !checkPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredential
	}

	return user, nil
}
