package model

import "time"

// MatchComment 比赛评论
type MatchComment struct {
	ID        int64     `json:"id"`
	MatchID   int64     `json:"match_id"`
	StudentID string    `json:"student_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	// Join fields
	UserName  string `json:"user_name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	MatchID int64  `json:"match_id"`
	Content string `json:"content"`
}
