package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建带默认中间件（日志与恢复）的 Gin 路由器实例
	r := gin.Default()

	// 定义简单的 GET 请求路由
	r.GET("/ping", func(c *gin.Context) {
		// 返回 JSON 响应
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 默认端口 8080 启动服务器
	// 监听 0.0.0.0:8080（Windows 下为 localhost:8080）
	r.Run()
}
