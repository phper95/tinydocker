package network

import (
	"fmt"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/db"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/vishvananda/netlink"
	"log"
	"net"
	"os/exec"
	"strings"
)

// 存储每个子网的 IP 分配状态
var allocatedIP map[string]string

func init() {
	InitBoltDB()
	err := LoadIP()
	if err != nil {
		logger.Error("Failed to load allocated IP", allocatedIP, "err:", err)
	} else {
		logger.Info("Loaded allocated IP: ", allocatedIP)
	}
}
func InitBoltDB() {
	err := db.InitBoltDBClient(db.DefaultBoltDBClientName, enum.DefaultNetworkDBPath)
	if err != nil {
		logger.Error("init bolt db error", err)
		panic(err)
	}
	err = db.GetBoltDBClient(db.DefaultBoltDBClientName).CreateBucketIfNotExists(enum.DefaultNetworkTable)
	if err != nil {
		logger.Error("create network table error", err)
	}
	err = db.GetBoltDBClient(db.DefaultBoltDBClientName).CreateBucketIfNotExists(enum.AllocatedIPKey)
	if err != nil {
		logger.Error("create allocated ip table error", err)
	}
	log.Println("init bolt db finished", db.DefaultBoltDBClientName)
}

// 从指定子网中分配一个可用的IP地址
func AllocateIP(subnet net.IP) (ip net.IP, err error) {
	return subnet, nil
}

// ReleaseIP 释放指定子网中的指定IP地址
// subnet: 要释放IP地址的子网
// ip: 要释放的IP地址
func ReleaseIP(subnet *net.IPNet, ip net.IP) error {
	return nil
}

// 将当前的 allocatedIP 状态保存到 BoltDB 数据库中
func SaveIP() error {
	return nil
}

// 从 BoltDB 数据库中加载分配的 IP 状态
func LoadIP() (err error) {
	return nil
}

// SetInterfaceIP 为指定网络接口设置IP地址
// name: 网络接口名称
// ip: 要设置的IP地址和子网掩码，格式如"192.168.1.1/24"
func SetInterfaceIP(name string, ip string) error {
	// 通过名称获取网络接口(目的是校验接口是否存在)
	link, err := netlink.LinkByName(name)
	logger.Debug("Link:", link, "err:", err)
	if err != nil {
		// 如果获取网络接口失败
		logger.Error("Failed to get link by name: %s", name)
		return err
	}

	// 解析IP
	ipNet, err := netlink.ParseIPNet(ip)
	if err != nil {
		logger.Error("Failed to parse IP net: %s", ip)
		return err
	}

	// 创建地址对象
	addr := &netlink.Addr{IPNet: ipNet, Peer: ipNet}
	// 添加IP地址到网络接口
	if err = netlink.AddrAdd(link, addr); err != nil {
		logger.Error("Failed to add IP address: %s", ip)
		return err
	}
	return nil
}

// SetInterfaceUp 将指定网络接口设置为up状态（激活状态）
// name: 网络接口名称
// 返回值: 设置成功返回nil，否则返回错误信息
func SetInterfaceUp(name string) error {
	// 通过名称获取网络接口链接对象
	link, err := netlink.LinkByName(name)
	logger.Debug("Link:", link, "err:", err)
	if err != nil {
		logger.Error("Failed to get link by name: %s", name)
		return err
	}

	// 设置网络接口为up状态
	if err = netlink.LinkSetUp(link); err != nil {
		logger.Error("Failed to set link up: %s", name)
		return err
	}

	return nil
}

// SetupIptables 设置iptables规则，实现NAT功能
// 容器内部使用私有IP地址，无法直接访问互联网，需要通过NAT转换源地址
// -t nat：指定使用nat表，这是处理NAT转换的表
// -A POSTROUTING：：将规则添加到POSTROUTING链，这是数据包离开本机前的最后一道处理链
// -s：指定源地址范围，即容器子网，只有来自这个子网的数据包才会匹配该规则
// ! -o：指定不从指定网桥接口出去的数据包，"!"表示否定，即不是从该网桥出去的流量
// 通过-s指定源子网和! -o排除特定接口，确保只有容器访问外网的流量才会进行地址伪装
// -j MASQUERADE：执行MASQUERADE动作，将数据包的源IP地址替换为出口接口的IP地址
func SetupIptables(name string, subnet *net.IPNet) error {
	iptableCmd := fmt.Sprintf("-t nat -A POSTROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), name)
	output, err := exec.Command("iptables", strings.Split(iptableCmd, " ")...).Output()
	if err != nil {
		logger.Error("Failed to execute iptables command: %s", output)
		return err
	}
	return nil
}
