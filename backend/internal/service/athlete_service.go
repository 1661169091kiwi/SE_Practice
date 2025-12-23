package service

import (
	"errors"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
)

type AthleteService struct {
	userRepo *repo.UserRepo
	teamRepo *repo.TeamRepo
}

func NewAthleteService() *AthleteService {
	return &AthleteService{
		userRepo: repo.NewUserRepo(),
		teamRepo: repo.NewTeamRepo(),
	}
}

var (
	ErrAthleteNotFound = errors.New("运动员信息不存在")
	ErrNotTeamMember   = errors.New("您不是该队伍的成员")
	ErrNotCaptain      = errors.New("您不是该队伍的队长，无权执行此操作")
)

// GetMyAthleteInfo 获取当前用户的运动员信息
func (s *AthleteService) GetMyAthleteInfo(studentID string) ([]model.Athlete, error) {
	return s.userRepo.ListAthletesByStudentID(studentID)
}

// GetMyAthleteInfoByTeam 获取当前用户在指定队伍中的运动员信息
func (s *AthleteService) GetMyAthleteInfoByTeam(studentID string, teamID int64) (*model.Athlete, error) {
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(studentID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrAthleteNotFound
	}
	return athlete, nil
}

// UpdateMyAthleteInfo 更新当前用户的运动员信息（如球衣号码）
func (s *AthleteService) UpdateMyAthleteInfo(studentID string, teamID int64, jerseyNumber string) error {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(studentID, teamID)
	if err != nil {
		return err
	}
	if athlete == nil {
		return ErrNotTeamMember
	}

	// 检查是否已批准
	teamMember, err := s.userRepo.GetTeamMemberByAthleteID(athlete.ID, teamID)
	if err != nil {
		return err
	}
	if teamMember == nil || !teamMember.IsApproved {
		return errors.New("您尚未被批准加入该队伍")
	}

	// 更新信息
	return s.userRepo.UpdateAthleteInfo(athlete.ID, jerseyNumber)
}

// LeaveTeam 退出队伍
func (s *AthleteService) LeaveTeam(studentID string, teamID int64) error {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(studentID, teamID)
	if err != nil {
		return err
	}
	if athlete == nil {
		return ErrNotTeamMember
	}

	// 检查是否是队长，队长不能退出（需要先转让队长）
	if athlete.IsCaptain {
		return errors.New("队长不能直接退出队伍，请先转让队长身份")
	}

	// 退出队伍
	return s.userRepo.LeaveTeam(studentID, teamID)
}

// IsTeamCaptain 检查用户是否是某个队伍的队长
func (s *AthleteService) IsTeamCaptain(studentID string, teamID int64) (bool, error) {
	return s.teamRepo.IsTeamCaptain(studentID, teamID)
}

// GetTeamMembersForAthlete 运动员查看队伍成员（包含队长信息）
func (s *AthleteService) GetTeamMembersForAthlete(teamID int64, studentID string) ([]model.TeamMemberDetail, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(studentID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 检查是否已批准
	teamMember, err := s.userRepo.GetTeamMemberByAthleteID(athlete.ID, teamID)
	if err != nil {
		return nil, err
	}
	if teamMember == nil || !teamMember.IsApproved {
		return nil, errors.New("您尚未被批准加入该队伍")
	}

	// 获取队伍成员列表
	return s.teamRepo.GetTeamMembers(teamID)
}

// ApproveTeamMemberAsCaptain 队长批准队员申请
func (s *AthleteService) ApproveTeamMemberAsCaptain(teamMemberID int64, captainID string) error {
	// 获取成员信息
	member, err := s.teamRepo.GetTeamMemberByID(teamMemberID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("成员申请不存在")
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(captainID, member.TeamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return ErrNotCaptain
	}

	// 批准成员
	return s.teamRepo.ApproveMember(teamMemberID)
}

// RemoveTeamMemberAsCaptain 队长移除队员
func (s *AthleteService) RemoveTeamMemberAsCaptain(teamMemberID int64, captainID string) error {
	// 获取成员信息
	member, err := s.teamRepo.GetTeamMemberByID(teamMemberID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("成员不存在")
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(captainID, member.TeamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return ErrNotCaptain
	}

	// 检查是否是移除自己
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(captainID, member.TeamID)
	if err != nil {
		return err
	}
	if athlete != nil && member.AthleteID == athlete.ID {
		return errors.New("不能移除自己，请使用退出队伍功能")
	}

	// 移除成员
	return s.teamRepo.RemoveMember(teamMemberID)
}

