package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Response 定义从服务器接收的 JSON 结构
type Response struct {
	ClientIP   string `json:"client_ip"`
	AccessTime string `json:"access_time"`
}

func main() {
	// 发起 HTTP GET 请求
	resp, err := http.Get("http://192.168.1.5/api/info")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// 解析 JSON 响应
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// 输出结果
	fmt.Printf("Client IP: %s\n", response.ClientIP)
	fmt.Printf("Access Time: %s\n", response.AccessTime)

	// 解析时间并显示本地时间
	if t, err := time.Parse(time.RFC3339, response.AccessTime); err == nil {
		fmt.Printf("Local Time: %s\n", t.Local().Format("2006-01-02 15:04:05"))
	}
}
