package main

import (
	"fmt"
	"sync"
)

type OrderService struct {
	mu sync.Mutex
}

func (s *OrderService) ProcessOrder(orderID string, isUrgent bool, hasDiscount bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if isUrgent {
		s.mu.Lock()
		fmt.Println("Handling urgent order:", orderID)
		s.mu.Unlock()
	}

	if hasDiscount {
		s.mu.Lock()
		fmt.Println("Applying discount for order:", orderID)
		s.mu.Unlock()
	}

	fmt.Println("Order processed successfully:", orderID)
	return nil
}

func main() {
	service := &OrderService{}
	err := service.ProcessOrder("12345", true, true)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
