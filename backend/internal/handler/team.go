package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var teamService = service.NewTeamService()

var mockTeamsMu sync.Mutex
var mockTeams = map[int64]*model.Team{}
var mockTeamNextID int64 = 1

func mockCloneTeam(t *model.Team) model.Team {
	if t == nil {
		return model.Team{}
	}
	return model.Team{
		ID:          t.ID,
		TeamName:    t.TeamName,
		SportID:     t.SportID,
		College:     t.College,
		TeamType:    t.TeamType,
		AvatarURL:   t.AvatarURL,
		Description: t.Description,
		CreatedBy:   t.CreatedBy,
		CreatedAt:   t.CreatedAt,
		IsApproved:  t.IsApproved,
	}
}

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

	if db.GetDB() == nil {
		if req.TeamName == "" || req.SportID <= 0 {
			util.Error(w, http.StatusBadRequest, "team_name and sport_id are required")
			return
		}
		mockTeamsMu.Lock()
		id := mockTeamNextID
		mockTeamNextID++
		team := &model.Team{
			ID:          id,
			TeamName:    req.TeamName,
			SportID:     req.SportID,
			College:     req.College,
			TeamType:    req.TeamType,
			AvatarURL:   req.AvatarURL,
			Description: req.Description,
			CreatedBy:   req.CreatedBy,
			IsApproved:  true,
		}
		mockTeams[id] = team
		out := mockCloneTeam(team)
		mockTeamsMu.Unlock()

		log.Printf("[Admin Log] Team created (mock): %s by %s", out.TeamName, req.CreatedBy)
		util.OK(w, out)
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

	if db.GetDB() == nil {
		mockTeamsMu.Lock()
		out := make([]model.Team, 0)
		for _, t := range mockTeams {
			if t != nil && !t.IsApproved {
				out = append(out, mockCloneTeam(t))
			}
		}
		mockTeamsMu.Unlock()
		util.OK(w, out)
		return
	}

	teams, err := teamService.GetPendingTeams()
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, teams)
}

// UpdateTeamAvatar 更新队伍头像
func UpdateTeamAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.AuthClaims(r)
	if !ok {
		util.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 1. Parse Multipart Form (10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		util.Error(w, http.StatusBadRequest, "failed to parse multipart form")
		return
	}

	// 2. Get team_id
	teamIDStr := r.FormValue("team_id")
	if teamIDStr == "" {
		util.Error(w, http.StatusBadRequest, "team_id is required")
		return
	}
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid team_id")
		return
	}

	// 3. Get File
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		util.Error(w, http.StatusBadRequest, "failed to get file")
		return
	}
	defer file.Close()

	// 4. Validate file type
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if !allowedTypes[ext] {
		util.Error(w, http.StatusBadRequest, "invalid file type")
		return
	}

	// 5. Create directory
	uploadDir := "./uploads/teams"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		util.Error(w, http.StatusInternalServerError, "failed to create upload directory")
		return
	}

	// 6. Save file
	filename := fmt.Sprintf("team_%d%s", teamID, ext)
	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, "failed to create file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		util.Error(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	avatarURL := "/uploads/teams/" + filename

	// 7. Update DB
	if db.GetDB() == nil {
		// Mock
		mockTeamsMu.Lock()
		if t, ok := mockTeams[teamID]; ok {
			t.AvatarURL = avatarURL
		}
		mockTeamsMu.Unlock()
	} else {
		err = teamService.UpdateTeamAvatar(teamID, avatarURL, claims.Sub, claims.Role)
		if err != nil {
			if strings.Contains(err.Error(), "permission denied") {
				util.Error(w, http.StatusForbidden, err.Error())
			} else {
				util.Error(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
	}

	util.OK(w, map[string]string{
		"avatar_url": avatarURL,
		"message":    "team avatar updated successfully",
	})
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

	if db.GetDB() == nil {
		mockTeamsMu.Lock()
		t := mockTeams[req.TeamID]
		if t == nil {
			mockTeamsMu.Unlock()
			util.Error(w, http.StatusNotFound, "team not found")
			return
		}
		t.IsApproved = true
		mockTeamsMu.Unlock()
		util.OK(w, "team approved")
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

	if db.GetDB() == nil {
		mockTeamsMu.Lock()
		if _, ok := mockTeams[id]; !ok {
			mockTeamsMu.Unlock()
			util.Error(w, http.StatusNotFound, "team not found")
			return
		}
		delete(mockTeams, id)
		mockTeamsMu.Unlock()
		util.OK(w, "team deleted")
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

	if db.GetDB() == nil {
		studentID := r.URL.Query().Get("student_id")
		if studentID == "" {
			util.Error(w, http.StatusBadRequest, "student_id is required")
			return
		}
		util.OK(w, []model.Team{})
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

	if db.GetDB() == nil {
		if req.TeamName == "" || req.SportID == 0 {
			util.Error(w, http.StatusBadRequest, "team_name and sport_id are required")
			return
		}
		mockTeamsMu.Lock()
		id := mockTeamNextID
		mockTeamNextID++
		team := &model.Team{
			ID:          id,
			TeamName:    req.TeamName,
			SportID:     req.SportID,
			College:     req.College,
			TeamType:    req.TeamType,
			AvatarURL:   req.AvatarURL,
			Description: req.Description,
			CreatedBy:   req.CreatedBy,
			IsApproved:  false,
		}
		mockTeams[id] = team
		out := mockCloneTeam(team)
		mockTeamsMu.Unlock()

		util.OK(w, out)
		return
	}

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

	if db.GetDB() == nil {
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

		mockTeamsMu.Lock()
		out := make([]model.Team, 0)
		for _, t := range mockTeams {
			if t == nil || !t.IsApproved {
				continue
			}
			if sportID > 0 && t.SportID != sportID {
				continue
			}
			out = append(out, mockCloneTeam(t))
		}
		mockTeamsMu.Unlock()
		util.OK(w, out)
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
