package main

import (
	"fmt"

	// 引入配置和路由包
	"my_project/config"
	"my_project/router"
	// 如果你想让 GORM 自动帮你建表，可以引入 model 包
	// "my_project/model"
)

func main() {
	// 1. 初始化数据库连接
	// 程序启动时先连上冰箱(数据库)，连不上直接报错退出
	db := config.InitDB()

	// 【可选】自动迁移数据库表结构
	// 这行代码会检查 model 里的 Struct 和数据库里的表是不是对得上的
	// 如果数据库里没表，它会自动创建。如果你已经手动导入了 SQL，这行可以注释掉。
	// db.AutoMigrate(&model.User{}, &model.Sport{}, &model.Event{}, &model.Team{}, &model.Match{}, &model.Athlete{})

	// 2. 初始化路由
	// 把连好的 db 传给路由层
	r := router.SetupRouter(db)

	// 3. 启动 Web 服务
	// 默认监听 8080 端口
	port := ":8080"
	fmt.Printf("服务启动成功! 接口地址: http://localhost%s/api/...\n", port)

	// Run() 会阻塞在这里一直运行，直到你按 Ctrl+C 停止
	if err := r.Run(port); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
