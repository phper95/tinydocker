package main

/*
// 这里编写C代码或C预处理指令
#include <stdio.h>

// 定义一个C函数
int add(int a, int b) {
    return a + b;
}
*/
import "C" // 必须紧跟在C代码注释后，不能有空行
import "fmt"

func main() {
	// 调用C函数：注意类型转换（Go类型 -> C类型）
	a := C.int(2)
	b := C.int(3)
	res := C.add(a, b)

	fmt.Printf("2 + 3 = %d\n", res) // 输出：2 + 3 = 5
}
