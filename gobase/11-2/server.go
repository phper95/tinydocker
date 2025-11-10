package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Response 定义返回的 JSON 结构
type Response struct {
	ClientIP   string `json:"client_ip"`
	AccessTime string `json:"access_time"`
}

// handler 处理 HTTP 请求
func handler(w http.ResponseWriter, r *http.Request) {
	// 获取客户端 IP
	clientIP := getClientIP(r)

	// 获取当前时间
	accessTime := time.Now().Format(time.RFC3339)

	// 创建响应结构体
	response := Response{
		ClientIP:   clientIP,
		AccessTime: accessTime,
	}

	// 设置响应头为 JSON 格式
	w.Header().Set("Content-Type", "application/json")

	// 将响应结构体编码为 JSON 并返回
	json.NewEncoder(w).Encode(response)
}

// getClientIP 获取客户端真实 IP 地址
func getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For 头部（通常用于代理服务器）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// 如果有多个 IP，取第一个
		ipList := splitIPs(forwarded)
		if len(ipList) > 0 {
			return ipList[0]
		}
	}

	// 检查 X-Real-IP 头部
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 直接从连接获取远程地址
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// splitIPs 分割 IP 列表字符串
func splitIPs(ips string) []string {
	var result []string
	for _, ip := range splitString(ips, ",") {
		trimmed := trimSpace(ip)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitString 简单分割字符串实现
func splitString(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i:i+1] == sep {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	return result
}

// trimSpace 去除字符串前后空格
func trimSpace(s string) string {
	start := 0
	end := len(s)

	// 找到开始位置（跳过前导空格）
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}

	// 找到结束位置（跳过后导空格）
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}

	return s[start:end]
}

func main() {
	// 注册路由处理器
	http.HandleFunc("/api/info", handler)

	// 启动服务器
	fmt.Println("Server starting on :80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
