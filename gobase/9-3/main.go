package main

import (
	"log"
	"net"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ip, ipNet, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ip)
	log.Println(ipNet)
	ip = net.ParseIP("192.168.1.1")
	log.Println(ip)
}
