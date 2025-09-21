package network

import (
	"fmt"
	"net"
	"strings"

	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/vishvananda/netlink"
)

// 用于生成对端虚拟网络接口的名称前缀
const BridgePeerNamePrefix = "br-peer-"

// 用作桥接网络驱动的方法接收者
type BridgeNetworkDriver struct {
}

// 返回驱动名称"bridge"，用于标识这个网络驱动类型
func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

// 此方法用于创建一个新的桥接网络：
// 1. 创建一个Network对象，包含网络名称、IP范围和驱动类型
// 2. 调用init方法初始化网络
// 3. 返回创建的网络对象
func (b *BridgeNetworkDriver) Create(subnet *net.IPNet, name string) (*Network, error) {
	logger.Debug("create bridge network subnet:", subnet, "name: ", name)
	nw := &Network{
		Name:    name,
		IPRange: subnet,
		Driver:  b.Name(),
	}
	err := b.init(nw)
	if err != nil {
		logger.Error("init bridge error: ", err)
		return nil, err
	}
	return nw, nil
}

// Delete 删除指定的网络
// 1.根据网络名称获取对应的网络接口
// 2.使用netlink库删除该网络接口
// network: 要删除的网络对象，包含网络名称等信息
// 返回值: 删除成功返回nil，否则返回错误信息
func (b *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name
	// 通过网络名称获取对应的网络接口
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		logger.Error("get bridge error: ", err)
		return err
	}
	// 删除网络接口
	return netlink.LinkDel(br)
}

func (b *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	bridgeName := network.Name
	// 获取网桥网络接口
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		logger.Error("get bridge error: ", err)
		return err
	}

	// 创建一个网络链接属性对象，用于配置网络接口的基本属性
	la := netlink.NewLinkAttrs()
	// 将网络端点ID的前5个字符作为网络接口的名称。这是一种命名约定，用于标识不同的网络端点。
	la.Name = endpoint.ID[:5]
	logger.Debug("endpoint name: ", la.Name, "endpoint.ID", endpoint.ID)
	// 设置主接口索引，将这个虚拟网络设备绑定到之前获取的网桥上。MasterIndex指定了该接口的主设备（即网桥）的索引。
	la.MasterIndex = br.Attrs().Index

	// 创建Veth虚拟以太网设备对
	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		// PeerName指定了Veth设备的另一端，后面会连接到网络命名空间
		PeerName: BridgePeerNamePrefix + endpoint.ID[:5],
	}

	// 将Veth设备一端连接到网桥
	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}

	// 设置Veth设备为up状态
	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Add Endpoint Device: %v", err)
	}
	return nil
}

func (b *BridgeNetworkDriver) Disconnect(network *Network, endpoint *Endpoint) error {
	bridgeName := network.Name
	// 获取网桥网络接口
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		logger.Error("get bridge error: ", err)
		return err
	}

	la := netlink.NewLinkAttrs()
	la.Name = endpoint.ID[:5]
	logger.Debug("endpoint name: ", la.Name, "endpoint.ID", endpoint.ID)
	la.MasterIndex = br.Attrs().Index

	endpoint.Device = netlink.Veth{
		LinkAttrs: la,
		PeerName:  BridgePeerNamePrefix + endpoint.ID[:5],
	}
	// 删除Veth设备
	if err := netlink.LinkDel(&endpoint.Device); err != nil {
		return fmt.Errorf("Error Delete Endpoint Device: %v", err)
	}
	return nil
}

// 初始化bridge网络
// 1.创建bridge网络接口
// 2.设置bridge网络接口的IP地址
// 3.设置bridge网络接口为up状态
// 4.配置iptables的NAT规则，实现网络地址转换功能
func (b *BridgeNetworkDriver) init(network *Network) error {
	// 创建bridge网络接口
	err := CreateBridgeInterface(network.Name)
	if err != nil {
		logger.Error("create bridge error: ", err)
		return err
	}

	// 创建network.IPRange的副本，这个副本包含了子网的IP地址和子网掩码
	bridgeIP := *network.IPRange
	// 确保bridgeIP的IP字段设置为子网的基准IP地址
	// 目的是为了确保网桥使用子网的第一个IP地址作为其IP地址。
	// 例如，如果子网是192.168.3.2/24，那么网桥将使用192.168.3.0作为其IP地址
	// 然后通过bridgeIP.String()方法，会得到一个标准的CIDR格式字符串（如"192.168.3.0/24"）
	// 这样就可以被SetInterfaceIP函数中的netlink.ParseIPNet正确解析。
	bridgeIP.IP = network.IPRange.IP

	// 设置bridge网络接口的IP地址
	err = SetInterfaceIP(network.Name, bridgeIP.String())
	if err != nil {
		logger.Error("set interface ip error: ", err)
		return err
	}

	// 设置bridge网络接口为up状态
	err = SetInterfaceUp(network.Name)
	if err != nil {
		logger.Error("set interface up error: ", err)
		return err
	}

	// 配置iptables的NAT规则，实现网络地址转换功能
	err = SetupIptables(network.Name, network.IPRange)
	if err != nil {
		logger.Error("setup iptables error: ", err)
		return err
	}
	return nil
}

// 创建网桥接口，接收网桥名称作为参数，返回错误信息
func CreateBridgeInterface(name string) error {
	inter, err := net.InterfaceByName(name)
	if inter != nil {
		return nil
	}
	// 如果接口已经存在，直接返回，无需重复创建
	if err != nil && !strings.Contains(err.Error(), "no such network interface") {
		logger.Error("get interface error: ", err)
		return err
	}
	// 创建一个新的网络链接属性对象，用于配置网络接口的基本属性
	la := netlink.NewLinkAttrs()
	// 设置网桥接口的名称
	la.Name = name
	// 创建一个Bridge对象，表示一个网桥网络接口，并将之前配置的属性应用到该对象上
	bridge := netlink.Bridge{LinkAttrs: la}
	// 使用netlink库的LinkAdd方法将创建的网桥接口添加到系统中
	if err = netlink.LinkAdd(&bridge); err != nil {
		logger.Error("create bridge error: ", err)
		return err
	}
	return nil

}
