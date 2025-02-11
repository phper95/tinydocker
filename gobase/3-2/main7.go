package main

import "fmt"

func main() {
	//根据分数输出等级
	score := 85
	switch {
	case score >= 90:
		println("A")
	case score >= 70:
		println("C")
		//break
		fallthrough
	case score >= 80:
		println("B")

	case score >= 60:
		println("D")
	default:
		println("E")
	}

	day := "Monday123"
	switch day {
	case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
		fmt.Println("Working day")
	case "Saturday":
		fmt.Println("Weekend")
	default:
		fmt.Println("other day")
	}
}
