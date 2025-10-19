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
			containers.GET("/:id", handlers.GetContainerInfo)
			// containers.POST("create", handlers.CreateContainer)
			// containers.POST("/:id/start", handlers.StartContainer)
			// containers.POST("/:id/stop", handlers.StopContainer)
			// containers.DELETE("/:id", handlers.DeleteContainer)
			// containers.GET("/:id/logs", handlers.GetContainerLogs)
		}
		// 镜像相关路由
		images := v1.Group("/images")
		{
			images.GET("", handlers.ListImages)
			// images.GET("/:id", handlers.GetImage)
			// images.DELETE("/:id", handlers.DeleteImage)
		}

		// 网络相关路由
		networks := v1.Group("/networks")
		{
			networks.GET("list", handlers.ListNetworks)
			// networks.POST("create", handlers.CreateNetwork)
			// networks.GET("/:id", handlers.GetNetwork)
			// networks.DELETE("/:id", handlers.DeleteNetwork)
		}
	}
}
