package handler

import (
	"net/http"
	"strconv"
	"strings"

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

	standings, err := standingsService.GetStandings(eventID)
	if err != nil {
		util.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.OK(w, standings)
}
