package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CORS() gin.HandlerFunc {
	// 定义CORS（跨域资源共享）配置
	config := cors.Config{
		// 允许所有来源访问，生产环境中可以根据需要进行限制
		// 例如：[]string{"http://example.com", "http://another.com"}
		// 这里使用"*"表示允许所有来源访问
		AllowOrigins: []string{"*"},
		// 允许的HTTP方法，可以根据需要进行调整
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 允许的HTTP头，可以根据需要进行调整
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		// 允许浏览器访问的响应头，可以根据需要进行调整
		ExposeHeaders: []string{"Content-Length"},
		// 是否允许携带凭证（如Cookies），根据需要进行调整
		// 这里设置为true，表示允许携带凭证
		AllowCredentials: true,
		// 预检请求的缓存时间
		// 预检请求（Preflight Request）
		// 预检请求是CORS机制中的一种安全措施，具有以下特点：
		// 触发条件：当浏览器发起跨域请求时，如果请求满足以下条件之一，会先发送预检请求：
		// 使用PUT、DELETE等非简单方法
		// 包含自定义请求头（如 Authorization）
		// Content-Type为 application/json 等非简单类型
		// 请求方式：浏览器会先自动发送一个 OPTIONS 方法的请求到目标服务器
		// 目的：询问服务器是否允许即将发送的实际请求，包含：
		// 请求方法是否被允许（Access-Control-Request-Method）
		// 请求头是否被允许（Access-Control-Request-Headers）
		// 缓存机制：服务器可以通过 MaxAge 响应头告诉浏览器预检结果可以缓存多久，避免重复发送预检请求
		// 这里设置为12小时，表示预检请求的结果可以缓存12小时
		MaxAge: 12 * time.Hour,
	}

	return cors.New(config)
}
