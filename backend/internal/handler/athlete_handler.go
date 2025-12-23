package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var athleteService = service.NewAthleteService()

// GetMyAthleteInfo 获取当前用户的运动员信息
func GetMyAthleteInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	athletes, err := athleteService.GetMyAthleteInfo(claims.Sub)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, athletes)
}

// GetMyAthleteInfoByTeam 获取当前用户在指定队伍中的运动员信息
func GetMyAthleteInfoByTeam(w http.ResponseWriter, r *http.Request) {
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

	athlete, err := athleteService.GetMyAthleteInfoByTeam(claims.Sub, teamID)
	if err != nil {
		if err == service.ErrAthleteNotFound || err == service.ErrNotTeamMember {
			util.Error(w, http.StatusNotFound, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, athlete)
}

// UpdateMyAthleteInfo 更新当前用户的运动员信息
func UpdateMyAthleteInfo(w http.ResponseWriter, r *http.Request) {
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
		TeamID       int64  `json:"team_id" binding:"required"`
		JerseyNumber string `json:"jersey_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := athleteService.UpdateMyAthleteInfo(claims.Sub, req.TeamID, req.JerseyNumber); err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "运动员信息更新成功"})
}

// LeaveTeam 退出队伍
func LeaveTeam(w http.ResponseWriter, r *http.Request) {
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
		TeamID int64 `json:"team_id" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := athleteService.LeaveTeam(claims.Sub, req.TeamID); err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "已成功退出队伍"})
}

// GetTeamMembersForAthlete 运动员查看队伍成员
func GetTeamMembersForAthlete(w http.ResponseWriter, r *http.Request) {
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

	members, err := athleteService.GetTeamMembersForAthlete(teamID, claims.Sub)
	if err != nil {
		if err == service.ErrNotTeamMember {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, members)
}

// CheckIfCaptain 检查是否是队长
func CheckIfCaptain(w http.ResponseWriter, r *http.Request) {
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

	isCaptain, err := athleteService.IsTeamCaptain(claims.Sub, teamID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]bool{"is_captain": isCaptain})
}

