package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// 引入你的 controller 包
	"my_project/controller"
)

// SetupRouter 初始化路由
func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 1. 配置 CORS (跨域中间件)
	// 这是为了防止前端(比如 localhost:8080) 访问 后端(localhost:3000) 时被浏览器拦截
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许任何来源
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 浏览器会在发送 POST 前发一个 OPTIONS 请求试探，直接返回成功即可
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 2. 初始化控制器
	// 把数据库连接传给 AdminController，这样所有接口都能用 db 了
	adminCtrl := &controller.AdminController{DB: db}

	// 3. 注册路由
	// 所有接口都以 /api 开头
	api := r.Group("/api")
	{
		// --- 赛事相关接口 ---
		api.POST("/publish-event", adminCtrl.PublishEvent) // 发布赛事
		api.POST("/arrange-event", adminCtrl.ArrangeEvent) // 赛事编排
		api.POST("/save-grade", adminCtrl.SaveGrade)       // 成绩录入
		api.POST("/save-rule", adminCtrl.SaveRule)         // 规则配置

		// --- 队伍相关接口 ---
		api.POST("/save-team", adminCtrl.SaveTeam) // 保存队伍
	}

	return r
}
