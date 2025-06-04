package main

import (
	"fmt"
	"net"
)

func main() {
	// 定义UDP服务器的地址和端口
	addr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 8081,
	}
	// 启动UDP服务器，监听指定的地址和端口
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("UDP server is running...")
	buffer := make([]byte, 1024)
	for {
		// 从UDP连接中读取数据
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Received from", clientAddr.String(), ":", string(buffer[:n]))
		//响应客户端的请求
		_, err = conn.WriteToUDP([]byte("Message received."+string(buffer[:n])), clientAddr)
		if err != nil {
			fmt.Println(err)
		}
	}
}
