package main

import "fmt"

func main() {
	num := 10
	if num > 5 {
		//x := 10
		print("num is greater than 5")
	} else {
		//fmt.Println(x)
		print("num is less than or equal to 5")
	}

	//fmt.Println(x)

	if num%2 == 0 {
		print("num is even")
	} else {
		print("num is odd")
	}

	if val := num / 2; val > 5 {
		fmt.Println("num divided by 2 is greater than 5")
	}

}
