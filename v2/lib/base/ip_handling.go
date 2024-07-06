package base

import (
	"fmt"
	"log"
	"net"
)

// Get preferred outbound ip of this machine
func GetOutboundIP() (net.IP, error) {
	addrs, err1 := net.InterfaceAddrs()
	if err1 != nil {
		return nil, err1
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("no iip found")
}

// Get preferred outbound ip of this machine within 192.168.0.0/24
func GetLocalNetOutboundIP() (net.IP, error) {
	addrs, err1 := net.InterfaceAddrs()
	if err1 != nil {
		return nil, err1
	}

	_, localnet, err2 := net.ParseCIDR("192.168.0.0/24")
	if err2 != nil {
		return nil, err2
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && localnet.Contains(ipnet.IP.To4()) {
				log.Println("local ip address: ", ipnet)
				return ipnet.IP, nil
			}
		}
	}
	return nil, fmt.Errorf("no iip found")
}
