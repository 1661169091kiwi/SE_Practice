package middleware

import "net/http"

// Auth is a minimal no-op middleware placeholder.
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)
    })
}


