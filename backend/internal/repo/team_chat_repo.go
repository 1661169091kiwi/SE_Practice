package repo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"

)

type TeamChatRepo struct{}

func NewTeamChatRepo() *TeamChatRepo {
	return &TeamChatRepo{}
}

// CreateMessage 创建消息
func (r *TeamChatRepo) CreateMessage(msg *model.TeamMessage) (int64, error) {
	query := `INSERT INTO team_messages (team_id, sender_id, message_type, content) 
	          VALUES (?, ?, ?, ?)`
	return db.Insert(query, msg.TeamID, msg.SenderID, msg.MessageType, msg.Content)
}

// GetTeamMessages 获取队伍消息列表
func (r *TeamChatRepo) GetTeamMessages(teamID int64, limit, offset int) ([]model.TeamMessage, error) {
	query := `
		SELECT m.message_id, m.team_id, m.sender_id, COALESCE(u.name, ''), m.message_type, m.content, m.created_at
		FROM team_messages m
		LEFT JOIN users u ON m.sender_id = u.student_id
		WHERE m.team_id = ?
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := db.Query(query, teamID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.TeamMessage
	for rows.Next() {
		var msg model.TeamMessage
		if err := rows.Scan(&msg.MessageID, &msg.TeamID, &msg.SenderID, &msg.SenderName,
			&msg.MessageType, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

// MarkMessageAsRead 标记消息为已读
func (r *TeamChatRepo) MarkMessageAsRead(messageID int64, readerID string) error {
	query := `INSERT INTO message_read_status (message_id, reader_id) 
	          VALUES (?, ?) 
	          ON DUPLICATE KEY UPDATE read_at = CURRENT_TIMESTAMP`
	_, err := db.Exec(query, messageID, readerID)
	return err
}

// GetMessageReadStatus 获取消息已读状态
func (r *TeamChatRepo) GetMessageReadStatus(messageID int64) ([]model.MessageReadStatus, error) {
	query := `
		SELECT rs.read_id, rs.message_id, rs.reader_id, u.name, rs.read_at
		FROM message_read_status rs
		LEFT JOIN users u ON rs.reader_id = u.student_id
		WHERE rs.message_id = ?
		ORDER BY rs.read_at DESC
	`
	rows, err := db.Query(query, messageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []model.MessageReadStatus
	for rows.Next() {
		var status model.MessageReadStatus
		if err := rows.Scan(&status.ReadID, &status.MessageID, &status.ReaderID,
			&status.ReaderName, &status.ReadAt); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, rows.Err()
}

// GetTeamMemberCount 获取队伍成员总数
func (r *TeamChatRepo) GetTeamMemberCount(teamID int64) (int, error) {
	query := `
		SELECT COUNT(DISTINCT tm.athlete_id)
		FROM team_members tm
		JOIN athletes a ON tm.athlete_id = a.athlete_id
		WHERE tm.team_id = ? AND tm.is_approved = TRUE AND tm.is_active = TRUE
	`
	var count int
	err := db.QueryRow(query, teamID).Scan(&count)
	return count, err
}

// CreateVote 创建投票
func (r *TeamChatRepo) CreateVote(vote *model.TeamVote) (int64, error) {
	optionsJSON, _ := json.Marshal(vote.Options)
	query := `INSERT INTO team_votes (team_id, creator_id, message_id, title, description, options, is_multiple, deadline, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query, vote.TeamID, vote.CreatorID, vote.MessageID, vote.Title,
		vote.Description, string(optionsJSON), vote.IsMultiple, vote.Deadline, vote.Status)
}

