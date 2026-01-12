package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	// 请修改下面的账号密码
	dsn := "root:SYSU@tcp(127.0.0.1:3306)/test_for_project?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("数据库连接成功！")
	return db
}
