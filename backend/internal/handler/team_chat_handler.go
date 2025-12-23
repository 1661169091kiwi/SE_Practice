package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var teamChatService = service.NewTeamChatService()

// CreateMessage 创建消息
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	msg, err := teamChatService.CreateMessage(&req, claims.Sub)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, msg)
}

// GetTeamMessages 获取队伍消息列表
func GetTeamMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}

	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := teamChatService.GetTeamMessages(teamID, claims.Sub, limit, offset)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, messages)
}

// MarkMessageAsRead 标记消息为已读
func MarkMessageAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		MessageID int64 `json:"message_id" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := teamChatService.MarkMessageAsRead(req.MessageID, claims.Sub); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "marked as read"})
}

// CreateVote 创建投票
func CreateVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	vote, err := teamChatService.CreateVote(&req, claims.Sub)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, vote)
}

// GetTeamVotes 获取队伍投票列表
func GetTeamVotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}

	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	status := r.URL.Query().Get("status") // active, closed, 或空（全部）

	votes, err := teamChatService.GetTeamVotes(teamID, claims.Sub, status)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, votes)
}

// Vote 投票
func Vote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := teamChatService.Vote(&req, claims.Sub); err != nil {
		if err == service.ErrVoteNotFound || err == service.ErrVoteClosed || err == service.ErrVoteExpired {
			util.Error(w, http.StatusBadRequest, err.Error())
		} else if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "voted successfully"})
}

// CreateNotification 创建通知
func CreateNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	notif, err := teamChatService.CreateNotification(&req, claims.Sub)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, notif)
}

// GetTeamNotifications 获取队伍通知列表
func GetTeamNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}

	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	notifications, err := teamChatService.GetTeamNotifications(teamID, claims.Sub, limit)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, notifications)
}

// CreateLeaveRequest 创建请假申请
func CreateLeaveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	leaveReq, err := teamChatService.CreateLeaveRequest(&req, claims.Sub)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, leaveReq)
}

// GetLeaveRequests 获取请假申请列表
func GetLeaveRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	teamIDStr := r.URL.Query().Get("team_id")
	if teamIDStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}

	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	applicantID := r.URL.Query().Get("applicant_id") // 可选，如果提供则只查询该申请人的
	status := r.URL.Query().Get("status")             // 可选，pending, approved, rejected

	requests, err := teamChatService.GetLeaveRequests(teamID, claims.Sub, applicantID, status)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, requests)
}

// ReviewLeaveRequest 审核请假申请
func ReviewLeaveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.ReviewLeaveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := teamChatService.ReviewLeaveRequest(&req, claims.Sub); err != nil {
		if err == service.ErrNotCaptain {
			util.Error(w, http.StatusForbidden, err.Error())
		} else if err == service.ErrLeaveNotFound {
			util.Error(w, http.StatusNotFound, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "reviewed successfully"})
}

