package main

import (
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/network"
	"github.com/phper95/tinydocker/pkg/db"
	"log"
	"net"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := db.InitBoltDBClient(db.DefaultBoltDBClientName, enum.DefaultNetworkDBPath)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer func() {
		err := db.GetBoltDBClient(db.DefaultBoltDBClientName).Close()
		if err != nil {
			log.Println(err)
		}
	}()
	nw, err := network.GetNetworkFromDB("mybridge")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("network: %+v \n", nw)
	subnet := nw.IPRange
	log.Printf("subnet: %s , subnet.IP: %s \n", subnet.String(), subnet.IP.String())

	Ip, netIP, _ := net.ParseCIDR(subnet.String())
	log.Printf("Ip %s, netip %s \n", Ip.String(), netIP.String())
	network.LoadIP()
	ip, err := network.AllocateIP(netIP)
	if err != nil {
		log.Println(err)
	}
	log.Printf("allocated ip: %s \n", ip.String())

}
