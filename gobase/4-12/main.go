package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

func main() {

}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	body, _ := io.ReadAll(r.Body)
	var user struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.Unmarshal(body, &user)

	// 连接数据库
	db, _ := sql.Open("mysql", "user:pass@tcp(127.0.0.1:3306)/dbname")
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", user.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 验证密码
	if user.Password != storedPassword {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 生成 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("secret-key"))

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
