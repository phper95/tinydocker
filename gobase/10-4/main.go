package main

import "fmt"

func main() {
	a := 1
	{
		// 可以访问和修改代码块外部的变量 a
		a = 2
		// 代码块中声明的变量 b 只能在代码块中使用，外部无法访问
		// b := 2
		fmt.Println(a)
	}
	fmt.Println(a)
	// fmt.Println(b)
}
