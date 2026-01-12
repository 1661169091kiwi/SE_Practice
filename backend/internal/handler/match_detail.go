package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/util"
)

// MatchComments 处理 /api/matches/{id}/comments
func MatchComments(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) < 2 || parts[1] != "comments" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}
	matchID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}

	if r.Method == http.MethodGet {
		if db.GetDB() == nil {
			util.OK(w, []model.MatchComment{})
			return
		}
		comments, err := matchService.GetMatchComments(matchID)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if comments == nil {
			comments = []model.MatchComment{}
		}
		util.OK(w, comments)
		return
	}

	if r.Method == http.MethodPost {
		if db.GetDB() == nil {
			util.OK(w, map[string]string{"message": "comment created successfully"})
			return
		}
		claims, ok := middleware.AuthClaims(r)
		if !ok {
			util.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		studentID := claims.Sub
		var req model.CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		req.MatchID = matchID
		if req.Content == "" {
			util.Error(w, http.StatusBadRequest, "content is required")
			return
		}
		err := matchService.CreateMatchComment(&req, studentID)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, map[string]string{"message": "comment created successfully"})
		return
	}
	util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
}

// MatchStats 处理 /api/matches/{id}/stats
func MatchStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) < 2 || parts[1] != "stats" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}
	matchID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}
	if db.GetDB() == nil {
		util.OK(w, &model.MatchStats{
			MatchID:           matchID,
			HomePossession:    50,
			AwayPossession:    50,
			HomeShots:         0,
			AwayShots:         0,
			HomeShotsOnTarget: 0,
			AwayShotsOnTarget: 0,
			HomeFouls:         0,
			AwayFouls:         0,
			HomeCorners:       0,
			AwayCorners:       0,
		})
		return
	}
	stats, err := matchService.GetMatchStats(matchID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, stats)
}

// GetMatchLineups 保留查询比赛阵容（通过 query 的旧接口）
// 新接口请使用 /api/matches/{id}/lineups 由 Lineups 处理
