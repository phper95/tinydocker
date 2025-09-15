package main

import "net"

func main() {
	ip, ipNet, _ := net.ParseCIDR("192.168.1.2/29")
	println(ip.String(), ipNet.String())
}
