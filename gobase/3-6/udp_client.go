package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 定义服务器的地址和端口
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8081,
	}
	//连接UDP服务器
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		log.Println("连接服务器失败：", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to server,type something to send:")
	buffer := make([]byte, 1024)
	for {
		// 读取用户输入
		fmt.Print("> ")
		var input string
		_, err := fmt.Fscanf(os.Stdin, "%s\n", &input)
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}
		// 将消息发送到服务器
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Println("Error writing to server:", err)
			continue
		}
		// 接收服务器的响应
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error reading from server:", err)
			continue
		}
		// 打印服务器的响应
		fmt.Printf("Received from %s: %s\n", clientAddr, string(buffer[:n]))

	}
}
