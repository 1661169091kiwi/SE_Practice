package handler

import (
	"net/http"
	"strconv"
	"strings"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var standingsService = service.NewStandingsService()

// GetEventStandings 获取赛事积分榜
func GetEventStandings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 路径应为 /api/events/{event_id}/standings
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) < 2 || parts[1] != "standings" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	eventID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if db.GetDB() == nil {
		util.OK(w, []model.Standings{})
		return
	}

	standings, err := standingsService.GetStandings(eventID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, standings)
}

func GetEventStandingsOverview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/events/"), "/")
	if len(parts) < 3 || parts[1] != "standings" || parts[2] != "overview" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}

	eventID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if db.GetDB() == nil {
		util.OK(w, map[string]any{
			"event_id": eventID,
			"format":   "points",
		})
		return
	}

	ov, err := standingsService.GetStandingsOverview(eventID)
	if err != nil {
		if err.Error() == "event not found" {
			util.Error(w, http.StatusNotFound, "event not found")
			return
		}
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, ov)
}
