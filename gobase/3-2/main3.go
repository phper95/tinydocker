package main

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	slice := []int{0, 1, 2, 3}
	//var fns []func()
	//for _, val := range slice {
	//	fns = append(fns, func() {
	//		log.Println(val)
	//	})
	//}
	//for _, fn := range fns {
	//	fn()
	//}

	for _, val := range slice {
		func() {
			log.Println(val)
		}()
	}
}
