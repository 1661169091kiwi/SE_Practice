package handler

import (
    "net/http"

    "se_practice/backend/internal/util"
)

func Health(w http.ResponseWriter, r *http.Request) {
    util.OK(w, map[string]string{"status": "ok"})
}


