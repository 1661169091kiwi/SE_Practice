package main

import (
	"multimodal-notes/controller"
	"multimodal-notes/service"
	"multimodal-notes/util"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	gptClient := util.NewGPTClient("sk-x") // 替换为你的 OpenAI API 密钥
	noteService := service.NewNoteService(gptClient)
	noteController := controller.NewNoteController(noteService)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/notes", noteController.UploadNote)
	}
}
