package main

import (
	"github.com/gin-gonic/gin"
	local_pkg "github.com/phper95/localpkg"
	//local_pkg "gitlab.com/phper95/localpkg"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	local_pkg.LocalFunc()
}

//Jaeger和Zipkin。
