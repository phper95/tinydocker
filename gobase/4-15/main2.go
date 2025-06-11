package main

import (
	"fmt"
	"regexp"
)

// 定义一个函数 validate_email，检查字符串是否为合法邮箱格式
func validate_email(email string) bool {
	// 定义邮箱的正则表达式，这里用的是最基础的格式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func main() {
	fmt.Println(validate_email("example@gmail.com")) // 输出：true
	fmt.Println(validate_email("example@.com"))      // 输出：false
}
