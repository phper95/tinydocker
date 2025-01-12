package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 文件路径
	filePath := "./example.txt"

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个缓冲区来存储 multipart/form-data 请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件到 multipart/form-data 请求体
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}

	// 关闭 writer 以完成 multipart/form-data 请求体
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/upload", body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// 打印响应
	fmt.Println("Response:", string(respBody))
}
