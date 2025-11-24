package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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

	// 调用服务层获取比赛列表
	var matches []model.Match
	if eventID > 0 {
		// 按赛事ID筛选
		matches, err = matchService.ListMatchesByEvent(eventID)
	} else {
		// 获取所有比赛
		matches, err = matchService.ListAllMatches()
	}

	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, matches)
}

// MatchDetail 获取比赛详情
func MatchDetail(w http.ResponseWriter, r *http.Request) {
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
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, match)
}

// UpdateMatchScore 更新比赛分数
func UpdateMatchScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径参数: /api/matches/{id}/score
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) < 2 || parts[0] == "" || parts[1] != "score" {
		util.Error(w, http.StatusBadRequest, "invalid url path, expected /api/matches/{id}/score")
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
