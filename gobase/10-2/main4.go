package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	// r := gin.Default()

	// 等价于
	// 创建不包含默认中间件的引擎
	r := gin.New()
	// 添加日志中间件
	r.Use(gin.Logger())
	// 添加恢复中间件（当处理 HTTP 请求的 goroutine 中发生 panic 时，Recovery 中间件可以捕获这个 panic，向客户端返回一个 HTTP 500 Internal Server Error 响应，而不是让连接挂起或断开。）
	// r.Use(gin.Recovery())

	r.GET("/index/:id", func(c *gin.Context) {
		if c.Param("id") == "1" {
			panic("An error occurred!")
		}
		c.String(200, "Index Page")
	})
	err := r.Run(":80")
	if err != nil {
		fmt.Println("server start failed", err)
	}
}
