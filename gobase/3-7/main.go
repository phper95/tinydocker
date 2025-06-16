package main

import "fmt"

func main() {
	f1()
	sum(1, 2, 3, 4, 5)
	sum2([]int{1, 2, 3, 4, 5})
	s := []interface{}{1, 2, 3, 4, 5}
	f(s...)
	calculate(1, 2)
}

// 定义一个函数，返回两个值：和与积
func calculate(a, b interface{}) (sum int, product int) {
	sum = a.(int) + b.(int)
	product = a.(int) * b.(int)
	return
}

func f1() (m map[int]int) {
	fmt.Println(m == nil)
	m = make(map[int]int)
	m[1] = 1
	return
}

func sum(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}
func sum2(nums []int) int {
	total := 0
	for _, num := range nums {
		total += num
	}
	return total
}

func f(params ...interface{}) {
	for i, param := range params {
		fmt.Println("i", i, "param", param)
	}
}
