package main

func main() {
	ch := make(chan int, 2)
	ch1 := make(chan int, 2)
	go func() {
		for {
			select {
			case <-ch:
				ch1 <- 1
			}
		}
	}()

	for {
		select {
		case <-ch1:
			ch <- 1
		}
	}

}
