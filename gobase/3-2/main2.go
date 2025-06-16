package main

import "fmt"

const Name = "John"

var name = "John"

func main() {
	fmt.Println(name)
	name = "Jane"
	const Name = "xxx"
	fmt.Println(Name)
	readConst()
	name = "Bob"
	fmt.Println(name)

	//基本类型
	var a int = 10
	var in8 int8 = 10
	var in16 int16 = 10
	var in32 int32 = 10
	var in64 int64 = 10
	var u8 uint8 = 10
	var u16 uint16 = 10
	var u32 uint32 = 10
	var u64 uint64 = 10
	var f32 float32 = 10.5
	var f64 float64 = 10.5
	var b bool = true
	var s string = "Hello"
	fmt.Println(a, in8, in16, in32, in64, u8, u16, u32, u64, f32, f64, b, s)
	//复合类型
	var arr [5]int = [5]int{1, 2, 3, 4, 5}
	arr[2] = 10
	var slice []int = []int{1, 2, 3, 4, 5}
	slice = append(slice, 6)
	var map1 map[string]int = map[string]int{"a": 1, "b": 2, "c": 3}
	map1["d"] = 4
	var ch chan int = make(chan int)
	//ch <- 10
	fmt.Println(arr, slice, map1, ch)

	// 匿名函数赋值给变量
	var add = func(a, b int) int {
		return a + b
	}
	fmt.Println(add(1, 2))
	// 调用匿名函数
	fmt.Println(func(a, b int) int {
		return a + b
	}(1, 2))

	//闭包函数
	cnt := func() func() int {
		counter := 0
		return func() int {
			counter++
			return counter
		}
	}()
	fmt.Println(cnt())
	fmt.Println(cnt())
}

func readConst() {
	fmt.Println(Name)
}

// 求和函数
func sum(a, b int) int {
	return a + b
}

// 定义一个多返回值函数
func swap(a, b int) (int, int) {
	return b, a
}
