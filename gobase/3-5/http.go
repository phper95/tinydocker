package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	http.HandleFunc("/api/v1/test", testHandler)
	http.HandleFunc("/api/v1/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	//接收GET请求的参数
	fmt.Println(r.URL.Query())
	//判断是否POST请求
	if r.Method == "POST" {
		//读取POST请求的body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		//打印POST请求的body
		fmt.Println(string(body))
	}
	//设置响应头为json格式
	w.Header().Set("Content-Type", "application/json")

	//设置响应内容
	w.Write([]byte(`{"code": 0, "msg": "success"}`))
}

// 文件上传处理函数
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 设置最大上传文件大小为10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Println("Error parsing the multipart form:", err)
	}

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 创建文件保存路径
	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		log.Println("Error creating the file for writing:", err)
		http.Error(w, "Error creating the file for writing", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 复制文件内容到目标文件
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Println("Error copying the file:", err)
		http.Error(w, "Error copying the file", http.StatusInternalServerError)
		return
	}

	// 设置响应头为json格式
	w.Header().Set("Content-Type", "application/json")

	// 设置响应内容
	w.Write([]byte(`{"code": 0, "msg": "file uploaded successfully"}`))
}
