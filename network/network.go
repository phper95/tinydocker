package network

import (
	"encoding/json"
	"fmt"
	"github.com/phper95/tinydocker/enum"
	"github.com/phper95/tinydocker/pkg/db"
	"github.com/phper95/tinydocker/pkg/logger"
	"github.com/vishvananda/netlink"
	"net"
	"os"
	"text/tabwriter"
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
	allocateIP, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		logger.Error("parse subnet error: ", err)
		return err
	}
	// 在子网中分配IP
	ip, err := AllocateIP(allocateIP)
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
func (n *Network) Save() error {
	data, _ := json.Marshal(n)
	err := db.GetBoltDBClient("").Put(enum.DefaultNetworkTable, n.Name, data)
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
