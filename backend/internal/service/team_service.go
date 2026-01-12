package service

import (
	"errors"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

type TeamService struct {
	teamRepo *repo.TeamRepo
}

func NewTeamService() *TeamService {
	return &TeamService{
		teamRepo: repo.NewTeamRepo(),
	}
}

// CreateTeam 创建队伍
func (s *TeamService) CreateTeam(req *model.CreateTeamRequest) (*model.Team, error) {
	if req.TeamName == "" {
		return nil, errors.New("team name is required")
	}
	if req.SportID <= 0 {
		return nil, errors.New("sport id is required")
	}
	if req.CreatedBy == "" {
		return nil, errors.New("creator is required")
	}

	team := &model.Team{
		TeamName:    req.TeamName,
		SportID:     req.SportID,
		College:     req.College,
		TeamType:    req.TeamType,
		AvatarURL:   req.AvatarURL,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		IsApproved:  true, // Admin created teams are approved by default
	}

	id, err := s.teamRepo.CreateTeam(team)
	if err != nil {
		return nil, err
	}
	if err := s.teamRepo.ApproveTeam(id); err != nil {
		return nil, err
	}
	team.ID = id
	team.IsApproved = true
	return team, nil
}

// GetMyTeams 获取用户加入的队伍
func (s *TeamService) GetMyTeams(studentID string) ([]model.Team, error) {
	return s.teamRepo.GetTeamsByStudentID(studentID)
}

// ApplyCreateTeam 申请创建队伍
func (s *TeamService) ApplyCreateTeam(req *model.CreateTeamRequest) (*model.Team, error) {
	if req.TeamName == "" {
		return nil, errors.New("team name is required")
	}
	if req.SportID <= 0 {
		return nil, errors.New("sport id is required")
	}
	if req.CreatedBy == "" {
		return nil, errors.New("creator is required")
	}

	team := &model.Team{
		TeamName:    req.TeamName,
		SportID:     req.SportID,
		College:     req.College,
		TeamType:    req.TeamType,
		AvatarURL:   req.AvatarURL,
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
		IsApproved:  false, // Applied teams need approval
	}

	id, err := s.teamRepo.CreateTeam(team)
	if err != nil {
		return nil, err
	}
	team.ID = id
	return team, nil
}

// ListTeams 列出队伍
func (s *TeamService) ListTeams(sportID int64) ([]model.Team, error) {
	return s.teamRepo.ListTeams(sportID)
}

// GetTeamDetail 获取队伍详情
func (s *TeamService) GetTeamDetail(id int64) (*model.Team, error) {
	team, err := s.teamRepo.GetTeamByID(id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}
	return team, nil
}

// GetPendingTeams 获取待审核队伍
func (s *TeamService) GetPendingTeams() ([]model.Team, error) {
	return s.teamRepo.ListPendingTeams()
}

// ApproveTeam 批准队伍
func (s *TeamService) ApproveTeam(teamID int64) error {
	return s.teamRepo.ApproveTeam(teamID)
}

// DeleteTeam 删除队伍
func (s *TeamService) DeleteTeam(teamID int64) error {
	return s.teamRepo.DeleteTeam(teamID)
}

// GetTeamMembers 获取队伍成员
func (s *TeamService) GetTeamMembers(teamID int64) ([]model.TeamMemberDetail, error) {
	return s.teamRepo.GetTeamMembers(teamID)
}

// GetPendingTeamMembers 获取待审核成员
func (s *TeamService) GetPendingTeamMembers(teamID int64, userID, role string) ([]model.TeamMemberDetail, error) {
	if role == "admin" {
		return s.teamRepo.GetPendingTeamMembers(teamID)
	}

	// 检查是否是队伍创建者或队长
	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	// 检查是否是创建者
	if team.CreatedBy == userID {
		return s.teamRepo.GetPendingTeamMembers(teamID)
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(userID, teamID)
	if err != nil {
		return nil, err
	}
	if !isCaptain {
		return nil, errors.New("permission denied")
	}

	return s.teamRepo.GetPendingTeamMembers(teamID)
}

// ApproveMemberWithAuth 批准成员（支持管理员、创建者、队长）
func (s *TeamService) ApproveMemberWithAuth(teamMemberID int64, userID, role string) error {
	if role == "admin" {
		return s.teamRepo.ApproveMember(teamMemberID)
	}

	member, err := s.teamRepo.GetTeamMemberByID(teamMemberID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("member not found")
	}

	team, err := s.teamRepo.GetTeamByID(member.TeamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	// 检查是否是创建者
	if team.CreatedBy == userID {
		return s.teamRepo.ApproveMember(teamMemberID)
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(userID, member.TeamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return errors.New("permission denied: only admin, team creator, or captain can approve members")
	}

	return s.teamRepo.ApproveMember(teamMemberID)
}

// RemoveMember 移除成员
func (s *TeamService) RemoveMember(memberID int64) error {
	return s.teamRepo.RemoveMember(memberID)
}

// RemoveMemberWithAuth 移除成员（带权限检查，支持管理员、创建者、队长）
func (s *TeamService) RemoveMemberWithAuth(memberID int64, userID, role string) error {
	if role == "admin" {
		return s.teamRepo.RemoveMember(memberID)
	}

	member, err := s.teamRepo.GetTeamMemberByID(memberID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("member not found")
	}

	team, err := s.teamRepo.GetTeamByID(member.TeamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	// 检查是否是创建者
	if team.CreatedBy == userID {
		return s.teamRepo.RemoveMember(memberID)
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(userID, member.TeamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return errors.New("permission denied: only admin, team creator, or captain can remove members")
	}

	return s.teamRepo.RemoveMember(memberID)
}

// UpdateTeamName 更新队伍名称（支持管理员、创建者、队长）
func (s *TeamService) UpdateTeamName(teamID int64, name, userID, role string) error {
	// Verify permission: Admin, Creator, or Captain
	if role == "admin" {
		return s.teamRepo.UpdateTeamName(teamID, name)
	}

	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	// 检查是否是创建者
	if team.CreatedBy == userID {
		return s.teamRepo.UpdateTeamName(teamID, name)
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(userID, teamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return errors.New("permission denied: only admin, team creator, or captain can update team name")
	}

	return s.teamRepo.UpdateTeamName(teamID, name)
}

// UpdateTeamAvatar 更新队伍头像（支持管理员、创建者、队长）
func (s *TeamService) UpdateTeamAvatar(teamID int64, avatarURL, userID, role string) error {
	// Verify permission: Admin, Creator, or Captain
	if role == "admin" {
		return s.teamRepo.UpdateTeamAvatar(teamID, avatarURL)
	}

	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}

	// 检查是否是创建者
	if team.CreatedBy == userID {
		return s.teamRepo.UpdateTeamAvatar(teamID, avatarURL)
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(userID, teamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return errors.New("permission denied: only admin, team creator, or captain can update team avatar")
	}

	return s.teamRepo.UpdateTeamAvatar(teamID, avatarURL)
}
