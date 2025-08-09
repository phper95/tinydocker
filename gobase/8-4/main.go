package main

/*
#include <stdio.h>
#include <stdlib.h>

// C函数：计算两个整数的和
int add(int a, int b) {
    return a + b;
}

// C函数：打印字符串
void print_message(char* msg) {
    printf("C says: %s\n", msg);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	// 调用C函数add
	result := C.add(10, 20)
	fmt.Printf("10 + 20 = %d\n", int(result))

	// 调用C函数print_message
	msg := C.CString("Hello from Go!")
	defer C.free(unsafe.Pointer(msg)) // 释放C字符串内存
	C.print_message(msg)
}
