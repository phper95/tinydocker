package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	//连接到TCP服务器
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()
	fmt.Println("连接到服务器成功，请输入消息：")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		// 将用户输入的消息发送到服务器
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
			return
		}
		// 接收服务器的响应
		reply, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("服务器响应:", reply)
	}
}
