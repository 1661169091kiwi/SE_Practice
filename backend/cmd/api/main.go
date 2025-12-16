package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/router"
)

func main() {
	// 1. 读取端口
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatalf("DB_DSN is required")
	}
	if _, err := db.Init(dsn); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// 初始化表结构
	if err := db.InitSchema(); err != nil {
		log.Printf("Warning: failed to init schema: %v", err)
	}

	// 3. 注册路由并启动 HTTP 服务
	r := router.New()

	addr := ":" + port
	log.Printf("server starting at %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
	fmt.Println()
}
