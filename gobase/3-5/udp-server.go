package main

import (
	"fmt"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("0.0.0.0"),
	}

	// 启动 UDP 服务器
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()
	fmt.Println("UDP server is running on port 8081")

	buffer := make([]byte, 1024)
	for {
		// 接收数据
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP client:", err)
			continue
		}
		fmt.Printf("Message received from %s: %s \n", clientAddr.String(), string(buffer[:n]))

		// 发送响应
		_, err = conn.WriteToUDP([]byte("Message received"), clientAddr)
		if err != nil {
			fmt.Println("Error sending response:", err)
		}
	}
}
