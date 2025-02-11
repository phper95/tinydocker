package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/api/v1/test", handleTestfunc)
	http.HandleFunc("/api/v1/upload", uploadTestfunc)

	// 启动服务
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// 定义处理函数
func handleTestfunc(resp http.ResponseWriter, request *http.Request) {
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
}

// 文件上传处理函数
func uploadTestfunc(resp http.ResponseWriter, request *http.Request) {
	// 设置最大上传文件大小为10MB
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取上传的文件
	file, header, err := request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// 创建文件保存路径
	dst, err := os.Create(filepath.Join("upload", header.Filename))
	if err != nil {
		fmt.Println(err)
		http.Error(resp, "上传失败", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 复制文件内容到目标文件
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println(err)
		http.Error(resp, "上传失败", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	resp.Header().Set("Content-Type", "application/json")
	// 设置响应体
	resp.Write([]byte(`{"code": 0, "msg": "success"}`))

}
