package handler

import (
    "net/http"

    "se_practice/backend/internal/util"
)

func Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, map[string]string{"result": "register stub"})
}

func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        util.Error(w, http.StatusMethodNotAllowed, "method not allowed")
        return
    }
    util.OK(w, map[string]string{"token": "stub.jwt.token"})
}


