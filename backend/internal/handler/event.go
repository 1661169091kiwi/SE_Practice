package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var eventService = service.NewEventService()

func Events(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		status := r.URL.Query().Get("status")
		events, err := eventService.ListEvents(status)
		if err != nil {
			if err == service.ErrInvalidStatus {
				util.Error(w, http.StatusBadRequest, "invalid status parameter")
				return
			}
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, events)
		return
	}

	util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
}

// CreateEvent 创建赛事
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	var req model.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 基础参数验证
	if req.Format != "points" && req.Format != "knockout" && req.Format != "group_knockout" {
		util.Error(w, http.StatusBadRequest, "invalid format_type: must be 'points' or 'knockout' or 'group_knockout'")
		return
	}

	event, err := eventService.CreateEvent(&req)
	if err != nil {
		switch err.Error() {
		case "event name and sport id are required", "start date must be before end date", "invalid format type: must be 'points' or 'knockout' or 'group_knockout'",
			"teams are required", "invalid team_id", "duplicate team_id", "team not found", "team is not approved", "team sport_id mismatch",
			"invalid event format", "not enough teams", "knockout slot required", "duplicate slot", "slots must be 1..N without gaps", "team count must be power of two":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	log.Printf("[Admin Log] Event created: %s", event.EventName)

	util.OK(w, event)
}

// GetEventDetail 获取赛事详情
func GetEventDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPut {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径参数: /api/events/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "event id is required")
		return
	}

	if len(parts) >= 2 && parts[1] == "teams" {
		if r.Method == http.MethodPut {
			ReplaceEventTeams(w, r, parts[0])
			return
		}
		ListEventTeams(w, r, parts[0])
		return
	}

	if len(parts) >= 3 && parts[1] == "standings" && parts[2] == "overview" {
		if r.Method != http.MethodGet {
			util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		GetEventStandingsOverview(w, r)
		return
	}
	if len(parts) >= 2 && parts[1] == "standings" {
		if r.Method != http.MethodGet {
			util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		GetEventStandings(w, r)
		return
	}

	if len(parts) >= 2 && parts[1] == "knockout-schedule" {
		UpdateKnockoutSchedule(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if r.Method == http.MethodPut {
		UpdateEvent(w, r)
		return
	}

	event, err := eventService.GetEventDetail(id)
	if err != nil {
		if err == service.ErrEventNotFound {
			util.Error(w, http.StatusNotFound, "event not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, event)
}

func ListEventTeams(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	id, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}
	list, err := eventService.ListEventTeams(id)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.OK(w, list)
}

func ReplaceEventTeams(w http.ResponseWriter, r *http.Request, eventIDStr string) {
	if !middleware.CheckRole(w, r, "admin") {
		return
	}
	id, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}
	var body struct {
		Teams []model.EventTeamInput `json:"teams"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := eventService.ReplaceEventTeams(id, body.Teams); err != nil {
		if err == service.ErrEventNotFound {
			util.Error(w, http.StatusNotFound, "event not found")
			return
		}
		switch err.Error() {
		case "teams are required", "invalid team_id", "duplicate team_id", "team not found", "team is not approved", "team sport_id mismatch",
			"invalid event format", "not enough teams", "knockout slot required", "duplicate slot", "slots must be 1..N without gaps", "team count must be power of two":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	util.OK(w, map[string]string{"message": "event teams updated successfully"})
}

// UpdateEvent 更新赛事信息
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 权限检查
	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	// 解析路径参数: /api/events/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "event id is required")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req model.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = eventService.UpdateEvent(id, &req)
	if err != nil {
		if err == service.ErrEventNotFound {
			util.Error(w, http.StatusNotFound, "event not found")
			return
		}
		switch err.Error() {
		case "event name and sport id are required", "start date must be before end date", "invalid format type: must be 'points' or 'knockout' or 'group_knockout'",
			"teams are required", "invalid team_id", "duplicate team_id", "team not found", "team is not approved", "team sport_id mismatch",
			"invalid event format", "not enough teams", "knockout slot required", "duplicate slot", "slots must be 1..N without gaps", "team count must be power of two":
			util.Error(w, http.StatusBadRequest, err.Error())
		default:
			util.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	util.OK(w, map[string]string{"message": "event updated successfully"})
}

// ListEvents 列出赛事
func ListEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 获取查询参数 status
	status := r.URL.Query().Get("status")
	events, err := eventService.ListEvents(status)
	if err != nil {
		if err == service.ErrInvalidStatus {
			util.Error(w, http.StatusBadRequest, "invalid status parameter")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, events)
}

// UpdateKnockoutSchedule 更新淘汰赛赛程
func UpdateKnockoutSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) < 3 || parts[1] != "knockout-schedule" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}
	// eventID is not strictly needed for the repo method as it uses knockout_match_id, but good for validation if needed
	// _, err := strconv.ParseInt(parts[0], 10, 64)

	var req map[string]string // map[knockout_match_id_string]time_string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updates := make(map[int64]string)
	for k, v := range req {
		id, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		updates[id] = v
	}

	eventRepo := service.NewEventService().Repo()
	if err := eventRepo.UpdateKnockoutSchedule(updates); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "schedule updated"})
}

// UpdateEventStatus 更新赛事状态
func UpdateEventStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !middleware.CheckRole(w, r, "admin") {
		return
	}

	// path: /api/events/{id}/status
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) < 2 {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := eventService.UpdateEventStatus(id, req.Status); err != nil {
		if err == service.ErrEventNotFound {
			util.Error(w, http.StatusNotFound, "event not found")
			return
		}
		if err == service.ErrInvalidStatus {
			util.Error(w, http.StatusBadRequest, "invalid status")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, map[string]string{"message": "event status updated"})
}

// DeleteEvent 删除赛事
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
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

	if err := eventService.DeleteEvent(id); err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, "event deleted")
}
