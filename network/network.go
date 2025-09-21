package network

import (
	"encoding/json"
	"fmt"
	"github.com/phper95/tinydocker/container/models"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/db"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

type NetworkDriver interface {
	Name() string // 返回驱动名称
	Create(subnet *net.IPNet, name string) (*Network, error)
	Delete(nework Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network *Network, endpoint *Endpoint) error
}
type Network struct {
	Name    string
	IPRange *net.IPNet
	Driver  string
}
type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"device"`
	IPAddress   net.IP           `json:"ip"`
	MACAddress  net.HardwareAddr `json:"mac"`
	Network     *Network         `json:"network"`
	PortMapping []string         `json:"port_mapping"`
}

const (
	InterfaceLoName = "lo"
)

// 1. 解析用户输入的子网信息，确保格式正确
// 2. 调用指定的网络驱动创建网络
// 3. 将网络信息保存到数据库中，便于后续管理
func CreateNetwork(name, driver, subnet string) error {
	// 判断网络是否存在
	nw, err := GetNetworkFromDB(name)
	if err != nil {
		logger.Error("get network from db error: ", err)
		return err
	}
	if nw != nil {
		logger.Error("network %s already exists", name)
		return fmt.Errorf("network %s already exists", name)
	}
	// 解析子网
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Error("parse subnet error: ", err)
		return err
	}
	err = LoadIP()
	if err != nil {
		logger.Error("Failed to load allocated IP", allocatedIP, "err:", err)
	} else {
		logger.Info("Loaded allocated IP: ", allocatedIP)
	}
	// 在子网中分配IP
	ip, err := AllocateIP(ipNet)
	if err != nil {
		logger.Error("allocate ip error: ", err)
		return err
	}
	ipNet.IP = ip
	dr, err := NewDriver(driver)
	if err != nil {
		logger.Error("create driver error: ", err)
		return err
	}
	nw, err = dr.Create(ipNet, name)
	if err != nil {
		logger.Error("create network error: ", err)
		return err
	}
	return nw.Save()
}

func ListNetwork() {
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tIPRANGE\tDRIVER")
	data, err := db.GetBoltDBClient("").GetAll(enum.DefaultNetworkTable)
	if err != nil {
		logger.Error("list network error: ", err)
		return
	}
	for _, v := range data {
		nw := &Network{}
		if err := json.Unmarshal(v, nw); err != nil {
			logger.Error("list network error: ", err)
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", nw.Name, nw.IPRange.String(), nw.Driver)
	}
	if err := w.Flush(); err != nil {
		logger.Error("list network error: ", err)
	}
	return
}

func DeleteNetwork(name string) error {
	nw, err := GetNetworkFromDB(name)
	if err != nil {
		logger.Error("get network from db error: ", err)
		return err
	}
	if nw == nil {
		logger.Error("network %s not exists", name)
		return fmt.Errorf("network %s not exists", name)
	}

	// 释放IP
	err = ReleaseIP(nw.IPRange, nw.IPRange.IP)
	if err != nil {
		logger.Error("release ip error: ", err)
		return err
	}

	dr, err := NewDriver(nw.Driver)
	if err != nil {
		logger.Error("create driver error: ", err)
		return err
	}
	err = dr.Delete(*nw)
	if err != nil {
		logger.Error("delete network error: ", err)
		return err
	}
	return nw.Delete(name)
}
func (nw *Network) Save() error {
	data, _ := json.Marshal(nw)
	err := db.GetBoltDBClient("").Put(enum.DefaultNetworkTable, nw.Name, data)
	if err != nil {
		logger.Error("save ip error: %v", err)
		return err
	}
	return nil
}

func (nw *Network) Delete(name string) error {
	err := db.GetBoltDBClient("").Delete(enum.DefaultNetworkTable, name)
	if err != nil {
		logger.Error("delete ip error: %v", err)
		return err
	}
	return nil
}

func GetNetworkFromDB(name string) (network *Network, err error) {
	data, err := db.GetBoltDBClient("").Get(enum.DefaultNetworkTable, name)
	if err != nil {
		logger.Error("load ip error: %v", err)
		return
	}
	if data == nil {
		return nil, nil
	}
	err = json.Unmarshal(data, &network)
	if err != nil {
		logger.Error("unmarshal ip error: %v", err)
	}
	return
}

func Connect(name string, containerInfo *models.Info) (ip net.IP, err error) {
	nw, err := GetNetworkFromDB(name)
	if err != nil {
		logger.Error("get network from db error: ", err)
		return
	}
	if nw == nil {
		logger.Error("network %s not exists", name)
		err = fmt.Errorf("network %s not exists", name)
		return
	}

	// 加载已分配的IP
	err = LoadIP()
	if err != nil {
		logger.Error("Failed to load allocated IP", allocatedIP, "err:", err)
		return nil, err
	}

	fmt.Println("nw.IPRange=", nw.IPRange)
	// 注意，这里需要解析子网
	_, subnet, err := net.ParseCIDR(nw.IPRange.String())
	if err != nil {
		logger.Error("parse subnet error: ", err)
		return
	}
	// 分配容器IP
	ip, err = AllocateIP(subnet)
	if err != nil {
		logger.Error("allocate ip error: ", err)
		return
	}
	logger.Info("Connect network: %+v, ip: %s", nw, ip.String())
	// 创建网络端点
	ep := &Endpoint{
		ID:          fmt.Sprintf("%s-%s", containerInfo.Id, name),
		IPAddress:   ip,
		Network:     nw,
		PortMapping: containerInfo.PortMapping,
	}
	dr, err := NewDriver(nw.Driver)
	if err != nil {
		logger.Error("create driver error: ", err)
		return
	}
	// 网络驱动挂载和配置网络端点
	err = dr.Connect(nw, ep)
	if err != nil {
		logger.Error("connect network error: ", err)
		return
	}

	// 在容器网络命名空间中配置网络接口
	err = configEndpointNetwork(ep, containerInfo)
	if err != nil {
		logger.Error("config endpoint network error: ", err)
		return
	}

	// 配置容器的端口映射
	err = configPortMapping(ep, containerInfo)
	if err != nil {
		logger.Error("config port mapping error: ", err)
		return
	}
	logger.Info("connect container %+v to endpoint %+v success", *containerInfo, *ep)
	return ip, nil
}

// 在容器网络命名空间中配置网络接口，包括设置IP地址、激活接口、激活回环设备以及添加默认路由，使容器能够正常进行网络通信。
func configEndpointNetwork(ep *Endpoint, containerInfo *models.Info) error {
	// 根据端点设备的PeerName获取对应的网络接口对象
	var peerLink netlink.Link
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		peerLink, err = netlink.LinkByName(ep.Device.PeerName)
		if err == nil {
			break
		}
		if i < maxRetries-1 {
			// 等待一段时间后重试
			time.Sleep(1 * time.Second)
		}
	}
	if err != nil {
		logger.Error("get peer link error: ", err)
		return err
	}
	// 将当前线程切换到容器的网络命名空间
	// 返回一个清理函数，用于恢复原始的网络命名空间(注意configNetNs函数中的逻辑会在当前行执行，但最终返回的函数会在configEndpointNetwork函数返回之前执行，具体执行机制可以参考tinydocker\gobase\9-8\main.go的示例代码)
	defer configNetNs(&peerLink, containerInfo)()
	ip := *ep.Network.IPRange
	ip.IP = ep.IPAddress
	log.Println("ip.string(): ", ip.String(), " ep.IPAddress: ", "ep.IPAddress", ip.IP.String())

	// 为端点设备的Peer接口设置IP地址
	err = SetInterfaceIP(ep.Device.PeerName, ip.String())
	if err != nil {
		logger.Error("set interface ip error: ", err)
		return err
	}

	// 将端点设备的Peer接口设置为up状态（激活状态）
	err = SetInterfaceUp(ep.Device.PeerName)
	if err != nil {
		logger.Error("set interface up error: ", err)
		return err
	}

	// 激活容器的回环接口（loopback interface）
	err = SetInterfaceUp(InterfaceLoName)
	if err != nil {
		logger.Error("set interface lo up error: ", err)
		return err
	}

	// 解析默认路由的CIDR表示，"0.0.0.0/0"代表所有IP地址
	_, ipNet, _ := net.ParseCIDR("0.0.0.0/0")

	// 添加默认路由规则
	defaultRoute := &netlink.Route{
		// 使用peerLink的接口索引
		LinkIndex: peerLink.Attrs().Index,
		// 网关地址设置为网络范围的IP（通常是网桥的IP）
		Gw: ep.Network.IPRange.IP,
		// 目标网络设置为所有IP地址（0.0.0.0/0）
		Dst: ipNet,
	}
	// 将默认路由添加到容器的网络路由表中(这里不做严格错误检查，如果路由已存在也会报错)
	e := netlink.RouteAdd(defaultRoute)
	if e != nil {
		logger.Error("add default route error: ", err)
	}
	return err

}

// 进入容器网络命名空间配置虚拟以太网对在容器内的一端的网络

func configNetNs(peerLink *netlink.Link, containerInfo *models.Info) func() {
	// 获取容器网络命名空间
	// 通过/proc/<pid>/ns/net文件获取容器网络命名空间的文件描述符来操作容器网络命名空间
	// containerInfo的Pid是容器在宿主机上的进程ID
	f, err := os.OpenFile(fmt.Sprintf("/proc/%d/ns/net", containerInfo.Pid), os.O_RDONLY, 0)
	if err != nil {
		logger.Error("open container netns error: ", err)
		return nil
	}
	// 将容器网络命名空间的文件描述符
	nsFd := f.Fd()
	// 锁定当前的OS线程,因为Go语言中Goroutine的调度是由Go的Runtime管理的，使用的是GMP调度模型，
	// GMP调度模型的本质是通过逻辑处理器（P）作为中间调度层，将Goroutine映射到操作系统线程（M）上执行，
	// 从而避免操作系统线程直接管理调度Goroutine导致频繁的上下文切换来提高Go语言的并发性能。
	// 因此goroutine可能会在不同的OS线程上执行，而网络命名空间的切换是针对OS线程的，
	// 如果不锁定当前线程，可能会导致网络命名空间切换后，Goroutine被调度到另一个线程上执行，从而无法正确访问容器的网络资源。
	// 通过runtime.LockOSThread()函数，可以确保当前的Goroutine在执行期间始终绑定到同一个OS线程，
	// 这样在切换网络命名空间后，Goroutine仍然能够正确访问和操作容器的网络资源，避免了潜在的错误和不一致性。
	runtime.LockOSThread()

	// 修改veth peer 另外一端移到容器的namespace中
	err = netlink.LinkSetNsFd(*peerLink, int(nsFd))
	if err != nil {
		logger.Error("set container netns error: ", err)
		return nil
	}
	// 获取当前进程的网络命名空间，保存原始命名空间以便后续恢复
	origns, err := netns.Get()
	if err != nil {
		logger.Error("get current netns error: ", err)
		return nil
	}
	// 将当前进程切换到容器的网络命名空间
	err = netns.Set(netns.NsHandle(nsFd))
	if err != nil {
		logger.Error("set netns error: ", err)
		return nil
	}
	// 在容器的网络命名空间中执行完容器的网络配置之后将程序恢复到原始的网络命名空间
	return func() {
		// 恢复到原始网络命名空间
		err = netns.Set(origns)
		if err != nil {
			logger.Error("set netns error: ", err)
		}
		// 关闭原始命名空间文件描述符
		err = origns.Close()
		if err != nil {
			logger.Error("close netns error: ", err)
		}
		// 解锁OS线程
		runtime.UnlockOSThread()
		// 关闭容器命名空间文件描述符
		err = f.Close()
		if err != nil {
			logger.Error("close netns file error: ", err)
		}
	}

}

// 配置宿主机到容器的端口映射
// 通过iptables的DNAT规则来实现宿主机上的请求转发到容器上
func configPortMapping(ep *Endpoint, containerInfo *models.Info) error {
	for _, pm := range ep.PortMapping {
		portMapping := strings.Split(pm, ":")
		if len(portMapping) != 2 {
			logger.Error("invalid port mapping: %s", pm)
			continue
		}
		// 添加iptables规则
		// -t nat：指定使用nat表，这是处理NAT转换的表
		// -A：在PREROUTING链上添加规则, PREROUTING链用于在数据包到达路由决策之前进行处理
		// -p tcp：指定协议为TCP
		// --dport：指定目标端口
		// -j DNAT：指定目标地址转换
		// --to-destination：指定目标地址和端口
		iptablesCmd := fmt.Sprintf("-t nat -A PREROUTING -p tcp -m tcp  --dport %s -j DNAT --to-destination %s:%s",
			portMapping[0], ep.IPAddress.String(), portMapping[1])
		output, err := exec.Command("iptables", strings.Split(iptablesCmd, " ")...).Output()
		if err != nil {
			logger.Error("add iptables rule error:", err, "output:", output)
			continue
		}
		logger.Info("add iptables rule: %s", output)
	}
	return nil
}
