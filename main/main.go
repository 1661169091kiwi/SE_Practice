package main

import (
	"multimodal-notes/controller"
	"multimodal-notes/model"
	"multimodal-notes/route"
	"multimodal-notes/service"
	"multimodal-notes/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:qqhyjy0514@tcp(localhost:3306)/multimodal_notes?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	db.AutoMigrate(&model.Note{}, &model.User{})

	//API
	gptClient := util.NewGPTClient("sk-gLJpO9I10cMfjn0PYz80SSwELl84fmTyKjhYlMUwkyANTfpf")

	//note Service和Controller
	noteService := service.NewNoteService(gptClient)
	noteController := controller.NewNoteController(noteService)
	//user Service和Controller
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	r := gin.Default()

	controllers := &route.Controllers{
		Note: noteController,
		User: userController,
		//后面的功能在这里继续添加
	}
	route.SetupRoutes(r, controllers)
	r.Run(":8080")
}
