package router

import (
    "net/http"

    "se_practice/backend/internal/handler"
)

// New registers routes and returns an http.Handler.
func New() http.Handler {
    mux := http.NewServeMux()

    // health
    mux.HandleFunc("/api/health", handler.Health)

    // auth
    mux.HandleFunc("/api/register", handler.Register)
    mux.HandleFunc("/api/login", handler.Login)

    // events
    mux.HandleFunc("/api/events", handler.Events)
    mux.HandleFunc("/api/events/subscribe", handler.SubscribeEvent)

    // matches
    mux.HandleFunc("/api/matches", handler.Matches)
    mux.HandleFunc("/api/matches/", handler.MatchDetail) // expects /api/matches/{id}

    return mux
}


