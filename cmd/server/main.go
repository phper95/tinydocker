package main

import (
	"fmt"
	"log"

	"github.com/phper95/tinydocker/internal/api/routes"
	"github.com/phper95/tinydocker/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting TinyDocker API server on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// 优雅关闭

}
