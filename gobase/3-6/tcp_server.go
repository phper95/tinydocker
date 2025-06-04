package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//启动一个TCP监听器，监听端口8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server is running...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("New client connected:", conn.RemoteAddr())
		//处理客户端的连接请求
		go handleTCPClient(conn)
	}
}

func handleTCPClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Received from client:", line)
		//向客户端发送响应数据
		_, err = conn.Write([]byte("Message received." + line))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
