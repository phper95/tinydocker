package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5}
	m := make(map[int]*int)
	for _, num := range nums {
		fmt.Printf("num%d； addr：%p\n ", num, &num)
		m[num] = &num
		num = num * 2
	}
	//fmt.Println(nums)
	fmt.Println(m)
	for k, v := range m {
		fmt.Printf("key:%d, value:%d\n", k, *v)
	}
	//for i, num := range nums {
	//	nums[i] = num * 2
	//}
	//fmt.Println(nums)
	var ap *[3]int
	fmt.Println(ap)
	//for i, p := range ap {
	//	fmt.Println(i, p)
	//}
	//for i, _ := range ap {
	//	fmt.Println(i)
	//}
	//var fns []func()
	//for i, num := range nums {
	//	fns = append(fns, func() {
	//		fmt.Println(i, num)
	//	})
	//}
	//
	//for _, fn := range fns {
	//	fn()
	//}

	//for i, num := range nums {
	//	func() {
	//		fmt.Println(i, num)
	//	}()
	//}
	//遍历 map 时的无序性
	mp := map[int]int{1: 1, 2: 2, 3: 3, 4: 4, 5: 5}
	for k, v := range mp {
		fmt.Println(k, v)
	}
}
