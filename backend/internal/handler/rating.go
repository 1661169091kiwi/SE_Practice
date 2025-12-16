package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"se_practice/backend/internal/middleware"
	"se_practice/backend/internal/service"
	"se_practice/backend/internal/util"
)

var ratingService = service.NewRatingService()

func Ratings(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
		if len(parts) < 2 || parts[1] != "ratings" {
			util.Error(w, http.StatusBadRequest, "invalid path")
			return
		}
		matchID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			util.Error(w, http.StatusBadRequest, "invalid match id")
			return
		}
		list, err := ratingService.ListByMatch(matchID)
		if err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, list)
		return
	}
	if r.Method == http.MethodPost {
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
		if len(parts) < 2 || parts[1] != "ratings" {
			util.Error(w, http.StatusBadRequest, "invalid path")
			return
		}
		matchID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			util.Error(w, http.StatusBadRequest, "invalid match id")
			return
		}
		claims, ok := middleware.AuthClaims(r)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var body struct {
			StudentID string  `json:"student_id"`
			RaterID   string  `json:"rater_id"`
			Score     float64 `json:"score"`
			Comment   string  `json:"comment"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			util.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if body.StudentID == "" || body.RaterID == "" {
			util.Error(w, http.StatusBadRequest, "student_id and rater_id are required")
			return
		}
		if claims.Sub != body.RaterID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if body.Score < 0 {
			util.Error(w, http.StatusBadRequest, "score invalid")
			return
		}
		if err := ratingService.Upsert(matchID, body.StudentID, body.RaterID, body.Score, body.Comment); err != nil {
			util.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		util.OK(w, map[string]any{"match_id": matchID})
		return
	}
	util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
}
