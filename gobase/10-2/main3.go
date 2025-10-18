package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	err := r.Run(":80")
	if err != nil {
		panic(err)
	}
}
