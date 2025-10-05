package main

import (
	"multimodal-notes/model"

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
	db.AutoMigrate(&model.Note{})

	r := gin.Default()
	SetupRoutes(r)
	r.Run(":8080")
}
