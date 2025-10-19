package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 错误的 MarshalJSON 实现，会导致无限递归
func (p *Person) MarshalJSON() ([]byte, error) {
	// 直接使用原类型会导致递归调用自身
	return json.Marshal(struct {
		*Person
		DisplayName string `json:"display_name"`
	}{
		Person:      p, // 这里直接引用了 Person 类型
		DisplayName: "Mr./Ms. " + p.Name,
	})
}

func main() {
	person := &Person{
		Name: "Alice",
		Age:  30,
	}

	// 这行代码会触发无限递归，最终导致栈溢出
	_, err := json.Marshal(person)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
