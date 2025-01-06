package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var clients = make(map[net.Conn]bool)
var broadcast = make(chan string)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		clients[conn] = true
		go handleClient(conn)
	}
}

func broadcaster() {
	for {
		message := <-broadcast
		for conn := range clients {
			_, err := fmt.Fprintln(conn, message)
			if err != nil {
				delete(clients, conn)
				conn.Close()
			}
		}
	}
}

func handleClient(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		broadcast <- scanner.Text()
	}
	delete(clients, conn)
	conn.Close()
}
