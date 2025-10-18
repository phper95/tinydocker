package routes

import (
	"github.com/phper95/tinydocker/internal/api/handlers"
	"github.com/phper95/tinydocker/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.Logger())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "TinyDocker API is running",
		})
	})

	// API 版本分组
	v1 := r.Group("/api/v1")
	// 放在代码块中，使代码结构更清晰
	{
		// 容器相关路由
		containers := v1.Group("/containers")
		{
			containers.GET("list", handlers.ListContainers)
		}
	}
}
