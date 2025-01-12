package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var clients = make(map[net.Conn]bool)
var broadcast = make(chan string)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	// 启动监听器
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("TCP server is running on port 8080")

	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Client connected")
		go handleTCPConnection(conn)
	}
}

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// 读取客户端发送的数据
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}
		fmt.Printf("Message received: %s", message)

		// 向客户端发送响应
		_, err = conn.Write([]byte("Message received: " + message))
		if err != nil {
			fmt.Println("Error writing to client:", err)
			return
		}
	}
}
