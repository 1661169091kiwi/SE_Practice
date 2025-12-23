package service

import (
	"errors"
	"fmt"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
	"time"
)

type TeamChatService struct {
	chatRepo *repo.TeamChatRepo
	teamRepo *repo.TeamRepo
	userRepo *repo.UserRepo
}

func NewTeamChatService() *TeamChatService {
	return &TeamChatService{
		chatRepo: repo.NewTeamChatRepo(),
		teamRepo: repo.NewTeamRepo(),
		userRepo: repo.NewUserRepo(),
	}
}

var (
	ErrVoteNotFound  = errors.New("投票不存在")
	ErrVoteClosed    = errors.New("投票已关闭")
	ErrVoteExpired   = errors.New("投票已过期")
	ErrLeaveNotFound = errors.New("请假申请不存在")
)

// CreateMessage 创建消息
func (s *TeamChatService) CreateMessage(req *model.CreateMessageRequest, senderID string) (*model.TeamMessage, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(senderID, req.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 检查是否已批准
	teamMember, err := s.userRepo.GetTeamMemberByAthleteID(athlete.ID, req.TeamID)
	if err != nil {
		return nil, err
	}
	if teamMember == nil || !teamMember.IsApproved {
		return nil, errors.New("您尚未被批准加入该队伍")
	}

	msg := &model.TeamMessage{
		TeamID:      req.TeamID,
		SenderID:    senderID,
		MessageType: req.MessageType,
		Content:     req.Content,
	}
	if msg.MessageType == "" {
		msg.MessageType = "text"
	}

	messageID, err := s.chatRepo.CreateMessage(msg)
	if err != nil {
		return nil, err
	}
	msg.MessageID = messageID

	// 获取发送者姓名
	user, _ := s.userRepo.GetUserByStudentID(senderID)
	if user != nil {
		msg.SenderName = user.Name
	}

	return msg, nil
}

// GetTeamMessages 获取队伍消息列表
func (s *TeamChatService) GetTeamMessages(teamID int64, userID string, limit, offset int) ([]model.TeamMessage, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	messages, err := s.chatRepo.GetTeamMessages(teamID, limit, offset)
	if err != nil {
		return nil, err
	}

	// 获取总人数和已读状态
	totalCount, _ := s.chatRepo.GetTeamMemberCount(teamID)
	for i := range messages {
		messages[i].TotalCount = totalCount
		// 检查当前用户是否已读
		readStatus, _ := s.chatRepo.GetMessageReadStatus(messages[i].MessageID)
		messages[i].ReadCount = len(readStatus)
		for _, status := range readStatus {
			if status.ReaderID == userID {
				messages[i].IsRead = true
				break
			}
		}
	}

	return messages, nil
}

// MarkMessageAsRead 标记消息为已读
func (s *TeamChatService) MarkMessageAsRead(messageID int64, readerID string) error {
	return s.chatRepo.MarkMessageAsRead(messageID, readerID)
}

// CreateVote 创建投票
func (s *TeamChatService) CreateVote(req *model.CreateVoteRequest, creatorID string) (*model.TeamVote, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(creatorID, req.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	if len(req.Options) < 2 {
		return nil, errors.New("至少需要2个选项")
	}

	// 先创建消息
	msgReq := &model.CreateMessageRequest{
		TeamID:      req.TeamID,
		MessageType: "vote",
		Content:     req.Title,
	}
	msg, err := s.CreateMessage(msgReq, creatorID)
	if err != nil {
		return nil, err
	}

	// 解析截止时间
	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		// 尝试多种时间格式
		formats := []string{
			"2006-01-02T15:04:05Z07:00", // ISO 8601 with timezone
			"2006-01-02T15:04:05",       // ISO 8601 without timezone
			"2006-01-02T15:04",           // datetime-local format
			"2006-01-02 15:04:05",        // Standard format
			"2006-01-02 15:04",           // Standard format without seconds
		}
		var t time.Time
		var err error
		parsed := false
		for _, format := range formats {
			t, err = time.Parse(format, *req.Deadline)
			if err == nil {
				parsed = true
				break
			}
		}
		if !parsed {
			return nil, fmt.Errorf("无效的截止时间格式: %s", *req.Deadline)
		}
		deadline = &t
	}

	vote := &model.TeamVote{
		TeamID:      req.TeamID,
		CreatorID:   creatorID,
		MessageID:   msg.MessageID,
		Title:       req.Title,
		Description: req.Description,
		Options:     req.Options,
		IsMultiple:  req.IsMultiple,
		Deadline:    deadline,
		Status:      "active",
	}

	voteID, err := s.chatRepo.CreateVote(vote)
	if err != nil {
		return nil, err
	}
	vote.VoteID = voteID

	// 获取创建者姓名
	user, _ := s.userRepo.GetUserByStudentID(creatorID)
	if user != nil {
		vote.CreatorName = user.Name
	}

	return vote, nil
}

// GetTeamVotes 获取队伍投票列表
func (s *TeamChatService) GetTeamVotes(teamID int64, userID string, status string) ([]model.TeamVote, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	votes, err := s.chatRepo.GetTeamVotes(teamID, status)
	if err != nil {
		return nil, err
	}

	// 获取投票结果和用户投票
	for i := range votes {
		records, _ := s.chatRepo.GetVoteRecords(votes[i].VoteID)
		votes[i].VoteCount = len(records)

		// 计算每个选项的投票数
		optionCounts := make(map[int]int)
		for _, options := range records {
			for _, optIdx := range options {
				optionCounts[optIdx]++
			}
		}

		// 构建结果
		totalVotes := len(records)
		votes[i].Results = make([]model.VoteOptionResult, len(votes[i].Options))
		for j, option := range votes[i].Options {
			count := optionCounts[j]
			percentage := 0.0
			if totalVotes > 0 {
				percentage = float64(count) * 100.0 / float64(totalVotes)
			}
			votes[i].Results[j] = model.VoteOptionResult{
				OptionIndex: j,
				OptionText:  option,
				VoteCount:   count,
				Percentage:  percentage,
			}
		}

		// 获取用户投票
		userVote, _ := s.chatRepo.GetUserVote(votes[i].VoteID, userID)
		votes[i].UserVote = userVote

		// 检查是否过期
		if votes[i].Deadline != nil && time.Now().After(*votes[i].Deadline) && votes[i].Status == "active" {
			s.chatRepo.CloseVote(votes[i].VoteID)
			votes[i].Status = "closed"
		}
	}

	return votes, nil
}

// GetVoteByMessageID 根据message_id获取投票详情
func (s *TeamChatService) GetVoteByMessageID(messageID int64, userID string) (*model.TeamVote, error) {
	vote, err := s.chatRepo.GetVoteByMessageID(messageID)
	if err != nil {
		return nil, err
	}
	if vote == nil {
		return nil, ErrVoteNotFound
	}

	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, vote.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 获取投票结果和用户投票
	records, _ := s.chatRepo.GetVoteRecords(vote.VoteID)
	vote.VoteCount = len(records)

	// 计算每个选项的投票数
	optionCounts := make(map[int]int)
	for _, options := range records {
		for _, optIdx := range options {
			optionCounts[optIdx]++
		}
	}

	// 构建结果
	totalVotes := len(records)
	vote.Results = make([]model.VoteOptionResult, len(vote.Options))
	for j, option := range vote.Options {
		count := optionCounts[j]
		percentage := 0.0
		if totalVotes > 0 {
			percentage = float64(count) * 100.0 / float64(totalVotes)
		}
		vote.Results[j] = model.VoteOptionResult{
			OptionIndex: j,
			OptionText:  option,
			VoteCount:   count,
			Percentage:  percentage,
		}
	}

	// 获取用户投票
	userVote, _ := s.chatRepo.GetUserVote(vote.VoteID, userID)
	vote.UserVote = userVote

	// 检查是否过期
	if vote.Deadline != nil && time.Now().After(*vote.Deadline) && vote.Status == "active" {
		s.chatRepo.CloseVote(vote.VoteID)
		vote.Status = "closed"
	}

	return vote, nil
}

// Vote 投票
func (s *TeamChatService) Vote(req *model.VoteRequest, voterID string) error {
	vote, err := s.chatRepo.GetVoteByID(req.VoteID)
	if err != nil {
		return err
	}
	if vote == nil {
		return ErrVoteNotFound
	}

	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(voterID, vote.TeamID)
	if err != nil {
		return err
	}
	if athlete == nil {
		return ErrNotTeamMember
	}

	if vote.Status == "closed" {
		return ErrVoteClosed
	}

	if vote.Deadline != nil && time.Now().After(*vote.Deadline) {
		return ErrVoteExpired
	}

	// 验证选项索引
	if !vote.IsMultiple && len(req.SelectedOptions) > 1 {
		return errors.New("此投票不允许多选")
	}

	for _, idx := range req.SelectedOptions {
		if idx < 0 || idx >= len(vote.Options) {
			return errors.New("无效的选项索引")
		}
	}

	return s.chatRepo.CreateVoteRecord(req.VoteID, voterID, req.SelectedOptions)
}

// CreateNotification 创建通知
func (s *TeamChatService) CreateNotification(req *model.CreateNotificationRequest, senderID string) (*model.TeamNotification, error) {
	// 检查是否是队伍成员（通知通常由队长或管理员创建）
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(senderID, req.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 先创建消息
	msgReq := &model.CreateMessageRequest{
		TeamID:      req.TeamID,
		MessageType: "notification",
		Content:     req.Content,
	}
	msg, err := s.CreateMessage(msgReq, senderID)
	if err != nil {
		return nil, err
	}

	notif := &model.TeamNotification{
		TeamID:           req.TeamID,
		SenderID:         senderID,
		MessageID:        msg.MessageID,
		Title:            req.Title,
		Content:          req.Content,
		NotificationType: req.NotificationType,
	}
	if notif.NotificationType == "" {
		notif.NotificationType = "info"
	}

	notifID, err := s.chatRepo.CreateNotification(notif)
	if err != nil {
		return nil, err
	}
	notif.NotificationID = notifID

	// 获取发送者姓名
	user, _ := s.userRepo.GetUserByStudentID(senderID)
	if user != nil {
		notif.SenderName = user.Name
	}

	return notif, nil
}

// GetTeamNotifications 获取队伍通知列表
func (s *TeamChatService) GetTeamNotifications(teamID int64, userID string, limit int) ([]model.TeamNotification, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	notifications, err := s.chatRepo.GetTeamNotifications(teamID, limit)
	if err != nil {
		return nil, err
	}

	// 获取已读状态
	for i := range notifications {
		readStatus, _ := s.chatRepo.GetMessageReadStatus(notifications[i].MessageID)
		notifications[i].ReadCount = len(readStatus)
		totalCount, _ := s.chatRepo.GetTeamMemberCount(teamID)
		notifications[i].TotalCount = totalCount
		for _, status := range readStatus {
			if status.ReaderID == userID {
				notifications[i].IsRead = true
				break
			}
		}
	}

	return notifications, nil
}

// CreateLeaveRequest 创建请假申请
func (s *TeamChatService) CreateLeaveRequest(req *model.CreateLeaveRequest, applicantID string) (*model.LeaveRequest, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(applicantID, req.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("无效的开始日期格式")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("无效的结束日期格式")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("结束日期不能早于开始日期")
	}

	// 先创建消息通知队长
	user, _ := s.userRepo.GetUserByStudentID(applicantID)
	applicantName := applicantID
	if user != nil {
		applicantName = user.Name
	}
	
	leaveTypeText := map[string]string{
		"sick":    "病假",
		"personal": "事假",
		"other":   "其他",
	}[req.LeaveType]
	if leaveTypeText == "" {
		leaveTypeText = "请假"
	}
	
	msgContent := fmt.Sprintf("%s 申请%s，时间：%s 至 %s，原因：%s", 
		applicantName, leaveTypeText, req.StartDate, req.EndDate, req.Reason)
	
	msgReq := &model.CreateMessageRequest{
		TeamID:      req.TeamID,
		MessageType: "leave_request",
		Content:     msgContent,
	}
	msg, err := s.CreateMessage(msgReq, applicantID)
	if err != nil {
		return nil, err
	}

	leaveReq := &model.LeaveRequest{
		TeamID:    req.TeamID,
		ApplicantID: applicantID,
		MessageID: &msg.MessageID,
		LeaveType: req.LeaveType,
		StartDate: startDate,
		EndDate:   endDate,
		Reason:    req.Reason,
		Status:    "pending",
	}

	leaveID, err := s.chatRepo.CreateLeaveRequest(leaveReq)
	if err != nil {
		return nil, err
	}
	leaveReq.LeaveID = leaveID

	// 获取申请人姓名
	if user != nil {
		leaveReq.ApplicantName = user.Name
	}

	return leaveReq, nil
}

// GetLeaveRequests 获取请假申请列表
func (s *TeamChatService) GetLeaveRequests(teamID int64, userID string, applicantID string, status string) ([]model.LeaveRequest, error) {
	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, teamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 目前放宽权限：
	// - 任何队伍成员在未指定 applicantID 时都可以查看该队伍的所有请假申请
	// - 如果指定了 applicantID，则按 applicantID 过滤

	requests, err := s.chatRepo.GetLeaveRequests(teamID, applicantID, status)
	if err != nil {
		return nil, err
	}

	// 获取申请人姓名
	for i := range requests {
		if requests[i].ApplicantName == "" {
			user, _ := s.userRepo.GetUserByStudentID(requests[i].ApplicantID)
			if user != nil {
				requests[i].ApplicantName = user.Name
			}
		}
		if requests[i].ReviewerName == "" && requests[i].ReviewerID != nil {
			user, _ := s.userRepo.GetUserByStudentID(*requests[i].ReviewerID)
			if user != nil {
				requests[i].ReviewerName = user.Name
			}
		}
	}

	return requests, nil
}

// ReviewLeaveRequest 审核请假申请
func (s *TeamChatService) ReviewLeaveRequest(req *model.ReviewLeaveRequest, reviewerID string) error {
	leaveReq, err := s.chatRepo.GetLeaveRequestByID(req.LeaveID)
	if err != nil {
		return err
	}
	if leaveReq == nil {
		return ErrLeaveNotFound
	}

	// 检查是否是队长
	isCaptain, err := s.teamRepo.IsTeamCaptain(reviewerID, leaveReq.TeamID)
	if err != nil {
		return err
	}
	if !isCaptain {
		return ErrNotCaptain
	}

	if req.Status != "approved" && req.Status != "rejected" {
		return errors.New("无效的审核状态")
	}

	// 更新请假申请状态
	err = s.chatRepo.ReviewLeaveRequest(req.LeaveID, reviewerID, req.Status, req.ReviewComment)
	if err != nil {
		return err
	}

	// 如果有关联的消息，更新消息内容以显示审核结果
	if leaveReq.MessageID != nil {
		reviewer, _ := s.userRepo.GetUserByStudentID(reviewerID)
		reviewerName := reviewerID
		if reviewer != nil {
			reviewerName = reviewer.Name
		}
		
		statusText := map[string]string{
			"approved": "已批准",
			"rejected": "已拒绝",
		}[req.Status]
		
		// 获取原始消息
		msg, err := s.chatRepo.GetMessageByID(*leaveReq.MessageID)
		if err == nil && msg != nil {
			// 更新消息内容，添加审核结果
			newContent := msg.Content
			if req.ReviewComment != "" {
				newContent = fmt.Sprintf("%s\n【审核结果】%s - %s（%s）", msg.Content, statusText, reviewerName, req.ReviewComment)
			} else {
				newContent = fmt.Sprintf("%s\n【审核结果】%s - %s", msg.Content, statusText, reviewerName)
			}
			s.chatRepo.UpdateMessageContent(*leaveReq.MessageID, newContent)
		}
	}

	return nil
}

// DeleteMessage 删除消息
func (s *TeamChatService) DeleteMessage(messageID int64, userID string) error {
	// 获取消息信息
	msg, err := s.chatRepo.GetMessageByID(messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return errors.New("消息不存在")
	}

	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, msg.TeamID)
	if err != nil {
		return err
	}
	if athlete == nil {
		return ErrNotTeamMember
	}

	// 检查是否有删除权限：
	// - 普通成员：只能删除自己发送的消息
	// - 队长：可以删除本队伍内的任意消息（包括投票、通知、请假等）
	if msg.SenderID != userID && !athlete.IsCaptain {
		return errors.New("只能删除自己发送的消息")
	}

	// 根据消息类型，可能需要删除关联的数据
	switch msg.MessageType {
	case "vote":
		// 删除投票及其记录
		vote, _ := s.chatRepo.GetVoteByMessageID(messageID)
		if vote != nil {
			// 删除投票记录
			_, _ = db.Exec("DELETE FROM vote_records WHERE vote_id = ?", vote.VoteID)
			// 删除投票
			_, _ = db.Exec("DELETE FROM team_votes WHERE vote_id = ?", vote.VoteID)
		}
	case "notification":
		// 删除通知
		_, _ = db.Exec("DELETE FROM team_notifications WHERE message_id = ?", messageID)
	case "leave_request":
		// 删除请假申请
		_, _ = db.Exec("DELETE FROM leave_requests WHERE message_id = ?", messageID)
	}

	// 此处不再按 sender_id 限制，权限已在上面检查
	return s.chatRepo.DeleteMessage(messageID)
}

// GetLeaveRequestByMessageID 根据message_id获取请假申请
func (s *TeamChatService) GetLeaveRequestByMessageID(messageID int64, userID string) (*model.LeaveRequest, error) {
	leaveReq, err := s.chatRepo.GetLeaveRequestByMessageID(messageID)
	if err != nil {
		return nil, err
	}
	if leaveReq == nil {
		return nil, ErrLeaveNotFound
	}

	// 检查是否是队伍成员
	athlete, err := s.userRepo.GetAthleteByStudentIDAndTeamID(userID, leaveReq.TeamID)
	if err != nil {
		return nil, err
	}
	if athlete == nil {
		return nil, ErrNotTeamMember
	}

	// 获取申请人姓名
	if leaveReq.ApplicantName == "" {
		user, _ := s.userRepo.GetUserByStudentID(leaveReq.ApplicantID)
		if user != nil {
			leaveReq.ApplicantName = user.Name
		}
	}
	// 获取审核人姓名
	if leaveReq.ReviewerName == "" && leaveReq.ReviewerID != nil {
		user, _ := s.userRepo.GetUserByStudentID(*leaveReq.ReviewerID)
		if user != nil {
			leaveReq.ReviewerName = user.Name
		}
	}

	return leaveReq, nil
}