// GetTeamVotes 获取队伍投票列表
func (r *TeamChatRepo) GetTeamVotes(teamID int64, status string) ([]model.TeamVote, error) {
	query := `
		SELECT v.vote_id, v.team_id, v.creator_id, COALESCE(u.name, ''), v.message_id, v.title, v.description, 
		       v.options, v.is_multiple, v.deadline, v.status, v.created_at
		FROM team_votes v
		LEFT JOIN users u ON v.creator_id = u.student_id
		WHERE v.team_id = ?
	`
	var args []interface{}
	args = append(args, teamID)
	if status != "" {
		query += " AND v.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY v.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []model.TeamVote
	for rows.Next() {
		var vote model.TeamVote
		var optionsJSON string
		var deadline sql.NullTime
		if err := rows.Scan(&vote.VoteID, &vote.TeamID, &vote.CreatorID, &vote.CreatorName,
			&vote.MessageID, &vote.Title, &vote.Description, &optionsJSON, &vote.IsMultiple,
			&deadline, &vote.Status, &vote.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(optionsJSON), &vote.Options)
		if deadline.Valid {
			vote.Deadline = &deadline.Time
		}
		votes = append(votes, vote)
	}
	return votes, rows.Err()
}

// GetVoteByID 根据ID获取投票
func (r *TeamChatRepo) GetVoteByID(voteID int64) (*model.TeamVote, error) {
	query := `
		SELECT v.vote_id, v.team_id, v.creator_id, COALESCE(u.name, ''), v.message_id, v.title, v.description, 
		       v.options, v.is_multiple, v.deadline, v.status, v.created_at
		FROM team_votes v
		LEFT JOIN users u ON v.creator_id = u.student_id
		WHERE v.vote_id = ?
	`
	var vote model.TeamVote
	var optionsJSON string
	var deadline sql.NullTime
	err := db.QueryRow(query, voteID).Scan(&vote.VoteID, &vote.TeamID, &vote.CreatorID, &vote.CreatorName,
		&vote.MessageID, &vote.Title, &vote.Description, &optionsJSON, &vote.IsMultiple,
		&deadline, &vote.Status, &vote.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(optionsJSON), &vote.Options)
	if deadline.Valid {
		vote.Deadline = &deadline.Time
	}
	return &vote, nil
}

// GetVoteByMessageID 根据message_id获取投票
func (r *TeamChatRepo) GetVoteByMessageID(messageID int64) (*model.TeamVote, error) {
	query := `
		SELECT v.vote_id, v.team_id, v.creator_id, COALESCE(u.name, ''), v.message_id, v.title, v.description, 
		       v.options, v.is_multiple, v.deadline, v.status, v.created_at
		FROM team_votes v
		LEFT JOIN users u ON v.creator_id = u.student_id
		WHERE v.message_id = ?
	`
	var vote model.TeamVote
	var optionsJSON string
	var deadline sql.NullTime
	err := db.QueryRow(query, messageID).Scan(&vote.VoteID, &vote.TeamID, &vote.CreatorID, &vote.CreatorName,
		&vote.MessageID, &vote.Title, &vote.Description, &optionsJSON, &vote.IsMultiple,
		&deadline, &vote.Status, &vote.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(optionsJSON), &vote.Options)
	if deadline.Valid {
		vote.Deadline = &deadline.Time
	}
	return &vote, nil
}

// CreateVoteRecord 创建投票记录
func (r *TeamChatRepo) CreateVoteRecord(voteID int64, voterID string, selectedOptions []int) error {
	optionsJSON, _ := json.Marshal(selectedOptions)
	query := `INSERT INTO vote_records (vote_id, voter_id, selected_options) 
	          VALUES (?, ?, ?) 
	          ON DUPLICATE KEY UPDATE selected_options = ?, voted_at = CURRENT_TIMESTAMP`
	_, err := db.Exec(query, voteID, voterID, string(optionsJSON), string(optionsJSON))
	return err
}

// GetVoteRecords 获取投票记录（返回投票者ID和选中的选项）
func (r *TeamChatRepo) GetVoteRecords(voteID int64) (map[string][]int, error) {
	query := `
		SELECT vr.voter_id, vr.selected_options
		FROM vote_records vr
		WHERE vr.vote_id = ?
	`
	rows, err := db.Query(query, voteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make(map[string][]int)
	for rows.Next() {
		var voterID string
		var optionsJSON string
		if err := rows.Scan(&voterID, &optionsJSON); err != nil {
			return nil, err
		}
		var options []int
		json.Unmarshal([]byte(optionsJSON), &options)
		records[voterID] = options
	}
	return records, rows.Err()
}

// GetUserVote 获取用户的投票
func (r *TeamChatRepo) GetUserVote(voteID int64, voterID string) ([]int, error) {
	query := `SELECT selected_options FROM vote_records WHERE vote_id = ? AND voter_id = ?`
	var optionsJSON string
	err := db.QueryRow(query, voteID, voterID).Scan(&optionsJSON)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var options []int
	json.Unmarshal([]byte(optionsJSON), &options)
	return options, nil
}

// CloseVote 关闭投票
func (r *TeamChatRepo) CloseVote(voteID int64) error {
	query := `UPDATE team_votes SET status = 'closed' WHERE vote_id = ?`
	_, err := db.Exec(query, voteID)
	return err
}

// CreateNotification 创建通知
func (r *TeamChatRepo) CreateNotification(notif *model.TeamNotification) (int64, error) {
	query := `INSERT INTO team_notifications (team_id, sender_id, message_id, title, content, notification_type) 
	          VALUES (?, ?, ?, ?, ?, ?)`
	return db.Insert(query, notif.TeamID, notif.SenderID, notif.MessageID, notif.Title,
		notif.Content, notif.NotificationType)
}

// GetTeamNotifications 获取队伍通知列表
func (r *TeamChatRepo) GetTeamNotifications(teamID int64, limit int) ([]model.TeamNotification, error) {
	query := `
		SELECT n.notification_id, n.team_id, n.sender_id, COALESCE(u.name, ''), n.message_id, n.title, n.content, 
		       n.notification_type, n.created_at
		FROM team_notifications n
		LEFT JOIN users u ON n.sender_id = u.student_id
		WHERE n.team_id = ?
		ORDER BY n.created_at DESC
		LIMIT ?
	`
	rows, err := db.Query(query, teamID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.TeamNotification
	for rows.Next() {
		var notif model.TeamNotification
		if err := rows.Scan(&notif.NotificationID, &notif.TeamID, &notif.SenderID, &notif.SenderName,
			&notif.MessageID, &notif.Title, &notif.Content, &notif.NotificationType, &notif.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notif)
	}
	return notifications, rows.Err()
}

// CreateLeaveRequest 创建请假申请
func (r *TeamChatRepo) CreateLeaveRequest(req *model.LeaveRequest) (int64, error) {
	query := `INSERT INTO leave_requests (team_id, applicant_id, message_id, leave_type, start_date, end_date, reason, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return db.Insert(query, req.TeamID, req.ApplicantID, req.MessageID, req.LeaveType,
		req.StartDate, req.EndDate, req.Reason, req.Status)
}

// GetLeaveRequests 获取请假申请列表
func (r *TeamChatRepo) GetLeaveRequests(teamID int64, applicantID string, status string) ([]model.LeaveRequest, error) {
	query := `
		SELECT l.leave_id, l.team_id, l.applicant_id, u1.name, l.message_id, l.leave_type, 
		       l.start_date, l.end_date, l.reason, l.status, l.reviewer_id, u2.name, 
		       l.review_comment, l.reviewed_at, l.created_at
		FROM leave_requests l
		LEFT JOIN users u1 ON l.applicant_id = u1.student_id
		LEFT JOIN users u2 ON l.reviewer_id = u2.student_id
		WHERE l.team_id = ?
	`
	var args []interface{}
	args = append(args, teamID)
	if applicantID != "" {
		query += " AND l.applicant_id = ?"
		args = append(args, applicantID)
	}
	if status != "" {
		query += " AND l.status = ?"
		args = append(args, status)
	}
	query += " ORDER BY l.created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []model.LeaveRequest
	for rows.Next() {
		var req model.LeaveRequest
		var messageID sql.NullInt64
		var reviewerID sql.NullString
		var reviewerName sql.NullString
		var reviewComment sql.NullString
		var reviewedAt sql.NullTime
		if err := rows.Scan(&req.LeaveID, &req.TeamID, &req.ApplicantID, &req.ApplicantName,
			&messageID, &req.LeaveType, &req.StartDate, &req.EndDate, &req.Reason, &req.Status,
			&reviewerID, &reviewerName, &reviewComment, &reviewedAt, &req.CreatedAt); err != nil {
			return nil, err
		}
		if messageID.Valid {
			req.MessageID = &messageID.Int64
		}
		if reviewerID.Valid {
			req.ReviewerID = &reviewerID.String
		}
		if reviewerName.Valid {
			req.ReviewerName = reviewerName.String
		}
		if reviewComment.Valid {
			req.ReviewComment = &reviewComment.String
		}
		if reviewedAt.Valid {
			req.ReviewedAt = &reviewedAt.Time
		}
		requests = append(requests, req)
	}
	return requests, rows.Err()
}

// GetLeaveRequestByID 根据ID获取请假申请
func (r *TeamChatRepo) GetLeaveRequestByID(leaveID int64) (*model.LeaveRequest, error) {
	query := `
		SELECT l.leave_id, l.team_id, l.applicant_id, u1.name, l.message_id, l.leave_type, 
		       l.start_date, l.end_date, l.reason, l.status, l.reviewer_id, u2.name, 
		       l.review_comment, l.reviewed_at, l.created_at
		FROM leave_requests l
		LEFT JOIN users u1 ON l.applicant_id = u1.student_id
		LEFT JOIN users u2 ON l.reviewer_id = u2.student_id
		WHERE l.leave_id = ?
	`
	var req model.LeaveRequest
	var messageID sql.NullInt64
	var reviewerID sql.NullString
	var reviewerName sql.NullString
	var reviewComment sql.NullString
	var reviewedAt sql.NullTime
	err := db.QueryRow(query, leaveID).Scan(&req.LeaveID, &req.TeamID, &req.ApplicantID, &req.ApplicantName,
		&messageID, &req.LeaveType, &req.StartDate, &req.EndDate, &req.Reason, &req.Status,
		&reviewerID, &reviewerName, &reviewComment, &reviewedAt, &req.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if messageID.Valid {
		req.MessageID = &messageID.Int64
	}
	if reviewerID.Valid {
		req.ReviewerID = &reviewerID.String
	}
	if reviewerName.Valid {
		req.ReviewerName = reviewerName.String
	}
	if reviewComment.Valid {
		req.ReviewComment = &reviewComment.String
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}
	return &req, nil
}

// GetLeaveRequestByMessageID 根据message_id获取请假申请
func (r *TeamChatRepo) GetLeaveRequestByMessageID(messageID int64) (*model.LeaveRequest, error) {
	query := `
		SELECT l.leave_id, l.team_id, l.applicant_id, u1.name, l.message_id, l.leave_type, 
		       l.start_date, l.end_date, l.reason, l.status, l.reviewer_id, u2.name, 
		       l.review_comment, l.reviewed_at, l.created_at
		FROM leave_requests l
		LEFT JOIN users u1 ON l.applicant_id = u1.student_id
		LEFT JOIN users u2 ON l.reviewer_id = u2.student_id
		WHERE l.message_id = ?
	`
	var req model.LeaveRequest
	var msgID sql.NullInt64
	var reviewerID sql.NullString
	var reviewerName sql.NullString
	var reviewComment sql.NullString
	var reviewedAt sql.NullTime
	err := db.QueryRow(query, messageID).Scan(&req.LeaveID, &req.TeamID, &req.ApplicantID, &req.ApplicantName,
		&msgID, &req.LeaveType, &req.StartDate, &req.EndDate, &req.Reason, &req.Status,
		&reviewerID, &reviewerName, &reviewComment, &reviewedAt, &req.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if msgID.Valid {
		req.MessageID = &msgID.Int64
	}
	if reviewerID.Valid {
		req.ReviewerID = &reviewerID.String
	}
	if reviewerName.Valid {
		req.ReviewerName = reviewerName.String
	}
	if reviewComment.Valid {
		req.ReviewComment = &reviewComment.String
	}
	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}
	return &req, nil
}

// ReviewLeaveRequest 审核请假申请
func (r *TeamChatRepo) ReviewLeaveRequest(leaveID int64, reviewerID string, status string, comment string) error {
	query := `UPDATE leave_requests 
	          SET status = ?, reviewer_id = ?, review_comment = ?, reviewed_at = CURRENT_TIMESTAMP 
	          WHERE leave_id = ?`
	_, err := db.Exec(query, status, reviewerID, comment, leaveID)
	return err
}

// DeleteMessage 删除消息
// 权限校验在 Service 层完成，这里只按 message_id 删除记录
func (r *TeamChatRepo) DeleteMessage(messageID int64) error {
	query := `DELETE FROM team_messages WHERE message_id = ?`
	result, err := db.Exec(query, messageID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("消息不存在或无权删除")
	}
	return nil
}

// GetMessageByID 根据ID获取消息
func (r *TeamChatRepo) GetMessageByID(messageID int64) (*model.TeamMessage, error) {
	query := `
		SELECT m.message_id, m.team_id, m.sender_id, COALESCE(u.name, ''), m.message_type, m.content, m.created_at
		FROM team_messages m
		LEFT JOIN users u ON m.sender_id = u.student_id
		WHERE m.message_id = ?
	`
	var msg model.TeamMessage
	err := db.QueryRow(query, messageID).Scan(&msg.MessageID, &msg.TeamID, &msg.SenderID, &msg.SenderName,
		&msg.MessageType, &msg.Content, &msg.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// UpdateMessageContent 更新消息内容
func (r *TeamChatRepo) UpdateMessageContent(messageID int64, content string) error {
	query := `UPDATE team_messages SET content = ? WHERE message_id = ?`
	_, err := db.Exec(query, content, messageID)
	return err
}

