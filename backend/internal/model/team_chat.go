package model

import "time"

// TeamMessage 队内聊天消息
type TeamMessage struct {
	MessageID   int64     `json:"message_id"`
	TeamID      int64     `json:"team_id"`
	SenderID    string    `json:"sender_id"`
	SenderName  string    `json:"sender_name,omitempty"`
	MessageType string    `json:"message_type"` // text, image, file, vote, notification, leave_request
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	IsRead      bool      `json:"is_read,omitempty"` // 当前用户是否已读
	ReadCount   int       `json:"read_count,omitempty"` // 已读人数
	TotalCount  int       `json:"total_count,omitempty"` // 总人数
}

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	TeamID      int64  `json:"team_id" binding:"required"`
	MessageType string `json:"message_type"` // 默认为text
	Content     string `json:"content" binding:"required"`
}

// MessageReadStatus 消息已读状态
type MessageReadStatus struct {
	ReadID    int64     `json:"read_id"`
	MessageID int64     `json:"message_id"`
	ReaderID  string    `json:"reader_id"`
	ReaderName string   `json:"reader_name,omitempty"`
	ReadAt    time.Time `json:"read_at"`
}

// TeamVote 投票
type TeamVote struct {
	VoteID      int64     `json:"vote_id"`
	TeamID      int64     `json:"team_id"`
	CreatorID   string    `json:"creator_id"`
	CreatorName string    `json:"creator_name,omitempty"`
	MessageID   int64     `json:"message_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Options     []string  `json:"options"` // 选项列表
	IsMultiple  bool      `json:"is_multiple"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	Status      string    `json:"status"` // active, closed
	CreatedAt   time.Time `json:"created_at"`
	VoteCount   int       `json:"vote_count,omitempty"` // 投票人数
	Results     []VoteOptionResult `json:"results,omitempty"` // 投票结果
	UserVote    []int     `json:"user_vote,omitempty"` // 当前用户的投票（选项索引）
}

// VoteOptionResult 投票选项结果
type VoteOptionResult struct {
	OptionIndex int    `json:"option_index"`
	OptionText  string `json:"option_text"`
	VoteCount   int    `json:"vote_count"`
	Percentage  float64 `json:"percentage"`
}

// CreateVoteRequest 创建投票请求
type CreateVoteRequest struct {
	TeamID      int64    `json:"team_id" binding:"required"`
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Options     []string `json:"options" binding:"required,min=2"`
	IsMultiple  bool     `json:"is_multiple"`
	Deadline    *string  `json:"deadline"` // ISO 8601格式
}

// VoteRequest 投票请求
type VoteRequest struct {
	VoteID         int64  `json:"vote_id" binding:"required"`
	SelectedOptions []int `json:"selected_options" binding:"required"` // 选项索引数组
}

// TeamNotification 通知
type TeamNotification struct {
	NotificationID int64     `json:"notification_id"`
	TeamID         int64     `json:"team_id"`
	SenderID       string    `json:"sender_id"`
	SenderName     string    `json:"sender_name,omitempty"`
	MessageID      int64     `json:"message_id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	NotificationType string  `json:"notification_type"` // info, warning, urgent
	CreatedAt      time.Time `json:"created_at"`
	IsRead         bool      `json:"is_read,omitempty"`
	ReadCount      int       `json:"read_count,omitempty"`
	TotalCount     int       `json:"total_count,omitempty"`
}

// CreateNotificationRequest 创建通知请求
type CreateNotificationRequest struct {
	TeamID           int64  `json:"team_id" binding:"required"`
	Title            string `json:"title" binding:"required"`
	Content          string `json:"content" binding:"required"`
	NotificationType string `json:"notification_type"` // info, warning, urgent
}

// LeaveRequest 请假申请
type LeaveRequest struct {
	LeaveID       int64     `json:"leave_id"`
	TeamID        int64     `json:"team_id"`
	ApplicantID   string    `json:"applicant_id"`
	ApplicantName string    `json:"applicant_name,omitempty"`
	MessageID     *int64    `json:"message_id,omitempty"`
	LeaveType     string    `json:"leave_type"` // sick, personal, other
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	Reason        string    `json:"reason"`
	Status        string    `json:"status"` // pending, approved, rejected
	ReviewerID    *string   `json:"reviewer_id,omitempty"`
	ReviewerName  string    `json:"reviewer_name,omitempty"`
	ReviewComment *string   `json:"review_comment,omitempty"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateLeaveRequest 创建请假申请请求
type CreateLeaveRequest struct {
	TeamID    int64  `json:"team_id" binding:"required"`
	LeaveType string `json:"leave_type" binding:"required"` // sick, personal, other
	StartDate string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate   string `json:"end_date" binding:"required"` // YYYY-MM-DD
	Reason    string `json:"reason" binding:"required"`
}

// ReviewLeaveRequest 审核请假申请请求
type ReviewLeaveRequest struct {
	LeaveID       int64  `json:"leave_id" binding:"required"`
	Status        string `json:"status" binding:"required"` // approved, rejected
	ReviewComment string `json:"review_comment"`
}

