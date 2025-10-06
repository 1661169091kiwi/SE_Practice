package route

import (
	"multimodal-notes/controller"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Note *controller.NoteController
	User *controller.UserController
	// 后续加新控制器（如标签、搜索）时，再在这里补充字段
}

func SetupRoutes(r *gin.Engine, c *Controllers) {
	// 1. API v1 分组：存放笔记相关接口
	v1 := r.Group("/api/v1")
	{
		notes := v1.Group("/notes")
		{
			notes.POST("/", c.Note.UploadNote)
		}
	}

	// 2. 存放用户注册/登录接口
	auth := r.Group("/auth")
	{
		auth.POST("/register", c.User.Register)
		auth.POST("/login", c.User.Login)
	}
}
