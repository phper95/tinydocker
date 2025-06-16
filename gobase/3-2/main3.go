package main

import (
	"fmt"
	"net/http"
	"sort"
)

func main() {
	nums := []int{1, 2, 3, 4, 5}
	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	go func() {
		// do something
	}()

	defer func() {
		// do something
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	dobule := multiply(2)
	fmt.Println(dobule(5))

	fmt.Println(multiply(2)(5))
	fmt.Println(multiply(3)(5))

}

func multiply(factor int) func(int) int {
	return func(num int) int {
		return num * factor
	}
}
