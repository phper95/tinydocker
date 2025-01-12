package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to TCP server. Type your message:")
	reader := bufio.NewReader(os.Stdin)

	for {
		// 从用户输入读取消息
		fmt.Print(">> ")
		message, _ := reader.ReadString('\n')

		// 发送消息到服务器
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// 接收服务器的响应
		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}
		fmt.Printf("Server reply: %s", reply)
	}
}
