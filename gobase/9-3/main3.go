package main

import (
	"github.com/vishvananda/netns"
	"log"
)

func main() {
	// 创建network namespace
	// 获取当前的网络命名空间
	origns, err := netns.Get()
	if err != nil {
		log.Fatal("Failed to get original namespace:", err)
	}
	log.Println("original ns:", origns.String())
	defer origns.Close() // 延迟关闭原始网络命名空间

	// 创建一个新的网络命名空间
	newns, err := netns.NewNamed("ns1")
	if err != nil {
		log.Fatal("Failed to create new namespace:", err)
	}
	defer newns.Close() // 延迟关闭新创建的网络命名空间

	// 切换到新的网络命名空间
	err = netns.Set(newns)
	if err != nil {
		log.Println("Failed to switch to new namespace:", err)
	}

	// 获取当前的网络命名空间
	ns, err := netns.Get()
	if err != nil {
		log.Println("Failed to get current namespace:", err)
	}
	log.Println("current ns:", ns.String())

	// 切换回原始的网络命名空间
	if err := netns.Set(origns); err != nil {
		log.Fatal("Failed to switch back to original namespace:", err)
	}
}
