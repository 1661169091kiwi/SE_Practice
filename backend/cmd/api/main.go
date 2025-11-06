package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "se_practice/backend/internal/router"
)

func main() {
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

    r := router.New()

    addr := ":" + port
    log.Printf("server starting at %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatalf("server failed: %v", err)
    }
    fmt.Println()
}


