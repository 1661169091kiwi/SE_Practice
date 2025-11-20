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

	// 2. 初始化数据库（使用 configs/config.yaml 里的默认 DSN，实际可改为读配置）
	dsn := "user:pass@tcp(localhost:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local"
	if envDSN := os.Getenv("DB_DSN"); envDSN != "" {
		dsn = envDSN
	}
	if _, err := db.Init(dsn); err != nil {
		log.Fatalf("failed to init database: %v", err)
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


