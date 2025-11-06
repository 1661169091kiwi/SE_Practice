package handler

import (
    "net/http"
    "strings"

    "se_practice/backend/internal/util"
)

func Matches(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, []interface{}{})
}

func MatchDetail(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    // naive id parse: /api/matches/{id}
    parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/matches/"), "/")
    id := ""
    if len(parts) > 0 {
        id = parts[0]
    }
    util.OK(w, map[string]string{"id": id})
}


