package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/phper95/tinydocker/docs" // 注意：替换为实际生成的 docs 目录路径
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
)

func main() {
	r := gin.Default()
	// 注册 Swagger 文档访问路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := r.Run(":8080") // 启动服务
	if err != nil {
		log.Fatal(err)
	}
}
