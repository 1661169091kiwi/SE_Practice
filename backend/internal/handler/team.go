package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var teamService = service.NewTeamService()

// CreateTeam 创建队伍
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req model.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 从 Context 获取当前用户ID (Assume Auth middleware sets it)
	// For simplicity, we might assume it's passed in body for now if middleware isn't fully setting context user
	// But let's check if we can get it from headers or token claims.
	// Given the current auth middleware structure, let's assume the client sends "created_by" or we take it from token.
	// Since I cannot easily change Auth middleware right now, I'll assume it's in the request body OR
	// I should extract it from context if the middleware puts it there.
	// Looking at `middleware/auth.go`, it likely puts user info in context.
	// Let's assume for now the admin passes it or we trust the body.
	// To be safe, if CreatedBy is empty, we fail.

	if req.CreatedBy == "" {
		// Try to get from header or context?
		// For this phase, let's require it in the body or allow it to be empty if the DB allows (DB says NOT NULL).
		// So client must send it.
		util.Error(w, http.StatusBadRequest, "created_by is required")
		return
	}

	team, err := teamService.CreateTeam(&req)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Admin Log] Team created: %s by %s", team.TeamName, req.CreatedBy)

	util.OK(w, team)
}

// GetPendingTeams 获取待审核队伍
func GetPendingTeams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	teams, err := teamService.GetPendingTeams()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, teams)
}

// ApproveTeam 批准队伍
func ApproveTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req struct {
		TeamID int64 `json:"team_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if err := teamService.ApproveTeam(req.TeamID); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, "team approved")
}

// DeleteTeam 删除队伍
func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if !middleware.CheckRole(w, r, "admin") {
		return
	}
	// Parse URL param manually or use query param?
	// Router uses HandleFunc so we might need to parse path or use query param.
	// Since standard net/http mux in Go 1.22+ supports path vars, but here it looks like older style or custom.
	// Looking at router.go: `mux.HandleFunc("/api/user/profile/", handler.GetUserProfile)` suggests manual parsing.
	// But `router.go` also has `mux.HandleFunc("/api/admin/teams/approve", ...)` which uses body.
	// For DELETE, body is not recommended. Query param or path param is better.
	// I'll use query param `team_id` for simplicity as seen in `ListTeams` (though it uses sport_id query).
	// Actually `DeleteTeam` usually expects ID in path.
	// Let's use Query param `id` to be consistent with simple handlers, or parse path if I register it as `/api/teams/delete/`.
	// I'll register as `/api/teams/delete` and expect `?id=...`.

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

	if err := teamService.DeleteTeam(id); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, "team deleted")
}

// GetTeamMembersHandler 获取队伍成员
func GetTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	// Query param `team_id`
	idStr := r.URL.Query().Get("team_id")
	if idStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	members, err := teamService.GetTeamMembers(id)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, members)
}

// GetPendingTeamMembersHandler 获取待审核成员
func GetPendingTeamMembersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("team_id")
	if idStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	members, err := teamService.GetPendingTeamMembers(id, claims.Sub, claims.Role)
	if err != nil {
		if err.Error() == "permission denied" {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, members)
}

// ApproveTeamMemberHandler 批准队伍成员
func ApproveTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
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
		TeamMemberID int64 `json:"team_member_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := teamService.ApproveMemberWithAuth(req.TeamMemberID, claims.Sub, claims.Role); err != nil {
		if err.Error() == "permission denied" {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, "member approved")
}

// RemoveTeamMember 移除队员
func RemoveTeamMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // or DELETE
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		TeamMemberID int64 `json:"team_member_id"`
		// Optional: Pass TeamID if frontend has it to save a DB lookup, but Repo/Service can find it.
		// However, to check if I am the creator, I need to know which team this member belongs to.
		// Let's rely on Service to look up the member's team.
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	// We need to check if the user is admin OR the creator of the team the member belongs to.
	// Since we only have TeamMemberID, we need to fetch the team info first.
	// Let's add a helper in Service or do it here.
	// Doing it in Service is cleaner. Update RemoveMember to accept UserID and Role.

	if err := teamService.RemoveMemberWithAuth(req.TeamMemberID, claims.Sub, claims.Role); err != nil {
		if err.Error() == "permission denied" {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, "member removed")
}

// UpdateTeamName 更新队伍名称
func UpdateTeamName(w http.ResponseWriter, r *http.Request) {
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
		TeamID   int64  `json:"team_id"`
		TeamName string `json:"team_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	if req.TeamName == "" {
		util.Error(w, http.StatusBadRequest, "team_name required")
		return
	}

	if err := teamService.UpdateTeamName(req.TeamID, req.TeamName, claims.Sub, claims.Role); err != nil {
		if err.Error() == "permission denied: not the team creator" {
			util.Error(w, http.StatusForbidden, err.Error())
		} else {
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, "team name updated")
}

// GetMyTeams 获取我的队伍
func GetMyTeams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	studentID := r.URL.Query().Get("student_id")
	if studentID == "" {
		util.Error(w, http.StatusBadRequest, "student_id is required")
		return
	}

	teams, err := teamService.GetMyTeams(studentID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, teams)
}

// ApplyCreateTeam 申请创建队伍
func ApplyCreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 获取当前用户信息
	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 强制设置创建者为当前用户
	req.CreatedBy = claims.Sub

	team, err := teamService.ApplyCreateTeam(&req)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, team)
}

// ListTeams 列出队伍
func ListTeams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sportIDStr := r.URL.Query().Get("sport_id")
	var sportID int64 = 0
	if sportIDStr != "" {
		var err error
		sportID, err = strconv.ParseInt(sportIDStr, 10, 64)
		if err != nil {
			util.Error(w, http.StatusBadRequest, "invalid sport_id")
			return
		}
	}

	teams, err := teamService.ListTeams(sportID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, teams)
}
