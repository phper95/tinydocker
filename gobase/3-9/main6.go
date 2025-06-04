package main

import "log"

func main() {
	var i int
	var str string
	var b bool
	var p *int
	var s []int
	var m map[string]int
	var c chan int
	var f func()
	log.Println("i", i, "str", str, "b", b, "p", p, "s", s == nil, "m", m == nil, "c", c, "f", f)

}
