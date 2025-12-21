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
	if req.Format != "points" && req.Format != "group_knockout" {
		util.Error(w, http.StatusBadRequest, "invalid format_type: must be 'points' or 'group_knockout'")
		return
	}

	event, err := eventService.CreateEvent(&req)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[Admin Log] Event created: %s", event.EventName)

	util.OK(w, event)
}

// GetEventDetail 获取赛事详情
func GetEventDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		UpdateEvent(w, r)
		return
	}
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径参数: /api/events/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		util.Error(w, http.StatusBadRequest, "event id is required")
		return
	}

	// 检查是否请求积分榜
	if len(parts) >= 2 && parts[1] == "standings" {
		GetEventStandings(w, r)
		return
	}

	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
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
		util.Error(w, http.StatusInternalServerError, err.Error())
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
