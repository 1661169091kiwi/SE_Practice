package main

import (
	"multimodal-notes/controller"
	"multimodal-notes/model"
	"multimodal-notes/repository"
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
	db.AutoMigrate(&model.Note{}, &model.User{}, &model.Collection{}, &model.Tag{})

	//API
	gptClient := util.NewGPTClient("sk-gLJpmTyKjhYlMUwkyANTfpf")

	//note Service和Controller
	noteService := service.NewNoteService(gptClient)
	noteController := controller.NewNoteController(noteService)
	//user Service和Controller
	userService := service.NewUserService(db)
	userController := controller.NewUserController(userService)

	//保存发表note
	noteRepo := repository.NewNoteRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// 3. 初始化收藏相关依赖（新增）
	// 3.1 收藏仓库
	colRepo := repository.NewCollectionRepository(db)
	// 3.3 收藏服务（注入收藏仓库和笔记仓库）
	colService := service.NewCollectionService(colRepo, noteRepo)
	// 3.4 收藏控制器
	colController := controller.NewCollectionController(colService)

	saveNoteService := service.NewSaveNoteService(noteRepo, tagRepo)
	saveNoteController := controller.NewSaveNoteController(saveNoteService)

	r := gin.Default()

	controllers := &route.Controllers{
		Note:       noteController,
		User:       userController,
		SaveNote:   saveNoteController,
		Collection: colController,

		//后面的功能在这里继续添加
	}
	route.SetupRoutes(r, controllers)
	r.Run(":8080")
}
