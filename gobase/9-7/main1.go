package main

import (
	"log"
	"math"
	"net"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	ip, ipNet, _ := net.ParseCIDR("192.168.1.1/29")
	ipFill := strings.Repeat("0", int(math.Pow(2, float64(32-29))))
	log.Println("ipFill:", ipFill)
	ipFill = "10000000"
	log.Println("[]byte(ip):", []byte(ip))
	log.Println("[]byte(ipNet.IP):", []byte(ipNet.IP))
	log.Println("ip[0], ip[1], ip[2], ip[3]", []byte(ipNet.IP)[0], []byte(ipNet.IP)[1], []byte(ipNet.IP)[2], []byte(ipNet.IP)[3])

	for c := range ipFill {
		if ipFill[c] == '0' {
			ipalloc := []byte(ipFill)
			log.Println(string(ipalloc))
			ipalloc[c] = '1'
			log.Println(string(ipalloc))
			ipFill = string(ipalloc)
			ip = ipNet.IP
			for t := uint(4); t > 0; t -= 1 {
				log.Println("t:", t, "[]byte(ip)[4-t]:", []byte(ip)[4-t], "uint8(c >> ((t - 1) * 8))", uint8(c>>((t-1)*8)))
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			ip[3] += 1
			break
		}
	}
	log.Println("ip", ip)
}
