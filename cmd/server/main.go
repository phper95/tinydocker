package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/phper95/tinydocker/internal/api/routes"
	"github.com/phper95/tinydocker/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	// 在 goroutine 中启动服务器
	go func() {
		log.Printf("Starting TinyDocker API server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// syscall.SIGINT,中断信号（Interrupt Signal）如Ctrl+C 触发， kill -INT <pid> 或 kill -2 <pid>
	// syscall.SIGTERM,终止信号（Termination Signal）如kill 触发: kill -TERM <pid> 或 kill -15 <pid>
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	// 设置超时 context 用于关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
