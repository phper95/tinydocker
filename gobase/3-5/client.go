package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// 连接服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// 读取欢迎信息和昵称提示
	reader := bufio.NewReader(conn)
	welcomeMessage, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading welcome message:", err)
		os.Exit(1)
	}
	fmt.Print(welcomeMessage)

	// 发送昵称
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		nickname := scanner.Text()
		_, err := fmt.Fprintln(conn, nickname)
		if err != nil {
			log.Println("Error sending nickname:", err)
			return
		}
	}
	// 开启 Goroutine 接收消息
	go receiveMessages(conn)
	// 发送消息到服务器
	sendMessages(conn)
}

func receiveMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error receiving messages:", err)
	}
}

func sendMessages(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_, err := fmt.Fprintln(conn, scanner.Text())
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error scanning input:", err)
	}
}
