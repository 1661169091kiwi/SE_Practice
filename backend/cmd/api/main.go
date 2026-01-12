package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"se_practice/backend/internal/db"
	"se_practice/backend/internal/model"
	"se_practice/backend/internal/repo"
	"se_practice/backend/internal/router"
)

func main() {
	// 1. 读取端口
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		if _, err := db.Init(dsn); err != nil {
			log.Fatalf("Error: failed to init database: %v", err)
		} else {
			if err := db.InitSchema(); err != nil {
				log.Printf("Warning: failed to init schema: %v", err)
			}
			if err := ensureAdminAccount(); err != nil {
				log.Printf("Warning: failed to ensure admin account: %v", err)
			}
		}
	} else {
		log.Printf("Warning: DB_DSN is not set, starting without database")
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

func ensureAdminAccount() error {
	if db.GetDB() == nil {
		return nil
	}

	studentID := "admin"
	password := "admin123"

	userRepo := repo.NewUserRepo()
	exists, err := userRepo.UserExists(studentID)
	if err != nil {
		return err
	}
	if !exists {
		sum := sha256.Sum256([]byte(password))
		hashed := hex.EncodeToString(sum[:])
		if err := userRepo.CreateUser(&model.User{
			StudentID: studentID,
			Password:  hashed,
			Name:      "管理员",
			College:   "系统管理部",
			Grade:     "system",
			AvatarURL: "",
		}); err != nil {
			return err
		}
	}

	admin, err := userRepo.GetAdminByStudentID(studentID)
	if err != nil {
		return err
	}
	if admin == nil {
		if _, err := userRepo.CreateAdmin(&model.Admin{
			StudentID:   studentID,
			Role:        "admin",
			Permissions: "",
		}); err != nil {
			return err
		}
	}

	return nil
}
