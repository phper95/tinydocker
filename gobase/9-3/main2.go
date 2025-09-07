package main

import (
	"encoding/json"
	"github.com/vishvananda/netlink"
	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
func main() {
	// 创建veth pair
	la := netlink.LinkAttrs{Name: "veth0"}
	veth := &netlink.Veth{LinkAttrs: la, PeerName: "veth1"}
	if err := netlink.LinkAdd(veth); err != nil {
		log.Fatal(err)
	}
	// 列出所有的网络接口
	links, err := netlink.LinkList()
	if err != nil {
		log.Fatal(err)
	}
	s, _ := json.Marshal(links)
	log.Println(string(s))

}
