package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// 连接到服务器
	conn, err := net.DialUDP("udp", nil, &serverAddr)
	if err != nil {
		fmt.Println("Error connecting to UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to UDP server. Type your message:")
	buffer := make([]byte, 1024)

	for {
		// 从用户输入读取消息
		var message string
		fmt.Print(">> ")
		fmt.Fscanf(os.Stdin, "%s\n", &message)

		// 发送消息到服务器
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// 接收服务器的响应
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Printf("Server reply: %s\n", string(buffer[:n]))
	}
}
