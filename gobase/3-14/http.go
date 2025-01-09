package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request")
		// 模拟耗时操作
		time.Sleep(5 * time.Second)

		// 设置响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := `{"message": "Response after 5 seconds"}`
		fmt.Fprintln(w, response)
	})

	port := ":80"
	log.Printf("Server is running at http://localhost%s\n", port)

	// 启动 HTTP 服务
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Failed to start server: %s\n", err)
	}
}
