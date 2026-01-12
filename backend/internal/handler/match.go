package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var matchService = service.NewMatchService()

// Matches 获取比赛列表（支持按赛事ID筛选）
func Matches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if db.GetDB() == nil {
		now := time.Now()
		util.OK(w, []model.Match{
			{
				ID:        1,
				EventID:   1,
				SportID:   1,
				Name:      "示例比赛 A",
				Round:     "小组赛",
				Time:      now.Add(2 * time.Hour),
				TeamAID:   101,
				TeamAName: "示例主队",
				TeamBID:   102,
				TeamBName: "示例客队",
				ScoreA:    0,
				ScoreB:    0,
				Status:    "not_started",
			},
			{
				ID:        2,
				EventID:   1,
				SportID:   1,
				Name:      "示例比赛 B",
				Round:     "小组赛",
				Time:      now.Add(-30 * time.Minute),
				TeamAID:   103,
				TeamAName: "示例主队 2",
				TeamBID:   104,
				TeamBName: "示例客队 2",
				ScoreA:    1,
				ScoreB:    0,
				Status:    "ongoing",
			},
		})
		return
	}

	// 从查询参数获取赛事ID（可选）
	eventIDStr := r.URL.Query().Get("event_id")
	var eventID int64 = 0
	var err error
	if eventIDStr != "" {
		eventID, err = strconv.ParseInt(eventIDStr, 10, 64)
		if err != nil {
			util.Error(w, http.StatusBadRequest, "invalid event_id parameter")
			return
		}
	}

	// 获取当前用户ID（可选，用于检查订阅状态）
	studentID := r.URL.Query().Get("student_id")

	// 调用服务层获取比赛列表
	var matches []model.Match
	view := r.URL.Query().Get("view")
	if eventID > 0 {
		// 按赛事ID筛选
		matches, err = matchService.ListMatchesByEvent(eventID, view, studentID)
	} else {
		// 获取所有比赛
		matches, err = matchService.ListAllMatches(view, studentID)
	}

	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, matches)
}

// AthleteMatches 获取运动员参加的比赛列表
func AthleteMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 从查询参数获取学号
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	// 获取 view 参数（可选）
	view := r.URL.Query().Get("view")

	// 调用服务层获取运动员的比赛列表
	matches, err := matchService.GetAthleteMatches(studentID, view)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, matches)
}

// AvailableMatchesForAthlete 获取运动员可以报名的比赛列表
func AvailableMatchesForAthlete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 从查询参数获取学号
	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	// 获取 view 参数（可选）
	view := r.URL.Query().Get("view")

	// 调用服务层获取可报名的比赛列表
	matches, err := matchService.GetAvailableMatchesForAthlete(studentID, view)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, matches)
}

// JoinMatch 运动员报名参加比赛
func JoinMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.JoinMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if req.StudentID == "" || req.MatchID == 0 || req.TeamID == 0 {
		util.Error(w, http.StatusBadRequest, "student_id, match_id and team_id are required")
		return
	}

	// 调用服务层报名
	if err := matchService.JoinMatch(&req); err != nil {
		util.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "报名成功"})
}

