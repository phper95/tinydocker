package routes

import (
	"crypto/rand"
	"github.com/phper95/tinydocker/pkg/logger"
	"log"
	"os"

	"github.com/phper95/tinydocker/internal/api/handlers"
	"github.com/phper95/tinydocker/internal/api/middleware"
	"github.com/phper95/tinydocker/internal/service"

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

	// 加载 JWT 密钥：优先从环境变量 JWT_SECRET 读取，否则随机生成
	var jwtSecret []byte
	if env := os.Getenv("JWT_SECRET"); env != "" {
		jwtSecret = []byte(env)
	} else {
		jwtSecret = make([]byte, 32)
		if _, err := rand.Read(jwtSecret); err != nil {
			logger.Error("生成 JWT 密钥失败: %v", err)
			log.Fatal(err)
		}
	}
	authService := service.NewAuthService(jwtSecret)
	// 创建处理器
	authHandler := handlers.NewAuthHandler(authService)
	// 认证相关路由（不需要认证）
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// 受保护的 API 版本分组
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware(authService))
	{
		// 认证相关受保护端点
		authProtected := v1.Group("/auth")
		{
			authProtected.GET("/profile", authHandler.GetProfile)
			authProtected.POST("/logout", authHandler.Logout)
		}
		// 容器相关路由
		containers := v1.Group("/containers")
		{
			containers.GET("list", middleware.RequirePermission("containers", "list"), handlers.ListContainers)
			containers.GET("/:id", middleware.RequirePermission("containers", "get"), handlers.GetContainerInfo)
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
