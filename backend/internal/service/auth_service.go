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
	ErrPasswordSame      = errors.New("新密码不能与旧密码相同")
	ErrFileTooLarge      = errors.New("file too large")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrUserNotFound      = errors.New("user not found")
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

// GetUserProfile 获取用户个人资料
func (s *AuthService) GetUserProfile(studentID string) (*model.User, error) {
	// 从数据库获取用户信息
	user, err := s.userRepo.GetUserByStudentID(studentID)
	if err != nil {
		return nil, err
	}

	// 如果用户不存在，直接返回nil
	if user == nil {
		return nil, nil
	}

	// 处理默认头像：如果用户未设置头像，使用默认头像链接
	if user.AvatarURL == "" {
		user.AvatarURL = "https://via.placeholder.com/150"
	}

	return user, nil
}

// UpdateAvatar 更新用户头像
func (s *AuthService) UpdateAvatar(studentID string, avatarURL string) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetUserByStudentID(studentID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 创建更新请求
	updateReq := &model.UpdateUserRequest{
		Name:      user.Name,
		College:   user.College,
		Grade:     user.Grade,
		AvatarURL: avatarURL,
	}

	// 更新用户头像
	return s.userRepo.UpdateUserProfile(studentID, updateReq)
}

// UpdateCollege 更新用户学院信息
func (s *AuthService) UpdateCollege(req *model.UpdateCollegeRequest) (*model.UpdateCollegeResponse, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetUserByStudentID(req.StudentID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 创建更新请求
	updateReq := &model.UpdateUserRequest{
		Name:      user.Name,
		College:   req.College,
		Grade:     user.Grade,
		AvatarURL: user.AvatarURL,
	}

	// 更新用户学院信息
	err = s.userRepo.UpdateUserProfile(req.StudentID, updateReq)
	if err != nil {
		return nil, err
	}

	// 返回更新后的学院信息
	return &model.UpdateCollegeResponse{
		StudentID: req.StudentID,
		College:   req.College,
		Message:   "学院信息更新成功",
	}, nil
}

// ChangePassword 修改用户密码
func (s *AuthService) ChangePassword(req *model.ChangePasswordRequest) error {
	// 1. 根据学号获取用户信息
	user, err := s.userRepo.GetUserByStudentID(req.StudentID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCredential
	}

	// 2. 校验旧密码
	if !checkPassword(user.Password, req.OldPassword) {
		return ErrInvalidCredential
	}

	// 3. 检查新密码是否与旧密码相同
	if req.OldPassword == req.NewPassword {
		return ErrPasswordSame
	}

	// 4. 哈希新密码
	newHashedPassword := hashPassword(req.NewPassword)

	// 5. 更新密码到数据库
	return s.userRepo.UpdatePassword(req.StudentID, newHashedPassword)
}
