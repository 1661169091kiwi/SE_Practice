package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var lineupService = NewLineupServiceProxy()

type lineupProxy struct {
	s *service.LineupService
}

func NewLineupServiceProxy() *lineupProxy {
	return &lineupProxy{s: service.NewLineupService()}
}

func (p *lineupProxy) Upsert(items []model.MatchLineup) error {
	return p.s.Upsert(items)
}

func (p *lineupProxy) ListByMatch(matchID int64) ([]model.MatchLineup, error) {
	return p.s.ListByMatch(matchID)
}

func Lineups(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
	if len(parts) < 2 || parts[1] != "lineups" {
		util.Error(w, http.StatusBadRequest, "invalid path")
		return
	}
	matchID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		util.Error(w, http.StatusBadRequest, "invalid match id")
		return
	}
	if r.Method == http.MethodGet {
		list, err := lineupService.ListByMatch(matchID)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, list)
		return
	}
	if r.Method == http.MethodPost {
		// 仅 POST 写入阵容需要采集员角色
		if !middleware.CheckRole(w, r, "collector") {
			return
		}
		var body struct {
			Items []model.MatchLineup `json:"items"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		for i := range body.Items {
			body.Items[i].MatchID = matchID
		}
		if err := lineupService.Upsert(body.Items); err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, map[string]any{"match_id": matchID, "updated": len(body.Items)})
		return
	}
	if r.Method == http.MethodDelete {
		// 仅 DELETE 需要采集员角色
		if !middleware.CheckRole(w, r, "collector") {
			return
		}
		var body struct {
			TeamID    int64  `json:"teamId"`
			StudentID string `json:"studentId"`
			Name      string `json:"name"`
			Position  string `json:"position"`
			Number    string `json:"number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if body.TeamID == 0 {
			util.Error(w, http.StatusBadRequest, "teamId required")
			return
		}
		lrepo := repo.NewLineupRepo()
		tx, err := db.BeginTransaction()
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		if strings.TrimSpace(body.StudentID) != "" {
			if err := lrepo.DeleteTx(tx, matchID, body.TeamID, strings.TrimSpace(body.StudentID)); err != nil {
				_ = tx.Rollback()
				util.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			if strings.TrimSpace(body.Name) == "" {
				_ = tx.Rollback()
				util.Error(w, http.StatusBadRequest, "name required for external deletion")
				return
			}
			if err := lrepo.DeleteExternalTx(tx, matchID, body.TeamID, strings.TrimSpace(body.Name), strings.TrimSpace(body.Position), strings.TrimSpace(body.Number)); err != nil {
				_ = tx.Rollback()
				util.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if err := tx.Commit(); err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, map[string]any{"match_id": matchID, "deleted": true})
		return
	}
	util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
}
