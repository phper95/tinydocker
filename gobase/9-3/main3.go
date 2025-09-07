package main

import (
	"github.com/vishvananda/netns"
	"log"
)

func main() {
	// 创建network namespace
	// 保存原始网络命名空间
	origns, err := netns.Get()
	if err != nil {
		log.Fatal("Failed to get original namespace:", err)
	}
	log.Println("original ns:", origns.String())
	defer origns.Close()
	newns, err := netns.NewNamed("ns1")
	if err != nil {
		log.Fatal("Failed to create new namespace:", err)
	}
	defer newns.Close()

	err = netns.Set(newns)
	if err != nil {
		log.Println("Failed to switch to new namespace:", err)
	}

	ns, err := netns.Get()
	if err != nil {
		log.Println("Failed to get current namespace:", err)
	}
	log.Println("current ns:", ns.String())
	// 切换回原始命名空间
	if err := netns.Set(origns); err != nil {
		log.Fatal("Failed to switch back to original namespace:", err)
	}
}