// MatchDetail 获取比赛详情
func MatchDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		UpdateMatch(w, r)
		return
	}
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径参数: /api/matches/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "match id is required")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}

	if db.GetDB() == nil {
		now := time.Now()
		detail := &model.MatchDetailResponse{
			MatchID:   id,
			EventID:   1,
			MatchName: "示例比赛",
			Round:     "小组赛",
			MatchTime: now.Add(2 * time.Hour),
			TeamA:     model.TeamBrief{ID: 101, Name: "示例主队"},
			TeamB:     model.TeamBrief{ID: 102, Name: "示例客队"},
			ScoreA:    0,
			ScoreB:    0,
			Status:    "not_started",
			Collectors: []model.UserBrief{
				{ID: 0},
				{ID: 0},
			},
		}
		util.OK(w, detail)
		return
	}

	// 调用服务层获取比赛详情
	detail, err := matchService.GetMatchDetail(id)
	if err != nil {
		if err.Error() == "match not found" {
			util.Error(w, http.StatusNotFound, "match not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, detail)
}

// CreateMatch 创建比赛
func CreateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req model.CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 基础参数验证
	if req.EventID <= 0 || req.TeamAID <= 0 || req.TeamBID <= 0 || req.MatchName == "" {
		util.Error(w, http.StatusBadRequest, "event_id, team_a_id, team_b_id and match_name are required")
		return
	}

	match, err := matchService.CreateMatch(&req)
	if err != nil {
		switch err.Error() {
		case "event not found":
			util.Error(w, http.StatusNotFound, "event not found")
		case "team a not found", "team b not found":
			util.Error(w, http.StatusNotFound, err.Error())
		case "team a is not approved", "team b is not approved", "event teams not configured", "team is not in event":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	log.Printf("[Admin Log] Match created: %s (EventID: %d)", match.Name, match.EventID)

	util.OK(w, match)
}

// UpdateMatchScore 更新比赛分数
func UpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径参数: /api/matches/update-score/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/update-score/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "match id is required")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}

	// 解析请求体中的分数
	var scoreReq struct {
		ScoreA int64 `json:"score_team_a"`
		ScoreB int64 `json:"score_team_b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&scoreReq); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 更新分数
	err = matchService.UpdateMatchScore(id, scoreReq.ScoreA, scoreReq.ScoreB)
	if err != nil {
		switch err.Error() {
		case "match not found":
			util.Error(w, http.StatusNotFound, "match not found")
		case "scores cannot be negative":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "score updated successfully"})
}

// UpdateMatch 更新比赛信息
func UpdateMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	// 解析路径参数: /api/matches/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "match id is required")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}

	var req model.CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = matchService.UpdateMatch(id, &req)
	if err != nil {
		switch err.Error() {
		case "match not found", "event not found":
			util.Error(w, http.StatusNotFound, err.Error())
		case "team a not found", "team b not found", "team a is not approved", "team b is not approved", "event teams not configured", "team is not in event":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "match updated successfully"})
}

// GetMatchData 获取赛事数据（积分榜/球员榜/赛程/历史）
func GetMatchData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析查询参数
	q := r.URL.Query()
	sportTypeStr := q.Get("sportType")
	matchTime := q.Get("matchTime")
	dataType := q.Get("dataType")

	if sportTypeStr == "" || matchTime == "" || dataType == "" {
		util.Error(w, http.StatusBadRequest, "sportType, matchTime, and dataType are required")
		return
	}

	sportType, err := strconv.Atoi(sportTypeStr)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid sportType")
		return
	}

	// 调用服务层
	data, err := matchService.GetMatchData(sportType, matchTime, dataType)
	if err != nil {
		if err.Error() == "invalid dataType" {
			util.Error(w, http.StatusBadRequest, "invalid dataType")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, data)
}

// GetUserSubscribedMatches 获取用户订阅的比赛
func GetUserSubscribedMatches(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userId := r.URL.Query().Get("userId")
	if userId == "" {
		util.Error(w, http.StatusBadRequest, "userId is required")
		return
	}

	data, err := matchService.GetUserSubscribedMatches(userId)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, data)
}

// SubscribeMatch 订阅/取消订阅比赛
func SubscribeMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req model.SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.StudentID == "" || (req.MatchID == "" && req.EventID == 0) {
		util.Error(w, http.StatusBadRequest, "studentId and (matchId or eventId) are required")
		return
	}

	msg, err := matchService.SubscribeMatch(&req)
	if err != nil {
		// 判断错误类型返回404或500
		if err.Error() == "match not found" {
			util.Error(w, http.StatusNotFound, err.Error())
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			// 如果是重复订阅错误，视为成功
			resp := model.SubscribeResponse{
				StudentID:   req.StudentID,
				MatchID:     req.MatchID,
				OperateType: req.OperateType,
			}
			util.JSON(w, http.StatusOK, util.APIResponse{
				Code:    200,
				Message: "已订阅该比赛",
				Data:    resp,
			})
			return
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := model.SubscribeResponse{
		StudentID:   req.StudentID,
		MatchID:     req.MatchID,
		OperateType: req.OperateType,
	}

	// 自定义返回结构以包含 message，需导入 util.APIResponse 结构体，
	// 若 response.go 中未导出 APIResponse，则需修改 util 或直接 map
	// 假设 response.go 中 APIResponse 已导出
	util.JSON(w, http.StatusOK, util.APIResponse{
		Code:    200,
		Message: msg,
		Data:    resp,
	})
}

// DeleteMatch 删除比赛
func DeleteMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		util.Error(w, http.StatusBadRequest, "id required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := matchService.DeleteMatch(id); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, "match deleted")
}
