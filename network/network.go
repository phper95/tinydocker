package network

import (
	"github.com/vishvananda/netlink"
	"net"
)

type NetworkDriver interface {
	Name() string // 返回驱动名称
	Create(subnet string, name string) (*Network, error)
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
