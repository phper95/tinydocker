package main

import "github.com/gin-gonic/gin"

func main() {
	// 创建Gin引擎（默认模式，生产环境建议使用gin.ReleaseMode）
	r := gin.Default()

	// 定义GET路由
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello Gin!")
	})

	// 定义带参数的路由（参数用:标识）
	r.GET("/user/:id", func(c *gin.Context) {
		// 获取路由参数
		id := c.Param("id")
		c.String(200, "User ID: %s", id)
	})

	// 测试：访问http://localhost/user/123，返回"User ID: 123"

	r.GET("/search", func(c *gin.Context) {
		// 获取查询参数，第二个参数为默认值（当参数不存在时使用）
		keyword := c.Query("keyword")
		page := c.DefaultQuery("page", "1")
		size := c.DefaultQuery("size", "10")

		c.String(200, "Keyword: %s, Page: %s, Size: %s", keyword, page, size)
	})

	// 测试：访问http://localhost/search?keyword=gin&page=2，返回"Keyword: gin, Page: 2, Size: 10"

	// 1. 定义结构体（用于绑定请求体数据）
	type UserRequest struct {
		Name  string `json:"name" binding:"required"` // required表示该字段为必填
		Age   int    `json:"age" binding:"min=18"`    // min=18表示年龄最小为18
		Email string `json:"email" binding:"email"`   // email表示需符合邮箱格式
	}

	// 2. 定义POST路由，解析JSON请求体
	r.POST("/user", func(c *gin.Context) {
		var req UserRequest
		// 解析JSON请求体到结构体，若解析失败（如字段缺失、格式错误），会返回400错误
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		// 解析成功，返回用户信息
		c.JSON(200, gin.H{
			"message": "create user success",
			"data": gin.H{
				"name":  req.Name,
				"age":   req.Age,
				"email": req.Email,
			},
		})
	})

	// 测试：用Postman发送POST请求到http://localhost:8080/user，请求体为
	// {
	//   "name": "Zhang San",
	//   "age": 20,
	//   "email": "zhangsan@example.com"
	// }
	// 会返回成功响应；若age设为17，会返回400错误

	// 结构体标签使用form
	type LoginRequest struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required,min=6"`
	}

	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		// 解析表单请求体
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "login success", "username": req.Username})
	})
	// 测试：用Postman以表单形式发送POST请求，填写username和password参数

	// 启动服务，监听80端口
	r.Run(":80")
}
