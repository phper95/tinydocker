package network

import (
	"encoding/json"
	"fmt"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/db"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/vishvananda/netlink"
	"math"
	"net"
	"os/exec"
	"strings"
	"sync"
)

// 存储每个子网的 IP 分配状态
var allocatedIP map[string]string
var lock sync.Mutex

// 从指定子网中分配一个可用的IP地址
func AllocateIP(subnet *net.IPNet) (ip net.IP, err error) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := allocatedIP[subnet.String()]; !ok {
		// bits 是 IP 地址的总位数
		// 对于 IPv4 是 32(4个字节 × 8位/字节)
		// IPv6 是 128位（16个字节 × 8位/字节）
		// ones 是子网掩码中 1 的个数（例如 /24 表示掩码中有 24 个 1）
		// ipv4 地址由4个字节组成，每个字节8位
		// 192.168.1.0/24 将子网掩码转成二进制是：
		// 前 3 个字节（共 24 位）：全为 1 → 每个字节的二进制是 11111111；
		// 第 4 个字节（剩余 8 位）：全为 0 → 二进制是 00000000
		// 11111111.11111111.11111111.00000000（十进制为 255.255.255.0）

		ones, bits := subnet.Mask.Size()
		// 初始化填充所有未分配的IP地址，用0表示未分配，1表示已分配
		// // 掩码 /24 表示前24位是网络位，后8位是主机位
		//		// bits-ones 计算的是主机位的数量
		//		// 如果子网是 192.168.1.0/24：
		//		// bits = 32（IPv4）                  ones = 24（/24 表示掩码中 24 个 1）
		//		// bits - ones = 32 - 24 = 8         2^8 = 256
		//		// 所以会创建一个包含 256 个 "0" 的字符串
		allocatedIP[subnet.String()] = strings.Repeat("0", int(math.Pow(2, float64(bits-ones))))
	} else {
		// 校验子网是否已经分配完毕
		if strings.Count(allocatedIP[subnet.String()], "0") == 0 {
			return nil, fmt.Errorf("subnet %s is full", subnet.String())
		}
	}
	for c := range allocatedIP[subnet.String()] {
		// 找到第一个未分配的IP地址
		if allocatedIP[subnet.String()][c] == '0' {
			ipalloc := []byte(allocatedIP[subnet.String()])
			ipalloc[c] = '1'
			allocatedIP[subnet.String()] = string(ipalloc)
			ip = subnet.IP
			// 循环处理IP地址的4个字节（IPv4地址由4个字节组成）
			for t := uint(4); t > 0; t -= 1 {
				// []byte(ip)[4-t]获取IP地址的每个字节
				// uint8(c >> ((t - 1) * 8))通过位运算将索引c转换为对应字节的值
				// ((t - 1) * 8)表示字节的偏移量，因为每个字节8位，所以需要乘以8
				// 假设子网是192.168.1.0/24，我们要分配第5个IP地址（索引c=4）：
				// 子网基础IP：192.168.1.0
				// 索引c=4表示这是第5个IP地址
				// 通过位运算将4转换为IP地址：
				// t=4: 192 + (4 >> 24) = 192 + 0 = 192
				// t=3: 168 + (4 >> 16) = 168 + 0 = 168
				// t=2: 1 + (4 >> 8) = 1 + 0 = 1
				// t=1: 0 + (4 >> 0) = 0 + 4 = 4
				// 最后将最后一个字节加1: 4 + 1 = 5
				// 	最终得到IP地址：192.168.1.5
				[]byte(ip)[4-t] += uint8(c >> ((t - 1) * 8))
			}
			// 分配的IP地址加1，跳过网络地址
			ip[3] += 1
			break
		}
	}
	err = SaveIP()
	fmt.Println("IP:", ip)
	return ip, err
}

// ReleaseIP 释放指定子网中的指定IP地址
// subnet: 要释放IP地址的子网
// ip: 要释放的IP地址
func ReleaseIP(subnet *net.IPNet, ip net.IP) error {
	// 检查子网是否存在
	subnetStr := subnet.String()
	ipStr := ip.String()

	if _, ok := allocatedIP[subnetStr]; !ok {
		return fmt.Errorf("subnet %s not found", subnetStr)
	}

	// 计算要释放的IP在位图中的索引
	// 需要减去1，因为在分配时有加1的操作
	ipIndex := 0
	ipBytes := []byte(ip.To4())
	subnetBytes := []byte(subnet.IP.To4())

	for i := range ipBytes {
		diff := uint(ipBytes[i] - subnetBytes[i])
		ipIndex += int(diff << ((3 - uint(i)) * 8))
	}

	// 减去网络地址偏移
	ipIndex -= 1

	// 检查索引是否有效
	ipAlloc := allocatedIP[subnetStr]
	if ipIndex < 0 || ipIndex >= len(ipAlloc) {
		return fmt.Errorf("invalid IP address %s for subnet %s", ipStr, subnetStr)
	}

	// 将对应位置标记为未分配('0')
	ipalloc := []byte(ipAlloc)
	ipalloc[ipIndex] = '0'
	allocatedIP[subnetStr] = string(ipalloc)

	// 保存更新后的分配状态
	err := SaveIP()
	if err != nil {
		logger.Error("Failed to save allocated IP", allocatedIP, "err:", err)
		return err
	}

	logger.Info("Released IP address %s from subnet %s", ipStr, subnetStr)
	return nil
}

// 将当前的 allocatedIP 状态保存到 BoltDB 数据库中
func SaveIP() error {
	jsonBytes, err := json.Marshal(allocatedIP)
	if err != nil {
		logger.Error("Failed to marshal allocated IP: %+v", allocatedIP, "err:", err)
		return err
	}
	err = db.GetBoltDBClient("").Put(enum.DefaultNetworkTable, enum.AllocatedIPKey, jsonBytes)
	if err != nil {
		logger.Error("Failed to save allocated IP: ", string(jsonBytes), "err:", err)
	}
	return err
}

// 从 BoltDB 数据库中加载分配的 IP 状态
func LoadIP() (err error) {
	jsonBytes, err := db.GetBoltDBClient("").Get(enum.DefaultNetworkTable, enum.AllocatedIPKey)
	if err != nil {
		logger.Error("Failed to load allocated IP: ", "err:", err)
		return err
	}

	if jsonBytes == nil {
		allocatedIP = make(map[string]string)
		return nil
	}

	if err = json.Unmarshal(jsonBytes, &allocatedIP); err != nil {
		logger.Error("Failed to unmarshal allocated IP: ", string(jsonBytes), "err:", err)
		return err
	}
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
	logger.Info("Setup iptables command: %s", iptableCmd)
	output, err := exec.Command("iptables", strings.Split(iptableCmd, " ")...).Output()
	if err != nil {
		logger.Error("Failed to execute iptables command: %s", output)
		return err
	}
	logger.Info("Setup iptables rule: %s", string(output))
	return nil
}
