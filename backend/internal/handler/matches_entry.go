package handler

import (
	"net/http"
	"strings"
)

func MatchesEntry(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/api/matches/")
	if strings.HasSuffix(p, "/ratings") || strings.Contains(p, "/ratings") {
		Ratings(w, r)
		return
	}
	if strings.HasSuffix(p, "/lineups") || strings.Contains(p, "/lineups") {
		Lineups(w, r)
		return
	}
	if strings.HasSuffix(p, "/comments") || strings.Contains(p, "/comments") {
		MatchComments(w, r)
		return
	}
	if strings.HasSuffix(p, "/stats") || strings.Contains(p, "/stats") {
		MatchStats(w, r)
		return
	}
	MatchDetail(w, r)
}
