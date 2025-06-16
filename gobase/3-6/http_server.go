package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/test", func(resp http.ResponseWriter, request *http.Request) {
		//接收GET请求的参数
		fmt.Println(request.URL.Query())
		if request.Method == "POST" {
			//读取POST请求的body
			body, err := io.ReadAll(request.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(body))
		}
		//设置响应头
		resp.Header().Set("Content-Type", "application/json")

		//设置响应体
		resp.Write([]byte(`{"code": 0, "msg": "success"}`))
	})

	// 启动服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
