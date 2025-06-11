package main

import "github.com/gin-gonic/gin"

// 使用gin框架创建一个 GET 接口 /hello，返回 JSON 格式的 Hello World
func main() {
	// 初始化 gin 引擎
	r := gin.Default()

	// 设置 GET /hello 的路由，返回 JSON 数据
	r.GET("/hello", func(c *gin.Context) {
		// 从查询参数中获取 name，默认为 "World"
		name := c.DefaultQuery("name", "World")
		c.JSON(200, gin.H{
			"message": "Hello " + name + "!",
		})
	})

	// 启动服务器并监听端口8080
	r.Run(":8080")
}
