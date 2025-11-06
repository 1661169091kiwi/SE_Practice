package handler

import (
    "net/http"

    "se_practice/backend/internal/util"
)

func Events(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, []interface{}{})
}

func SubscribeEvent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, map[string]bool{"subscribed": true})
}


