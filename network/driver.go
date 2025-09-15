package network

import "fmt"

func NewDriver(driverName string) (driver NetworkDriver, err error) {
	switch driverName {
	case "bridge":
		driver = &BridgeNetworkDriver{}
		return
	default:
		return nil, fmt.Errorf("unsupported network driver: %s", driverName)
	}
}
