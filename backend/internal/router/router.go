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
	mux.HandleFunc("/api/user/change-password", handler.ChangePassword) // 修改密码
	// 用户资料
	mux.HandleFunc("/api/user/profile/", handler.GetUserProfile) // expects /api/user/profile/{student_id}
	mux.HandleFunc("/api/user/avatar/", handler.UpdateAvatar)    // 更新用户头像 /api/user/avatar/{student_id}
	mux.HandleFunc("/api/user/college", handler.UpdateCollege)   // 更新用户学院信息

	// events （未实现）
	mux.HandleFunc("/api/events", handler.Events)
	mux.HandleFunc("/api/events/subscribe", handler.SubscribeEvent)

	// matches
	mux.HandleFunc("/api/matches", handler.Matches)
	mux.HandleFunc("/api/matches/", handler.MatchDetail) // expects /api/matches/{id}

	return mux
}
