package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func MiddlewareName() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 请求处理前的逻辑
		// ...

		// 调用下一个中间件或路由处理函数（必须调用，否则请求会被阻塞）
		c.Next()

		// 2. 请求处理后的逻辑（在路由处理函数执行后执行）
		// ...
	}
}
func main() {
	r := gin.Default()
	// 使用中间件（全局使用，所有路由都会经过该中间件）
	r.Use(RequestTimer())
	r.GET("/user/list", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "user list"})
	})

	// 路由组使用中间件（/admin下的所有路由都需要身份验证）
	adminGroup := r.Group("/admin")
	adminGroup.Use(AuthMiddleware())
	adminGroup.GET("/user/list", func(c *gin.Context) {
		c.JSON(200, gin.H{"data": "user list"})
	})
	adminGroup.POST("/user/add", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "user added"})
	})
	// 测试：访问http://localhost/admin/user/list，若请求头没有Authorization或值不是valid-token，返回401；否则返回成功

	err := r.Run(":80")
	if err != nil {
		log.Fatalf("Server start failed: %v", err)
	}

}

// 请求耗时统计中间件
func RequestTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 前处理：记录请求开始时间
		start := time.Now()

		// 调用下一个处理环节
		c.Next()

		// 后处理：计算耗时
		duration := time.Since(start).Milliseconds()
		// 获取请求路径
		path := c.FullPath()
		// 打印耗时日志
		log.Printf("Path: %s, Duration: %v", path, duration)
	}
}

// 身份验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		token := c.GetHeader("Authorization")
		// 简单验证token（实际项目中需结合JWT等方式）
		if token != "valid-token" {
			// 中断请求，返回401错误
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		// token验证通过，继续执行
		c.Next()
	}
}
