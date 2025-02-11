package main

import "fmt"

func main() {
	// 经典 for 循环
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue
		}
		println(i)
	}

	//条件 for 循环
	count := 0
	for count < 10 {
		fmt.Println(count)
		count++
	}

	//无限循环
	sum := 0
	for {
		fmt.Println(sum)
		sum++
		if sum >= 10 {
			break
		}
		//return
	}

	//for range 循环
	nums := []int{1, 2, 3, 4, 5}
	for i, num := range nums {
		fmt.Println(i, num)
	}
	for _, num := range nums {
		fmt.Println(num)
	}
	for i := range nums {
		fmt.Println(i)
	}

	str := "Go语言"
	for i, char := range str {
		fmt.Printf("%d %c\n", i, char)
	}
	map1 := map[string]int{"apple": 1, "banana": 2, "orange": 3}
	for key, value := range map1 {
		fmt.Printf("%s %d\n", key, value)
	}

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	for num := range ch {
		fmt.Println(num)
	}
}
