package main

import "fmt"

type Product struct {
	ID    int
	Name  string
	Price float64
}

type OrderItem struct {
	Product  *Product
	Quantity int
}

type Order struct {
	ID           int
	Items        []*OrderItem
	CustomerName string
}

func calculateTotal(order *Order) float64 {
	var total float64 = 0
	for _, item := range order.Items {
		total += item.Product.Price * float64(item.Quantity) // 可能出现空指针异常的地方
	}
	return total
}

func processOrder(order *Order) {
	total := calculateTotal(order)
	fmt.Printf("Processing order for %s, Total: %.2f\n", order.CustomerName, total)
}

func main() {
	product1 := &Product{ID: 1, Name: "Laptop", Price: 999.99}
	orderItem1 := &OrderItem{Product: product1, Quantity: 1}
	orderItem2 := &OrderItem{Product: nil, Quantity: 2} // 这里故意设置为nil模拟错误情况

	order := &Order{
		ID:           1,
		Items:        []*OrderItem{orderItem1, orderItem2},
		CustomerName: "John Doe",
	}

	processOrder(order)
}
