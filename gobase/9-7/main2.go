package main

import (
	"log"
	"net"
)

func main() {
	ip, subnet, _ := net.ParseCIDR("192.168.1.5/24")
	// 计算要释放的IP在位图中的索引
	// 需要减去1，因为在分配时有加1的操作
	ipIndex := 0
	ipBytes := []byte(ip.To4())
	subnetBytes := []byte(subnet.IP.To4())
	log.Println("ipBytes", ipBytes, "subnetBytes", subnetBytes)
	// 循环处理IP地址的4个字节
	// ipBytes[i] - 获取目标IP地址的第i个字节[192 168 1 5]
	// subnetBytes[i] - 获取子网起始IP的第i个字节[192 168 1 0]
	// ipBytes[i] - subnetBytes[i] - 计算目标IP与子网起始IP在该字节上的差值
	// ((3 - uint(i)) * 8)) - 计算该字节在整个IP地址中的位移量,因为IP地址是4个字节，所以位移量是24、16、8、0
	// int(diff << ((3 - uint(i)) * 8)) - 将字节差值左移相应的位移量，
	// 左移操作 diff << ((3 - uint(i)) * 8) 是为了将IP地址各个字节的差值放到正确的位置上，从而计算出目标IP地址相对于子网起始地址的偏移量
	// 举个具体例子：
	// 假设子网是 192.168.1.0/24，我们要计算IP地址 192.168.1.5 的偏移量：
	// 目标IP: 192.168.1.5
	// 子网起始IP: 192.168.1.0
	// 计算过程：
	// i=0 : diff = 192-192=0, ipIndex = 0 + (0 << 24) = 0
	// i=1 : diff = 168-168=0, ipIndex = 0 + (0 << 16) = 0
	// i=2 : diff = 1-1=0, ipIndex = 0 + (0 << 8) = 0
	// i=3 : diff = 5-0=5, ipIndex = 0 + (5 << 0) = 5
	for i := range ipBytes {
		diff := uint(ipBytes[i] - subnetBytes[i])
		log.Println("i", i, " ipBytes[i]", ipBytes[i], "subnetBytes[i]", subnetBytes[i], " diff", diff)

		ipIndex += int(diff << ((3 - uint(i)) * 8))
		log.Println("ipIndex", ipIndex, " ((3 - uint(i)) * 8))", ((3 - uint(i)) * 8))
	}

	// 减去网络地址偏移
	ipIndex -= 1
	log.Println("ipIndex", ipIndex)
}
